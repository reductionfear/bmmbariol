package server

import (
	"fmt"
	"strings"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

const (
	namespace = "chesshook-intermediary"
	version   = "1"
)

// MessageHandler handles incoming WebSocket messages
type MessageHandler struct {
	authManager     *AuthManager
	connManager     *ConnectionManager
	engineName      string
	engineInputChan chan string
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(authMgr *AuthManager, connMgr *ConnectionManager, engineInputChan chan string) *MessageHandler {
	return &MessageHandler{
		authManager:     authMgr,
		connManager:     connMgr,
		engineName:      "Unknown",
		engineInputChan: engineInputChan,
	}
}

// SetEngineName sets the engine name for whatengine responses
func (h *MessageHandler) SetEngineName(name string) {
	h.engineName = name
}

// HandleMessage handles an incoming message and returns the response
func (h *MessageHandler) HandleMessage(message string, connection *Connection) (string, bool) {
	utils.Logger.Debugf("Handling message: %s", message)

	message = strings.TrimSpace(message)
	if message == "" {
		return "", false
	}

	parts := strings.Split(message, " ")
	command := parts[0]

	session := connection.GetSession()

	// Handle commands
	switch command {
	case "whoareyou":
		return h.handleWhoAreYou()

	case "whatengine":
		return h.handleWhatEngine(session)

	case "auth":
		return h.handleAuth(parts, session)

	case "sub":
		return h.handleSubscribe(session)

	case "unsub":
		return h.handleUnsubscribe(session)

	case "lock":
		return h.handleLock(connection)

	case "unlock":
		return h.handleUnlock(connection)

	default:
		return h.handleUCICommand(message, session)
	}
}

// handleWhoAreYou handles the whoareyou command
func (h *MessageHandler) handleWhoAreYou() (string, bool) {
	return fmt.Sprintf("iam %sv%s", namespace, version), true
}

// handleWhatEngine handles the whatengine command
func (h *MessageHandler) handleWhatEngine(session *UserSession) (string, bool) {
	if h.authManager.RequiresAuthForRead() && !session.IsAuthenticated() {
		return "autherr", true
	}
	return fmt.Sprintf("engine %s", h.engineName), true
}

// handleAuth handles the auth command
func (h *MessageHandler) handleAuth(parts []string, session *UserSession) (string, bool) {
	if len(parts) < 2 {
		return "autherr", true
	}

	providedKey := parts[1]

	// Check if already blocked
	if session.IsBlocked() {
		return "autherr", true
	}

	// Validate passkey
	if h.authManager.ValidateAuth(providedKey) {
		session.Authenticate()
		return "authok", true
	}

	// Invalid passkey
	attempts := session.IncrementFailedAttempts()
	utils.Logger.Warnf("Failed authentication attempt (attempt %d/3)", attempts)
	return "autherr", true
}

// handleSubscribe handles the sub command
func (h *MessageHandler) handleSubscribe(session *UserSession) (string, bool) {
	if h.authManager.RequiresAuthForRead() && !session.IsAuthenticated() {
		return "autherr", true
	}

	if session.Subscribe() {
		utils.Logger.Debugf("User subscribed (Total subscribers: %d)", h.connManager.GetSubscriberCount())
		return "subok", true
	}

	return "suberr", true
}

// handleUnsubscribe handles the unsub command
func (h *MessageHandler) handleUnsubscribe(session *UserSession) (string, bool) {
	if h.authManager.RequiresAuthForRead() && !session.IsAuthenticated() {
		return "autherr", true
	}

	if session.Unsubscribe() {
		utils.Logger.Debugf("User unsubscribed (Remaining subscribers: %d)", h.connManager.GetSubscriberCount())
		return "unsubok", true
	}

	return "unsuberr", true
}

// handleLock handles the lock command
func (h *MessageHandler) handleLock(connection *Connection) (string, bool) {
	session := connection.GetSession()

	if h.authManager.RequiresAuthForWrite() && !session.IsAuthenticated() {
		return "autherr", true
	}

	if h.connManager.TryAcquireEngineLock(connection.conn) {
		utils.Logger.Infof("Engine lock acquired by %s", connection.GetRemoteAddr())
		return "lockok", true
	}

	return "lockerr", true
}

// handleUnlock handles the unlock command
func (h *MessageHandler) handleUnlock(connection *Connection) (string, bool) {
	session := connection.GetSession()

	if h.authManager.RequiresAuthForWrite() && !session.IsAuthenticated() {
		return "autherr", true
	}

	if !session.HasLock() {
		return "unlockerr", true
	}

	if h.connManager.ReleaseEngineLock(connection.conn) {
		session.ReleaseLock()
		utils.Logger.Infof("Engine lock released by %s", connection.GetRemoteAddr())
		return "unlockok", true
	}

	return "unlockerr", true
}

// handleUCICommand handles UCI commands
func (h *MessageHandler) handleUCICommand(message string, session *UserSession) (string, bool) {
	if h.authManager.RequiresAuthForWrite() && !session.IsAuthenticated() {
		return "autherr", true
	}

	// Forward to engine
	utils.Logger.Debugf("Forwarding UCI command to engine: %s", message)
	select {
	case h.engineInputChan <- message:
		// Command queued successfully
		return "", false // No immediate response
	default:
		utils.Logger.Warnf("Engine input channel full, command dropped: %s", message)
		return "", false
	}
}

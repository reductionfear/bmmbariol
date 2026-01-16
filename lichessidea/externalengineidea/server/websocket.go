package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WebSocketServer manages WebSocket connections and message routing
type WebSocketServer struct {
	authManager    *AuthManager
	connManager    *ConnectionManager
	messageHandler *MessageHandler
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(authMgr *AuthManager, engineInputChan chan string) *WebSocketServer {
	connMgr := NewConnectionManager()
	msgHandler := NewMessageHandler(authMgr, connMgr, engineInputChan)

	return &WebSocketServer{
		authManager:    authMgr,
		connManager:    connMgr,
		messageHandler: msgHandler,
	}
}

// SetEngineName sets the engine name for the server
func (s *WebSocketServer) SetEngineName(name string) {
	s.messageHandler.SetEngineName(name)
}

// BroadcastEngineOutput broadcasts engine output to all subscribers
func (s *WebSocketServer) BroadcastEngineOutput(output string) {
	s.connManager.BroadcastToSubscribers(output)
}

// HandleWebSocket handles incoming WebSocket connections
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Logger.Errorf("WebSocket upgrade failed: %v", err)
		return
	}

	// Set read limit to prevent abuse
	conn.SetReadLimit(64 * 1024) // 64 KB max message size

	remoteAddr := conn.RemoteAddr().String()
	utils.Logger.Infof("New WebSocket connection from: %s", remoteAddr)

	// Check if localhost for auto-auth
	autoAuth := s.authManager.IsLocalhostConnection(remoteAddr)
	if autoAuth {
		utils.Logger.Debugf("Auto-authenticating localhost connection: %s", remoteAddr)
	}

	// Create session and add connection
	session := NewUserSession(autoAuth)
	connection := s.connManager.AddConnection(conn, session)

	// Cleanup on disconnect
	defer func() {
		s.connManager.RemoveConnection(conn)
		connection.Close()
		utils.Logger.Infof("WebSocket connection closed: %s", remoteAddr)
	}()

	// Message loop
	for {
		// Read message
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Logger.Warnf("WebSocket error: %v", err)
			}
			break
		}

		messageStr := string(message)
		utils.Logger.Debugf("Received from %s: %s", remoteAddr, messageStr)

		// Handle message
		response, shouldRespond := s.messageHandler.HandleMessage(messageStr, connection)

		if shouldRespond && response != "" {
			if err := connection.Send(response); err != nil {
				utils.Logger.Errorf("Failed to send response: %v", err)
				break
			}
		}
	}
}

// GetConnectionCount returns the number of active connections
func (s *WebSocketServer) GetConnectionCount() int {
	return s.connManager.GetConnectionCount()
}

// GetSubscriberCount returns the number of subscribed connections
func (s *WebSocketServer) GetSubscriberCount() int {
	return s.connManager.GetSubscriberCount()
}

package server

import (
	"crypto/rand"
	"math/big"
	"net"
	"strings"
	"sync"
)

const (
	maxFailedAttempts = 3
	passKeyLength     = 10
)

// AuthManager handles authentication for WebSocket connections
type AuthManager struct {
	passKey          string
	requireAuthWrite bool
	requireAuthRead  bool
	localhostBypass  bool
	mu               sync.RWMutex
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(requireWrite, requireRead, localhostBypass bool, passKey string) *AuthManager {
	if passKey == "" {
		passKey = generatePassKey()
	}

	return &AuthManager{
		passKey:          passKey,
		requireAuthWrite: requireWrite,
		requireAuthRead:  requireRead,
		localhostBypass:  localhostBypass,
	}
}

// GetPassKey returns the current passkey
func (a *AuthManager) GetPassKey() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.passKey
}

// IsLocalhostConnection checks if the connection is from localhost
func (a *AuthManager) IsLocalhostConnection(remoteAddr string) bool {
	if !a.localhostBypass {
		return false
	}

	// Extract host from address (remove port)
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If no port, use the whole address
		host = remoteAddr
	}

	// Check for localhost addresses
	return host == "127.0.0.1" || host == "::1" || host == "localhost"
}

// ValidateAuth checks if the provided passkey is correct
func (a *AuthManager) ValidateAuth(providedKey string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return providedKey == a.passKey
}

// RequiresAuthForWrite checks if authentication is required for write operations
func (a *AuthManager) RequiresAuthForWrite() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.requireAuthWrite
}

// RequiresAuthForRead checks if authentication is required for read operations
func (a *AuthManager) RequiresAuthForRead() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.requireAuthRead
}

// generatePassKey generates a random passkey
func generatePassKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, passKeyLength)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a simpler method if crypto/rand fails
			result[i] = charset[i%len(charset)]
		} else {
			result[i] = charset[num.Int64()]
		}
	}

	return string(result)
}

// UserSession tracks a user's authentication and subscription state
type UserSession struct {
	isAuthenticated   bool
	incorrectAttempts int
	isSubscribed      bool
	hasLock           bool
	mu                sync.RWMutex
}

// NewUserSession creates a new user session
func NewUserSession(autoAuth bool) *UserSession {
	return &UserSession{
		isAuthenticated:   autoAuth,
		incorrectAttempts: 0,
		isSubscribed:      false,
		hasLock:           false,
	}
}

// IsAuthenticated checks if the user is authenticated
func (u *UserSession) IsAuthenticated() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.isAuthenticated
}

// Authenticate authenticates the user
func (u *UserSession) Authenticate() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.isAuthenticated = true
}

// IncrementFailedAttempts increments the failed authentication attempts counter
func (u *UserSession) IncrementFailedAttempts() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.incorrectAttempts++
	return u.incorrectAttempts
}

// IsBlocked checks if the user is blocked due to too many failed attempts
func (u *UserSession) IsBlocked() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.incorrectAttempts >= maxFailedAttempts
}

// IsSubscribed checks if the user is subscribed to engine output
func (u *UserSession) IsSubscribed() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.isSubscribed
}

// Subscribe subscribes the user to engine output
func (u *UserSession) Subscribe() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.isSubscribed {
		return false // Already subscribed
	}
	u.isSubscribed = true
	return true
}

// Unsubscribe unsubscribes the user from engine output
func (u *UserSession) Unsubscribe() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	if !u.isSubscribed {
		return false // Not subscribed
	}
	u.isSubscribed = false
	return true
}

// HasLock checks if the user has the engine lock
func (u *UserSession) HasLock() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.hasLock
}

// AcquireLock acquires the engine lock
func (u *UserSession) AcquireLock() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.hasLock = true
}

// ReleaseLock releases the engine lock
func (u *UserSession) ReleaseLock() {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.hasLock = false
}

// CheckCommand checks if the user is authorized to execute a command
func (u *UserSession) CheckCommand(cmd string, authMgr *AuthManager) (needsAuth bool, authFailed bool) {
	parts := strings.Split(cmd, " ")
	if len(parts) == 0 {
		return false, false
	}

	command := parts[0]

	// Commands that don't require authentication
	if command == "whoareyou" {
		return false, false
	}

	// Read commands
	readCommands := map[string]bool{
		"whatengine": true,
		"sub":        true,
		"unsub":      true,
	}

	// Write commands (everything else except auth)
	if command == "auth" {
		return false, false
	}

	if readCommands[command] {
		if authMgr.RequiresAuthForRead() && !u.IsAuthenticated() {
			return true, true
		}
		return false, false
	}

	// Write commands
	if authMgr.RequiresAuthForWrite() && !u.IsAuthenticated() {
		return true, true
	}

	return false, false
}

package server

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// Connection represents a WebSocket connection with its session
type Connection struct {
	conn    *websocket.Conn
	session *UserSession
	mu      sync.Mutex
}

// NewConnection creates a new connection wrapper
func NewConnection(conn *websocket.Conn, session *UserSession) *Connection {
	return &Connection{
		conn:    conn,
		session: session,
	}
}

// Send sends a message to the connection
func (c *Connection) Send(message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		utils.Logger.Errorf("Error sending message: %v", err)
		return err
	}
	return nil
}

// Close closes the connection
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}

// GetSession returns the user session
func (c *Connection) GetSession() *UserSession {
	return c.session
}

// GetRemoteAddr returns the remote address
func (c *Connection) GetRemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// ConnectionManager manages all active WebSocket connections
type ConnectionManager struct {
	connections map[*websocket.Conn]*Connection
	mu          sync.RWMutex
	engineLock  sync.Mutex
	engineLocked bool
	lockHolder   *websocket.Conn
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections:  make(map[*websocket.Conn]*Connection),
		engineLocked: false,
		lockHolder:   nil,
	}
}

// AddConnection adds a connection to the manager
func (m *ConnectionManager) AddConnection(conn *websocket.Conn, session *UserSession) *Connection {
	m.mu.Lock()
	defer m.mu.Unlock()

	connection := NewConnection(conn, session)
	m.connections[conn] = connection
	utils.Logger.Infof("New connection added: %s (Total: %d)", conn.RemoteAddr().String(), len(m.connections))
	return connection
}

// RemoveConnection removes a connection from the manager
func (m *ConnectionManager) RemoveConnection(conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if connection, exists := m.connections[conn]; exists {
		// Release engine lock if held
		if connection.session.HasLock() {
			m.releaseEngineLockInternal(conn)
		}
		
		delete(m.connections, conn)
		utils.Logger.Infof("Connection removed: %s (Remaining: %d)", conn.RemoteAddr().String(), len(m.connections))
	}
}

// GetConnection gets a connection by websocket conn
func (m *ConnectionManager) GetConnection(conn *websocket.Conn) (*Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connection, exists := m.connections[conn]
	return connection, exists
}

// BroadcastToSubscribers broadcasts a message to all subscribed connections
func (m *ConnectionManager) BroadcastToSubscribers(message string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, connection := range m.connections {
		if connection.session.IsSubscribed() {
			if err := connection.Send(message); err != nil {
				utils.Logger.Warnf("Failed to send to subscriber: %v", err)
			}
		}
	}
}

// GetSubscriberCount returns the number of subscribed connections
func (m *ConnectionManager) GetSubscriberCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, connection := range m.connections {
		if connection.session.IsSubscribed() {
			count++
		}
	}
	return count
}

// TryAcquireEngineLock attempts to acquire the engine lock
func (m *ConnectionManager) TryAcquireEngineLock(conn *websocket.Conn) bool {
	m.engineLock.Lock()
	defer m.engineLock.Unlock()

	if m.engineLocked {
		return false
	}

	m.engineLocked = true
	m.lockHolder = conn

	// Update the session
	m.mu.RLock()
	defer m.mu.RUnlock()
	if connection, exists := m.connections[conn]; exists {
		connection.session.AcquireLock()
	}

	return true
}

// ReleaseEngineLock releases the engine lock
func (m *ConnectionManager) ReleaseEngineLock(conn *websocket.Conn) bool {
	m.engineLock.Lock()
	defer m.engineLock.Unlock()

	if !m.engineLocked || m.lockHolder != conn {
		return false
	}

	m.engineLocked = false
	m.lockHolder = nil

	// Update the session to reflect lock release
	m.mu.RLock()
	if connection, exists := m.connections[conn]; exists {
		connection.session.ReleaseLock()
	}
	m.mu.RUnlock()

	return true
}

// releaseEngineLockInternal releases the engine lock without updating session
// Used internally when session is already being updated or cleaned up
func (m *ConnectionManager) releaseEngineLockInternal(conn *websocket.Conn) bool {
	if !m.engineLocked || m.lockHolder != conn {
		return false
	}

	m.engineLocked = false
	m.lockHolder = nil
	return true
}

// IsEngineLocked checks if the engine is currently locked
func (m *ConnectionManager) IsEngineLocked() bool {
	m.engineLock.Lock()
	defer m.engineLock.Unlock()
	return m.engineLocked
}

// GetConnectionCount returns the total number of connections
func (m *ConnectionManager) GetConnectionCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

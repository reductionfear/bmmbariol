package engine

import (
	"fmt"
	"sync"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// EngineType represents the type of chess engine
type EngineType string

const (
	EngineTypeStockfish EngineType = "stockfish"
	EngineTypeLeela     EngineType = "leela"
	EngineTypeMaia      EngineType = "maia"
	EngineTypeRodent    EngineType = "rodent"
)

// EngineInterface defines the interface all engines must implement
type EngineInterface interface {
	Start() error
	Stop() error
	SendCommand(cmd string) error
	GetOutputChannel() <-chan string
	GetEngineName() string
	IsReady() bool
	Initialize(config *config.Config) error
}

// EngineManager manages multiple chess engines
type EngineManager struct {
	engines      map[EngineType]EngineInterface
	activeEngine EngineInterface
	activeType   EngineType
	config       *config.Config
	mu           sync.RWMutex
}

// NewEngineManager creates a new engine manager
func NewEngineManager(cfg *config.Config) *EngineManager {
	return &EngineManager{
		engines:      make(map[EngineType]EngineInterface),
		activeEngine: nil,
		activeType:   "",
		config:       cfg,
	}
}

// RegisterEngine registers an engine with the manager
func (m *EngineManager) RegisterEngine(engineType EngineType, engine EngineInterface) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.engines[engineType]; exists {
		return fmt.Errorf("engine type %s already registered", engineType)
	}

	m.engines[engineType] = engine
	utils.Logger.Infof("Registered engine: %s", engineType)
	return nil
}

// SetActiveEngine sets the active engine
func (m *EngineManager) SetActiveEngine(engineType EngineType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	engine, exists := m.engines[engineType]
	if !exists {
		return fmt.Errorf("engine type %s not registered", engineType)
	}

	// Stop current active engine if any
	if m.activeEngine != nil {
		utils.Logger.Infof("Stopping current active engine: %s", m.activeType)
		if err := m.activeEngine.Stop(); err != nil {
			utils.Logger.Warnf("Error stopping engine: %v", err)
		}
	}

	// Start new engine
	utils.Logger.Infof("Starting engine: %s", engineType)
	if err := engine.Start(); err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}

	// Initialize engine
	if err := engine.Initialize(m.config); err != nil {
		engine.Stop()
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	m.activeEngine = engine
	m.activeType = engineType
	utils.Logger.Infof("Active engine set to: %s", engineType)
	return nil
}

// GetActiveEngine returns the currently active engine
func (m *EngineManager) GetActiveEngine() (EngineInterface, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeEngine == nil {
		return nil, fmt.Errorf("no active engine")
	}

	return m.activeEngine, nil
}

// GetActiveEngineName returns the name of the active engine
func (m *EngineManager) GetActiveEngineName() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeEngine == nil {
		return "Unknown"
	}

	return m.activeEngine.GetEngineName()
}

// SendCommand sends a command to the active engine
func (m *EngineManager) SendCommand(cmd string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeEngine == nil {
		return fmt.Errorf("no active engine")
	}

	return m.activeEngine.SendCommand(cmd)
}

// GetOutputChannel returns the output channel of the active engine
func (m *EngineManager) GetOutputChannel() (<-chan string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.activeEngine == nil {
		return nil, fmt.Errorf("no active engine")
	}

	return m.activeEngine.GetOutputChannel(), nil
}

// StopAllEngines stops all registered engines
func (m *EngineManager) StopAllEngines() {
	m.mu.Lock()
	defer m.mu.Unlock()

	utils.Logger.Info("Stopping all engines")
	for engineType, engine := range m.engines {
		utils.Logger.Infof("Stopping engine: %s", engineType)
		if err := engine.Stop(); err != nil {
			utils.Logger.Warnf("Error stopping engine %s: %v", engineType, err)
		}
	}

	m.activeEngine = nil
	m.activeType = ""
}

// ListEngines returns a list of registered engine types
func (m *EngineManager) ListEngines() []EngineType {
	m.mu.RLock()
	defer m.mu.RUnlock()

	types := make([]EngineType, 0, len(m.engines))
	for engineType := range m.engines {
		types = append(types, engineType)
	}
	return types
}

// GetActiveType returns the active engine type
func (m *EngineManager) GetActiveType() EngineType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeType
}

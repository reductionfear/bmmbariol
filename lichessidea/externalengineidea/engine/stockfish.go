package engine

import (
	"fmt"
	"path/filepath"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// StockfishEngine represents a Stockfish chess engine
type StockfishEngine struct {
	*UCIEngine
	config *config.Config
}

// NewStockfishEngine creates a new Stockfish engine
func NewStockfishEngine(path string) *StockfishEngine {
	return &StockfishEngine{
		UCIEngine: NewUCIEngine(path),
	}
}

// Initialize initializes the Stockfish engine with configuration
func (e *StockfishEngine) Initialize(cfg *config.Config) error {
	e.config = cfg
	utils.Logger.Info("Initializing Stockfish engine")

	// Send UCI command
	if err := e.SendCommand("uci"); err != nil {
		return err
	}

	// Set threads
	if cfg.Threads > 0 {
		if err := e.SetOption("Threads", fmt.Sprintf("%d", cfg.Threads)); err != nil {
			utils.Logger.Warnf("Failed to set Threads: %v", err)
		}
	}

	// Set hash size
	if cfg.Hash > 0 {
		if err := e.SetOption("Hash", fmt.Sprintf("%d", cfg.Hash)); err != nil {
			utils.Logger.Warnf("Failed to set Hash: %v", err)
		}
	}

	// Set MultiPV for analysis
	if cfg.MultiPV > 0 {
		if err := e.SetOption("MultiPV", fmt.Sprintf("%d", cfg.MultiPV)); err != nil {
			utils.Logger.Warnf("Failed to set MultiPV: %v", err)
		}
	}

	// Configure Syzygy tablebase if available
	if cfg.TablesbaseDir != "" {
		tbPath := filepath.Join(cfg.TablesbaseDir, "syzygy")
		if err := e.SetOption("SyzygyPath", tbPath); err != nil {
			utils.Logger.Warnf("Failed to set SyzygyPath: %v", err)
		} else {
			utils.Logger.Infof("Configured Syzygy tablebase: %s", tbPath)
		}
	}

	// Apply UCI arguments
	for _, arg := range cfg.UCIArgs {
		if arg != "" {
			if err := e.SendCommand(arg); err != nil {
				utils.Logger.Warnf("Failed to send UCI arg: %s", arg)
			}
		}
	}

	// Send isready
	if err := e.SendCommand("isready"); err != nil {
		return err
	}

	utils.Logger.Info("Stockfish engine initialized successfully")
	return nil
}

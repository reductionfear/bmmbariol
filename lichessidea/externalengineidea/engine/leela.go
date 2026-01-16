package engine

import (
	"fmt"
	"path/filepath"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// LeelaEngine represents a Leela Chess Zero (lc0) engine
type LeelaEngine struct {
	*UCIEngine
	config      *config.Config
	weightsFile string
}

// NewLeelaEngine creates a new Leela Chess Zero engine
func NewLeelaEngine(path, weightsFile string) *LeelaEngine {
	return &LeelaEngine{
		UCIEngine:   NewUCIEngine(path),
		weightsFile: weightsFile,
	}
}

// Initialize initializes the Leela engine with configuration
func (e *LeelaEngine) Initialize(cfg *config.Config) error {
	e.config = cfg
	utils.Logger.Info("Initializing Leela Chess Zero engine")

	// Send UCI command
	if err := e.SendCommand("uci"); err != nil {
		return err
	}

	// Set weights file (required for Leela)
	weightsPath := e.weightsFile
	if weightsPath == "" && cfg.WeightsDir != "" {
		// Try to find a default weights file
		weightsPath = filepath.Join(cfg.WeightsDir, "network.pb.gz")
	}

	if weightsPath != "" {
		if err := e.SetOption("WeightsFile", weightsPath); err != nil {
			return fmt.Errorf("failed to set weights file: %w", err)
		}
		utils.Logger.Infof("Leela weights file: %s", weightsPath)
	}

	// Set threads (Leela uses Threads option)
	if cfg.Threads > 0 {
		if err := e.SetOption("Threads", fmt.Sprintf("%d", cfg.Threads)); err != nil {
			utils.Logger.Warnf("Failed to set Threads: %v", err)
		}
	}

	// Set backend (cuda, cudnn, opencl, cpu)
	// Default to auto-detect if not specified
	if err := e.SetOption("Backend", "multiplexing"); err != nil {
		utils.Logger.Warnf("Failed to set Backend: %v", err)
	}

	// Set MultiPV for analysis
	if cfg.MultiPV > 0 {
		if err := e.SetOption("MultiPV", fmt.Sprintf("%d", cfg.MultiPV)); err != nil {
			utils.Logger.Warnf("Failed to set MultiPV: %v", err)
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

	utils.Logger.Info("Leela Chess Zero engine initialized successfully")
	return nil
}

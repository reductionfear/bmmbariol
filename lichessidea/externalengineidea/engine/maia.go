package engine

import (
	"fmt"
	"path/filepath"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// MaiaEngine represents a Maia chess engine (uses lc0 with Maia weights)
type MaiaEngine struct {
	*UCIEngine
	config      *config.Config
	weightsFile string
	rating      int // Maia rating level (1100-1900)
}

// NewMaiaEngine creates a new Maia engine
func NewMaiaEngine(path, weightsFile string, rating int) *MaiaEngine {
	return &MaiaEngine{
		UCIEngine:   NewUCIEngine(path),
		weightsFile: weightsFile,
		rating:      rating,
	}
}

// Initialize initializes the Maia engine with configuration
func (e *MaiaEngine) Initialize(cfg *config.Config) error {
	e.config = cfg
	utils.Logger.Infof("Initializing Maia engine (rating: %d)", e.rating)

	// Send UCI command
	if err := e.SendCommand("uci"); err != nil {
		return err
	}

	// Set Maia weights file
	weightsPath := e.weightsFile
	if weightsPath == "" && cfg.WeightsDir != "" {
		// Use default Maia weights based on rating
		weightsPath = filepath.Join(cfg.WeightsDir, fmt.Sprintf("maia-%d.pb.gz", e.rating))
	}

	if weightsPath != "" {
		if err := e.SetOption("WeightsFile", weightsPath); err != nil {
			return fmt.Errorf("failed to set Maia weights file: %w", err)
		}
		utils.Logger.Infof("Maia weights file: %s", weightsPath)
	}

	// Maia-specific settings for human-like play
	// Single thread for deterministic behavior
	if err := e.SetOption("Threads", "1"); err != nil {
		utils.Logger.Warnf("Failed to set Threads: %v", err)
	}

	// Minibatch size
	if err := e.SetOption("MinibatchSize", "1"); err != nil {
		utils.Logger.Warnf("Failed to set MinibatchSize: %v", err)
	}

	// Max prefetch
	if err := e.SetOption("MaxPrefetch", "0"); err != nil {
		utils.Logger.Warnf("Failed to set MaxPrefetch: %v", err)
	}

	// Nodes per second limit for human-like timing
	// Maia uses very low NPS to simulate human thinking time
	if err := e.SetOption("NodesPerSecondLimit", "0.001"); err != nil {
		utils.Logger.Warnf("Failed to set NodesPerSecondLimit: %v", err)
	}

	// SlowMover setting
	if err := e.SetOption("SlowMover", "0"); err != nil {
		utils.Logger.Warnf("Failed to set SlowMover: %v", err)
	}

	// Backend
	if err := e.SetOption("Backend", "multiplexing"); err != nil {
		utils.Logger.Warnf("Failed to set Backend: %v", err)
	}

	// MultiPV for analysis (Maia typically uses 1)
	multiPV := 1
	if cfg.MultiPV > 0 {
		multiPV = cfg.MultiPV
	}
	if err := e.SetOption("MultiPV", fmt.Sprintf("%d", multiPV)); err != nil {
		utils.Logger.Warnf("Failed to set MultiPV: %v", err)
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

	utils.Logger.Info("Maia engine initialized successfully")
	return nil
}

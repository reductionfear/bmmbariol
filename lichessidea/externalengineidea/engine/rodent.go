package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// RodentEngine represents a Rodent IV chess engine with personality support
type RodentEngine struct {
	*UCIEngine
	config          *config.Config
	personalityFile string
}

// NewRodentEngine creates a new Rodent engine
func NewRodentEngine(path, personalityFile string) *RodentEngine {
	return &RodentEngine{
		UCIEngine:       NewUCIEngine(path),
		personalityFile: personalityFile,
	}
}

// Initialize initializes the Rodent engine with configuration and personality
func (e *RodentEngine) Initialize(cfg *config.Config) error {
	e.config = cfg
	utils.Logger.Info("Initializing Rodent IV engine")

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

	// Set hash
	if cfg.Hash > 0 {
		if err := e.SetOption("Hash", fmt.Sprintf("%d", cfg.Hash)); err != nil {
			utils.Logger.Warnf("Failed to set Hash: %v", err)
		}
	}

	// Set MultiPV
	if cfg.MultiPV > 0 {
		if err := e.SetOption("MultiPV", fmt.Sprintf("%d", cfg.MultiPV)); err != nil {
			utils.Logger.Warnf("Failed to set MultiPV: %v", err)
		}
	}

	// Load personality file if specified
	if e.personalityFile != "" {
		if err := e.loadPersonality(e.personalityFile); err != nil {
			utils.Logger.Warnf("Failed to load personality: %v", err)
		} else {
			utils.Logger.Infof("Loaded Rodent personality: %s", e.personalityFile)
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

	utils.Logger.Info("Rodent IV engine initialized successfully")
	return nil
}

// loadPersonality loads a personality file and applies its settings
func (e *RodentEngine) loadPersonality(filename string) error {
	// Construct full path
	personalityPath := filename
	if !filepath.IsAbs(filename) && e.config.PersonalitiesDir != "" {
		personalityPath = filepath.Join(e.config.PersonalitiesDir, filename)
	}

	// Read personality file
	data, err := os.ReadFile(personalityPath)
	if err != nil {
		return fmt.Errorf("failed to read personality file: %w", err)
	}

	// Parse and apply personality settings
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Process setoption commands
		if strings.HasPrefix(strings.ToLower(line), "setoption name ") {
			// Send the command directly to the engine
			if err := e.SendCommand(line); err != nil {
				utils.Logger.Warnf("Failed to apply personality setting: %s - %v", line, err)
			}
		}
	}

	return nil
}

// SetPersonality changes the personality at runtime
func (e *RodentEngine) SetPersonality(filename string) error {
	utils.Logger.Infof("Changing Rodent personality to: %s", filename)
	e.personalityFile = filename
	return e.loadPersonality(filename)
}

package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/intelligence"
)

// Config holds the server configuration
type Config struct {
	// Server settings
	Address         string `json:"address"`
	TLS             bool   `json:"tls"`
	CertFile        string `json:"cert_file"`
	KeyFile         string `json:"key_file"`

	// Authentication
	AuthWrite       bool   `json:"auth_write"`
	AuthRead        bool   `json:"auth_read"`
	LocalhostBypass bool   `json:"localhost_bypass"`
	PassKey         string `json:"pass_key"` // Will be generated if empty

	// Engine settings
	EnginePath      string   `json:"engine_path"`
	UCIArgs         []string `json:"uci_args"`
	MultiPV         int      `json:"multipv"`
	Threads         int      `json:"threads"`
	Hash            int      `json:"hash"` // MB

	// Intelligence
	IntelligenceEnabled bool                               `json:"intelligence_enabled"`
	IntelligenceSettings *intelligence.IntelligenceSettings `json:"intelligence_settings"`

	// Paths
	EnginesDir       string `json:"engines_dir"`
	BooksDir         string `json:"books_dir"`
	WeightsDir       string `json:"weights_dir"`
	PersonalitiesDir string `json:"personalities_dir"`
	TablesbaseDir    string `json:"tablesbase_dir"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Address:              "localhost:8080",
		TLS:                  false,
		CertFile:             "",
		KeyFile:              "",
		AuthWrite:            true,
		AuthRead:             false,
		LocalhostBypass:      true,
		PassKey:              "",
		EnginePath:           "./stockfish",
		UCIArgs:              []string{},
		MultiPV:              3,
		Threads:              1,
		Hash:                 128,
		IntelligenceEnabled:  false,
		IntelligenceSettings: intelligence.NewDefaultIntelligenceSettings(),
		EnginesDir:           "./engines",
		BooksDir:             "./books",
		WeightsDir:           "./weights",
		PersonalitiesDir:     "./personalities",
		TablesbaseDir:        "./ctg",
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	config := DefaultConfig()

	if path == "" {
		return config, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(path, data, 0644)
}

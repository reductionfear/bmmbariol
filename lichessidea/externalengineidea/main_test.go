package main

import (
	"testing"
	"time"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

func TestLoadConfiguration(t *testing.T) {
	// Initialize logger for tests
	if err := utils.InitLogger(true); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger()

	// Test loading default config
	cfg, err := config.LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	if cfg.Address != "localhost:8080" {
		t.Errorf("Expected default address localhost:8080, got %s", cfg.Address)
	}

	if cfg.MultiPV != 3 {
		t.Errorf("Expected default MultiPV 3, got %d", cfg.MultiPV)
	}
}

func TestPrepareUCIArgs(t *testing.T) {
	cfg := &config.Config{
		MultiPV:  5,
		Threads:  4,
		Hash:     256,
		UCIArgs:  []string{"setoption name Contempt value 0"},
	}

	args := prepareUCIArgs(cfg)

	expectedArgs := []string{
		"setoption name MultiPV value 5",
		"setoption name Threads value 4",
		"setoption name Hash value 256",
		"setoption name Contempt value 0",
	}

	if len(args) != len(expectedArgs) {
		t.Fatalf("Expected %d args, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Arg %d: expected %s, got %s", i, expected, args[i])
		}
	}
}

func TestPrepareUCIArgsMinimal(t *testing.T) {
	cfg := &config.Config{
		MultiPV: 1, // Won't be included (default is 1)
		Threads: 1, // Won't be included (default is 1)
		Hash:    0,  // Won't be included (0 means not set)
		UCIArgs: []string{},
	}

	args := prepareUCIArgs(cfg)

	if len(args) != 0 {
		t.Errorf("Expected 0 args for minimal config, got %d: %v", len(args), args)
	}
}

func TestGracefulShutdownSetup(t *testing.T) {
	// Just verify the main package compiles and the types are correct
	// We can't fully test main() without starting a server, but we can
	// test that the components initialize correctly

	if err := utils.InitLogger(true); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.CloseLogger()

	// Test that we can create config
	cfg := config.DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	// Test that we can prepare UCI args
	args := prepareUCIArgs(cfg)
	if args == nil {
		t.Fatal("prepareUCIArgs returned nil")
	}
}

func TestBridgeChannels(t *testing.T) {
	// Test that channels work as expected
	engineInputChan := make(chan string, 10)

	// Send a test command
	select {
	case engineInputChan <- "test command":
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Failed to send to engineInputChan")
	}

	// Receive the test command
	select {
	case cmd := <-engineInputChan:
		if cmd != "test command" {
			t.Errorf("Expected 'test command', got '%s'", cmd)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Failed to receive from engineInputChan")
	}
}

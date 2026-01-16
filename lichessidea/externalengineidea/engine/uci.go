package engine

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

// UCIEngine represents a UCI chess engine
type UCIEngine struct {
	path            string
	cmd             *exec.Cmd
	stdin           *bufio.Writer
	stdout          *bufio.Scanner
	outputChan      chan string
	inputChan       chan string
	engineName      string
	isReady         bool
	mu              sync.RWMutex
	stopChan        chan struct{}
	stoppedChan     chan struct{}
}

// NewUCIEngine creates a new UCI engine instance
func NewUCIEngine(path string) *UCIEngine {
	return &UCIEngine{
		path:        path,
		outputChan:  make(chan string, 100),
		inputChan:   make(chan string, 10),
		isReady:     false,
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}
}

// Start starts the UCI engine
func (e *UCIEngine) Start() error {
	utils.Logger.Infof("Starting UCI engine: %s", e.path)

	// Start engine process
	e.cmd = exec.Command(e.path)

	// Setup stdin
	stdin, err := e.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	e.stdin = bufio.NewWriter(stdin)

	// Setup stdout
	stdout, err := e.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	e.stdout = bufio.NewScanner(stdout)

	// Start the process
	if err := e.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}

	// Start output reader goroutine
	go e.readOutput()

	// Start input writer goroutine
	go e.writeInput()

	utils.Logger.Info("UCI engine started successfully")
	return nil
}

// Stop stops the UCI engine
func (e *UCIEngine) Stop() error {
	utils.Logger.Info("Stopping UCI engine")

	// Send quit command
	e.SendCommand("quit")

	// Signal stop
	close(e.stopChan)

	// Wait for goroutines to finish
	<-e.stoppedChan

	// Wait for process to exit (with timeout would be better in production)
	if e.cmd != nil && e.cmd.Process != nil {
		e.cmd.Wait()
	}

	utils.Logger.Info("UCI engine stopped")
	return nil
}

// SendCommand sends a command to the engine
func (e *UCIEngine) SendCommand(cmd string) error {
	select {
	case e.inputChan <- cmd:
		return nil
	default:
		return fmt.Errorf("input channel full")
	}
}

// GetOutputChannel returns the output channel
func (e *UCIEngine) GetOutputChannel() <-chan string {
	return e.outputChan
}

// GetEngineName returns the engine name
func (e *UCIEngine) GetEngineName() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.engineName
}

// IsReady returns whether the engine is ready
func (e *UCIEngine) IsReady() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isReady
}

// SetOption sets a UCI option
func (e *UCIEngine) SetOption(name, value string) error {
	cmd := fmt.Sprintf("setoption name %s value %s", name, value)
	return e.SendCommand(cmd)
}

// NewGame sends ucinewgame command
func (e *UCIEngine) NewGame() error {
	return e.SendCommand("ucinewgame")
}

// SetPosition sets the position
func (e *UCIEngine) SetPosition(fen string, moves []string) error {
	var cmd string
	if fen == "" || fen == "startpos" {
		cmd = "position startpos"
	} else {
		cmd = fmt.Sprintf("position fen %s", fen)
	}

	if len(moves) > 0 {
		cmd += " moves " + strings.Join(moves, " ")
	}

	return e.SendCommand(cmd)
}

// Go sends the go command with parameters
func (e *UCIEngine) Go(params string) error {
	cmd := "go"
	if params != "" {
		cmd += " " + params
	}
	return e.SendCommand(cmd)
}

// readOutput reads output from the engine
func (e *UCIEngine) readOutput() {
	defer close(e.stoppedChan)

	for e.stdout.Scan() {
		line := e.stdout.Text()
		utils.Logger.Debugf("Engine output: %s", line)

		// Parse engine name
		if strings.HasPrefix(line, "id name ") {
			e.mu.Lock()
			e.engineName = strings.TrimPrefix(line, "id name ")
			e.mu.Unlock()
			utils.Logger.Infof("Engine identified: %s", e.engineName)
		}

		// Parse readyok
		if line == "readyok" {
			e.mu.Lock()
			e.isReady = true
			e.mu.Unlock()
		}

		// Send to output channel
		select {
		case e.outputChan <- line:
		default:
			utils.Logger.Warn("Output channel full, dropping line")
		}

		// Check if we should stop
		select {
		case <-e.stopChan:
			return
		default:
		}
	}

	if err := e.stdout.Err(); err != nil {
		utils.Logger.Errorf("Error reading engine output: %v", err)
	}
}

// writeInput writes input to the engine
func (e *UCIEngine) writeInput() {
	for {
		select {
		case cmd := <-e.inputChan:
			utils.Logger.Debugf("Sending to engine: %s", cmd)
			_, err := e.stdin.WriteString(cmd + "\n")
			if err != nil {
				utils.Logger.Errorf("Error writing to engine: %v", err)
				return
			}
			if err := e.stdin.Flush(); err != nil {
				utils.Logger.Errorf("Error flushing engine input: %v", err)
				return
			}

		case <-e.stopChan:
			return
		}
	}
}

// Initialize initializes the engine with UCI protocol
func (e *UCIEngine) Initialize(uciArgs []string) error {
	utils.Logger.Info("Initializing UCI engine")

	// Send uci command
	if err := e.SendCommand("uci"); err != nil {
		return err
	}

	// Wait for uciok (in practice, we should wait properly)
	// For now, we'll send additional commands after a brief moment
	
	// Apply UCI arguments
	for _, arg := range uciArgs {
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

	utils.Logger.Info("UCI engine initialized")
	return nil
}

package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/config"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/engine"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/server"
	"github.com/reductionfear/bmmbariol/lichessidea/externalengineidea/utils"
)

var (
	// Command-line flags
	addr            = flag.String("addr", "localhost:8080", "HTTP service address")
	enginePath      = flag.String("engine", "./stockfish", "Path to the engine binary")
	uciArgsFlag     = flag.String("uciargs", "", "UCI initialization commands (semicolon-separated)")
	authWrite       = flag.Bool("authwrite", true, "Require authentication for write operations")
	authRead        = flag.Bool("authread", false, "Require authentication for read operations")
	localhostBypass = flag.Bool("localhost", true, "Bypass authentication for localhost connections")
	tls             = flag.Bool("tls", false, "Enable TLS/WSS")
	certFile        = flag.String("cert", "", "TLS certificate file")
	keyFile         = flag.String("key", "", "TLS private key file")
	configFile      = flag.String("config", "", "Configuration JSON file")
	multiPV         = flag.Int("multipv", 3, "Multi-PV count for analysis")
	threads         = flag.Int("threads", 1, "Number of engine threads")
	hash            = flag.Int("hash", 128, "Hash table size in MB")
	intelligence    = flag.Bool("intelligence", false, "Enable intelligence system (experimental)")
	debug           = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Initialize logger
	if err := utils.InitLogger(*debug); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer utils.CloseLogger()

	utils.Logger.Info("=== Lichess External Engine Server ===")
	utils.Logger.Info("Starting server initialization...")

	// Load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		utils.Logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate TLS configuration
	if cfg.TLS {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			utils.Logger.Fatal("TLS enabled but cert/key files not specified")
		}
		if _, err := os.Stat(cfg.CertFile); os.IsNotExist(err) {
			utils.Logger.Fatalf("Certificate file not found: %s", cfg.CertFile)
		}
		if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
			utils.Logger.Fatalf("Key file not found: %s", cfg.KeyFile)
		}
	}

	// Create engine input channel
	engineInputChan := make(chan string, 10)

	// Initialize UCI engine
	utils.Logger.Infof("Initializing engine: %s", cfg.EnginePath)
	uciEngine := engine.NewUCIEngine(cfg.EnginePath)

	if err := uciEngine.Start(); err != nil {
		utils.Logger.Fatalf("Failed to start engine: %v", err)
	}
	defer func() {
		utils.Logger.Info("Shutting down engine...")
		uciEngine.Stop()
	}()

	// Prepare UCI initialization arguments
	uciArgs := prepareUCIArgs(cfg)

	// Initialize engine with UCI commands
	if err := uciEngine.Initialize(uciArgs); err != nil {
		utils.Logger.Fatalf("Failed to initialize engine: %v", err)
	}

	// Wait for engine to be ready (give it a moment to process uci and isready)
	time.Sleep(500 * time.Millisecond)

	// Create authentication manager
	authManager := server.NewAuthManager(cfg.AuthWrite, cfg.AuthRead, cfg.LocalhostBypass, cfg.PassKey)
	utils.Logger.Infof("Authentication passkey: %s", authManager.GetPassKey())
	utils.Logger.Infof("Auth required for writes: %t", cfg.AuthWrite)
	utils.Logger.Infof("Auth required for reads: %t", cfg.AuthRead)
	utils.Logger.Infof("Localhost bypass: %t", cfg.LocalhostBypass)

	// Create WebSocket server
	wsServer := server.NewWebSocketServer(authManager, engineInputChan)

	// Bridge engine output to WebSocket broadcasts
	go bridgeEngineOutput(uciEngine, wsServer)

	// Bridge WebSocket input to engine
	go bridgeEngineInput(uciEngine, engineInputChan)

	// Setup HTTP handlers
	http.HandleFunc("/ws", wsServer.HandleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeTemplate.Execute(w, cfg.Address)
	})

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		protocol := "ws"
		if cfg.TLS {
			protocol = "wss"
		}
		utils.Logger.Infof("Server listening on %s://%s/ws", protocol, cfg.Address)
		utils.Logger.Infof("Home page available at http%s://%s/", 
			map[bool]string{true: "s", false: ""}[cfg.TLS], cfg.Address)

		if cfg.TLS {
			serverErr <- http.ListenAndServeTLS(cfg.Address, cfg.CertFile, cfg.KeyFile, nil)
		} else {
			serverErr <- http.ListenAndServe(cfg.Address, nil)
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-sigChan:
		utils.Logger.Infof("Received signal: %v", sig)
		utils.Logger.Info("Initiating graceful shutdown...")
	case err := <-serverErr:
		if err != nil {
			utils.Logger.Fatalf("Server error: %v", err)
		}
	}

	utils.Logger.Info("Server shutdown complete")
}

// loadConfiguration loads configuration from file or command-line flags
func loadConfiguration() (*config.Config, error) {
	var cfg *config.Config
	var err error

	if *configFile != "" {
		utils.Logger.Infof("Loading configuration from: %s", *configFile)
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	} else {
		utils.Logger.Info("Using default configuration with command-line overrides")
		cfg = config.DefaultConfig()
	}

	// Override with command-line flags (flags take precedence)
	if flag.Lookup("addr").Value.String() != flag.Lookup("addr").DefValue {
		cfg.Address = *addr
	}
	if flag.Lookup("engine").Value.String() != flag.Lookup("engine").DefValue {
		cfg.EnginePath = *enginePath
	}
	if flag.Lookup("authwrite").Value.String() != flag.Lookup("authwrite").DefValue {
		cfg.AuthWrite = *authWrite
	}
	if flag.Lookup("authread").Value.String() != flag.Lookup("authread").DefValue {
		cfg.AuthRead = *authRead
	}
	if flag.Lookup("localhost").Value.String() != flag.Lookup("localhost").DefValue {
		cfg.LocalhostBypass = *localhostBypass
	}
	if flag.Lookup("tls").Value.String() != flag.Lookup("tls").DefValue {
		cfg.TLS = *tls
	}
	if flag.Lookup("cert").Value.String() != flag.Lookup("cert").DefValue {
		cfg.CertFile = *certFile
	}
	if flag.Lookup("key").Value.String() != flag.Lookup("key").DefValue {
		cfg.KeyFile = *keyFile
	}
	if flag.Lookup("multipv").Value.String() != flag.Lookup("multipv").DefValue {
		cfg.MultiPV = *multiPV
	}
	if flag.Lookup("threads").Value.String() != flag.Lookup("threads").DefValue {
		cfg.Threads = *threads
	}
	if flag.Lookup("hash").Value.String() != flag.Lookup("hash").DefValue {
		cfg.Hash = *hash
	}
	if flag.Lookup("intelligence").Value.String() != flag.Lookup("intelligence").DefValue {
		cfg.IntelligenceEnabled = *intelligence
	}

	// Parse UCI args from flag
	if *uciArgsFlag != "" {
		cfg.UCIArgs = strings.Split(*uciArgsFlag, ";")
	}

	return cfg, nil
}

// prepareUCIArgs prepares UCI initialization arguments from config
func prepareUCIArgs(cfg *config.Config) []string {
	args := make([]string, 0)

	// Add MultiPV option if specified
	if cfg.MultiPV > 1 {
		args = append(args, fmt.Sprintf("setoption name MultiPV value %d", cfg.MultiPV))
	}

	// Add Threads option if specified
	if cfg.Threads > 1 {
		args = append(args, fmt.Sprintf("setoption name Threads value %d", cfg.Threads))
	}

	// Add Hash option if specified
	if cfg.Hash > 0 {
		args = append(args, fmt.Sprintf("setoption name Hash value %d", cfg.Hash))
	}

	// Add user-provided UCI args
	args = append(args, cfg.UCIArgs...)

	return args
}

// bridgeEngineOutput bridges engine output to WebSocket server
func bridgeEngineOutput(uciEngine *engine.UCIEngine, wsServer *server.WebSocketServer) {
	utils.Logger.Info("Starting engine output bridge...")

	outputChan := uciEngine.GetOutputChannel()
	engineNameSet := false

	for output := range outputChan {
		// Set engine name when detected
		if !engineNameSet && uciEngine.GetEngineName() != "" {
			wsServer.SetEngineName(uciEngine.GetEngineName())
			engineNameSet = true
			utils.Logger.Infof("Engine name set: %s", uciEngine.GetEngineName())
		}

		// Broadcast info and bestmove lines to subscribers
		if strings.HasPrefix(output, "info") || strings.HasPrefix(output, "bestmove") {
			wsServer.BroadcastEngineOutput(output)
		}
	}

	utils.Logger.Info("Engine output bridge stopped")
}

// bridgeEngineInput bridges WebSocket input to engine
func bridgeEngineInput(uciEngine *engine.UCIEngine, inputChan <-chan string) {
	utils.Logger.Info("Starting engine input bridge...")

	for cmd := range inputChan {
		if err := uciEngine.SendCommand(cmd); err != nil {
			utils.Logger.Warnf("Failed to send command to engine: %v", err)
		}
	}

	utils.Logger.Info("Engine input bridge stopped")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Chesshook External Engine server</title>
</head>
<body>
<p>This server is running a websocket server for the Chesshook userscript external engine protocol. You can find the chesshook userscript <a href="https://github.com/0mlml/chesshook">here.</a> You can find the source for this server <a href="https://github.com/0mlml/chesshook-intermediary">here.</a></p>
<p>To use this server with the userscript, set the engine url to <code>ws://{{.}}/ws</code></p>
</body>
</html>
`))

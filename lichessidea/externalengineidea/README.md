# Lichess External Engine Server (Go)

A production-ready WebSocket server for integrating chess engines with the Lichess External Engine protocol. Built with Go for high performance and reliability.

## Features

- 🚀 **High Performance**: Built with Go's concurrency primitives for efficient multi-client handling
- 🔒 **Secure**: TLS/WSS support, authentication, and localhost bypass
- 🎯 **UCI Compatible**: Full UCI protocol support for any chess engine
- 🧠 **Intelligence System**: Experimental move modification system (optional)
- 📊 **Production Ready**: Structured logging with zap, graceful shutdown, configurable options
- 🔧 **Flexible Configuration**: Command-line flags and JSON config file support

## Architecture

The server is built with a modular architecture following the specification in `architecture-specsheet.md`:

```
├── main.go                    # Entry point & HTTP server
├── config/                    # Configuration management
│   └── config.go
├── server/                    # WebSocket server components
│   ├── websocket.go          # Main WebSocket handler
│   ├── auth.go               # Authentication logic
│   ├── connection.go         # Connection manager
│   └── message.go            # Message parsing/routing
├── engine/                    # UCI engine management
│   └── uci.go                # UCI protocol implementation
├── intelligence/              # Move intelligence system
│   └── settings.go
├── models/                    # Data models
│   ├── move.go
│   ├── game_state.go
│   └── analysis.go
├── chess/                     # Chess logic
│   └── board.go
└── utils/                     # Utilities
    ├── logger.go             # Structured logging
    └── helpers.go
```

## Installation

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- A UCI-compatible chess engine (e.g., Stockfish)

### Building from Source

```bash
cd lichessidea/externalengineidea
go build -o lichess-engine-server
```

## Usage

### Basic Usage

```bash
./lichess-engine-server -engine ./stockfish
```

The server will:
1. Generate a random authentication passkey (printed to console)
2. Start listening on `ws://localhost:8080/ws`
3. Initialize the chess engine

### Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-addr` | string | `localhost:8080` | HTTP service address |
| `-engine` | string | `./stockfish` | Path to engine binary |
| `-authwrite` | bool | `true` | Require auth for write operations |
| `-authread` | bool | `false` | Require auth for read operations |
| `-localhost` | bool | `true` | Bypass auth for localhost |
| `-tls` | bool | `false` | Enable TLS/WSS |
| `-cert` | string | - | TLS certificate file |
| `-key` | string | - | TLS private key file |
| `-config` | string | - | JSON configuration file |
| `-multipv` | int | `3` | Multi-PV count for analysis |
| `-threads` | int | `1` | Number of engine threads |
| `-hash` | int | `128` | Hash table size in MB |
| `-intelligence` | bool | `false` | Enable intelligence system |
| `-debug` | bool | `false` | Enable debug logging |
| `-uciargs` | string | - | UCI commands (semicolon-separated) |

### Examples

**Basic Stockfish with custom options:**
```bash
./lichess-engine-server -engine ./stockfish -multipv 5 -threads 4 -hash 512
```

**With UCI arguments:**
```bash
./lichess-engine-server -engine ./stockfish -uciargs "setoption name Contempt value 0;setoption name Skill Level value 20"
```

**Secure WebSocket (WSS):**
```bash
./lichess-engine-server -tls -cert server.crt -key server.key -addr :8443
```

**Using config file:**
```bash
./lichess-engine-server -config config.json
```

### Configuration File

Create a `config.json` file:

```json
{
  "address": "localhost:8080",
  "tls": false,
  "auth_write": true,
  "auth_read": false,
  "localhost_bypass": true,
  "engine_path": "./stockfish",
  "multipv": 5,
  "threads": 4,
  "hash": 512,
  "uci_args": [
    "setoption name Contempt value 0"
  ],
  "intelligence_enabled": false
}
```

## Protocol Documentation

### Version Checking

```
message: whoareyou 

reponse: iam <name>v<version>

message: whatengine

response: engine <name>

Useful for checking if the server is outdated compared to the userscript.
``` 
### Authentication
```
message: `auth <passkey>`
response: `authok` or `autherr`
If the server responds with `authok`, the connection has been authenticated.
If the server responds with `autherr`, the connection has not been authenticated. The passkey could be incorrect or the server could be stonewalling the connection.

If the user fails to authenticate three times, every attempt will continue to fail.
```
### Subscribing to engine output
```
message: `sub`
response: `subok` or `suberr` or `autherr`
If the server responds with `subok`, the client will recieve all engine output beginning with `bestmove` and `info` through the websocket.
If the server responds with `suberr`, the client is already subscribed
If the server responds with `autherr`, the client is expected to provide a passkey, as it is required for read access.  

message: `unsub`
response: `unsubok` or `unsuberr` or `autherr`
If the server responds with `unsubok`, the client will no longer recieve engine output.
If the server responds with `unsuberr`, the client is not subscribed
If the server responds with `autherr`, the client is expected to provide a passkey, as it is required for read access.  
```
### Locking the engine for others
```
message: `lock`
response: `lockok` or `lockerr` or `autherr`
If the server responds with `lockok`, the server will reject all uci commands from other connections.
If the server responds with `lockerr`, someone already has a lock on the server. the client is expected to maintain whether it is the one with the lock.
If the server responds with `autherr`, the client is expected to provide a passkey, as it is required for write access.  

message: `unlock`
response: `unlockok` or `unlockerr` or `autherr`
If the server responds with `unlockok`, the server will allow other connections to lock the server or send uci commands.
If the server responds with `unlockerr`, the client does not have lock on the server.
If the server responds with `autherr`, the client is expected to provide a passkey, as it is required for write access.  
```
### UCI Commands
```
message: `<uci command>`
response: `autherr`
If the server responds with `autherr`, the client is expected to provide a passkey, as it is required for write access.  

The client is expected to subscribe to engine output if it would like to recieve engine output.
```

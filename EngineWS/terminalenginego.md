# Engine Package & Intelligence Package - Deep Dive

## What Each File Does

### `/engine/` Package

```
engine/
├── manager.go      # Orchestra conductor for all engines
├── uci.go          # Universal translator (UCI protocol)
├── stockfish.go    # Stockfish specialist
├── leela.go        # Lc0/Neural net specialist
├── maia.go         # Human-like play specialist
└── rodent.go       # Personality engine specialist
```

---

### **manager.go** - Multi-Engine Pool Manager

**Purpose:** Controls multiple engines, switches between them, handles lifecycle

```go
// What it does:
type EngineManager struct {
    engines     map[string]UCIEngine  // All loaded engines
    activeID    string                // Currently active engine
    mutex       sync.RWMutex          // Thread safety
}

// Key functions:
- LoadEngine(config) → Start engine process, initialize UCI
- SwitchEngine(id)   → Change active engine on-the-fly
- GetActive()        → Return current engine
- ShutdownAll()      → Clean termination of all engines
- ListEngines()      → Return available engine IDs
```

---

### **uci.go** - UCI Protocol Implementation

**Purpose:** Speaks the Universal Chess Interface language with any engine

```go
// What it does: 
type UCIEngine interface {
    Start() error
    SendCommand(cmd string) error
    ReadOutput() <-chan string
    SetOption(name, value string) error
    SetPosition(fen string, moves []string) error
    Go(params GoParams) error
    Stop() error
    Quit() error
    IsReady() bool
    GetName() string
}

// Handles UCI commands:
- "uci"           → Initialize engine
- "isready"       → Check engine ready
- "setoption"     → Configure engine
- "position"      → Set board position
- "go"            → Start analysis
- "stop"          → Stop analysis
- "quit"          → Terminate engine
```

---

### **stockfish.go** - Stockfish-Specific Handler

**Purpose:** Knows Stockfish's unique options and behaviors

```go
// Stockfish-specific features:
- Skill Level limiting (0-20)
- UCI_Elo limiting (1320-3190)
- NNUE evaluation toggle
- Syzygy tablebase integration
- Contempt settings
- Analysis Contempt modes
```

---

### **leela.go** - Leela Chess Zero Handler

**Purpose:** Handles neural network engine specifics

```go
// Lc0-specific features:
- Backend selection (cuda, cudnn, opencl, eigen)
- Weights file loading
- NNCache configuration
- Minibatch/Prefetch tuning
- Temperature settings
- WDL (Win/Draw/Loss) scoring
```

---

### **maia.go** - Maia Engine Handler

**Purpose:** Human-like play configuration

```go
// Maia-specific features:
- Strength level via weight file selection
- Node limiting for human-like speed
- Temperature for move variety
- Disabled search features for authenticity
```

---

### **rodent.go** - Rodent Engine Handler

**Purpose:** Personality-based chess play

```go
// Rodent-specific features:
- Personality file loading
- Dynamic weight adjustments
- Style parameters (aggressive, defensive, etc.)
- Eval blur for variety
```

---

### `/intelligence/` Package

```
intelligence/
├── settings.go     # Configuration data structure
├── evaluator.go    # Modifies engine evaluations
├── analyzer.go     # Analyzes position characteristics
└── multipliers.go  # Applies preference multipliers
```

---

### **settings.go** - Intelligence Settings Model

**Purpose:** Data structure for all intelligence preferences

```go
type IntelligenceSettings struct {
    Enabled                  bool
    AvoidLowIntelligence     bool
    LowIntelligenceThreshold float64
    // ... all multipliers and preferences
}
```

---

### **evaluator.go** - Move Evaluation Modifier

**Purpose:** Takes engine moves and adjusts scores based on preferences

```go
// What it does:
- Receives candidate moves from engine
- Applies multipliers based on move characteristics
- Re-ranks moves according to personality
- Handles checkmate detection (absolute 1000. 0 score)
- Protects critical positions from bad modifications
```

---

### **analyzer.go** - Position Analysis

**Purpose:** Understands what's happening on the board

```go
// What it analyzes:
- Is position critical?  (check, forced moves)
- Game phase (opening, middlegame, endgame)
- Piece activity and mobility
- King safety status
- Tactical motifs present
```

---

### **multipliers.go** - Piece/Move Multipliers

**Purpose:** Applies the actual score modifications

```go
// Multiplier logic:
func ApplyMultiplier(eval, mult float64, isCritical bool) float64 {
    if isCritical {
        mult = 1.0 + (mult-1.0)*0.3  // Reduce effect
    }
    if eval > 0 {
        return eval * mult
    }
    return eval / mult  // Inverse for negative evals
}
```

---

## 🚀 Direct Terminal Commands for Rodent & Maia

Yes!  You can use these directly without config files:

### Rodent with Personality (Direct Command)

```bash
# Windows
./lichess-engine-server.exe \
  -engine "./engines/rodent/RodentIV.exe" \
  -personality "./personalities/fighter. txt"

# Or with inline UCI args
./lichess-engine-server.exe \
  -engine "./engines/rodent/RodentIV.exe" \
  -uciargs "setoption name PersonalityFile value ./personalities/fighter.txt;setoption name Contempt value 50;setoption name Mobility Weight value 130"
```

### Maia with Strength Level (Direct Command)

```bash
# Maia 1100 (Beginner)
./lichess-engine-server.exe \
  -engine "./engines/lc0/lc0.exe" \
  -weights "./engines/maia/maia-1100.pb. gz" \
  -backend "eigen"

# Maia 1500 (Club Player)
./lichess-engine-server.exe \
  -engine "./engines/lc0/lc0.exe" \
  -weights "./engines/maia/maia-1500.pb.gz" \
  -backend "eigen"

# Maia 1900 (Expert)
./lichess-engine-server. exe \
  -engine "./engines/lc0/lc0.exe" \
  -weights "./engines/maia/maia-1900.pb.gz" \
  -backend "eigen"
```

---

## Complete Terminal Flag Reference

```bash
./lichess-engine-server.exe [FLAGS]

# Server Flags
  -addr string          Server address (default "localhost: 8080")
  -tls                  Enable TLS/WSS
  -cert string          TLS certificate file
  -key string           TLS private key file

# Authentication Flags
  -authwrite            Require auth for writes (default true)
  -authread             Require auth for reads (default false)
  -localhost            Bypass auth for localhost (default true)
  -passkey string       Custom passkey (default:  random)

# Engine Flags
  -engine string        Engine executable path (default "./stockfish")
  -type string          Engine type:  stockfish|lc0|maia|rodent|generic

# UCI Options (Common)
  -threads int          CPU threads (default 1)
  -hash int             Hash table MB (default 128)
  -multipv int          Multi-PV lines (default 3)
  -uciargs string       Additional UCI commands (semicolon-separated)

# Lc0/Maia Specific
  -weights string       Neural network weights file
  -backend string       Lc0 backend:  cuda-fp16|cuda|cudnn|opencl|eigen

# Rodent Specific
  -personality string   Personality file path

# Stockfish Specific
  -skill int            Skill Level 0-20 (default 20)
  -elo int              UCI_Elo when limited (1320-3190)
  -limiteloelo          Enable ELO limiting

# Book & Tablebase
  -book string          Opening book path
  -booktype string      Book type: polyglot|ctg|abk
  -syzygy string        Syzygy tablebase path

# Intelligence System
  -intelligence         Enable intelligence system
  -aggressive float     Aggressiveness multiplier (default 1.0)
  -passive float        Passiveness multiplier (default 1.0)
  -capture float        Capture preference (default 1.0)
  -castle float         Castle preference (default 1.0)

# Logging
  -debug                Enable debug logging
  -logfile string       Log file path
```

---

## Quick Examples

### 1. Stockfish Full Strength
```bash
./lichess-engine-server.exe -engine ./stockfish. exe -threads 8 -hash 1024 -multipv 5
```

### 2. Stockfish Limited to 1500 ELO
```bash
./lichess-engine-server.exe -engine ./stockfish.exe -limitelo -elo 1500
```

### 3. Stockfish Skill Level 10
```bash
./lichess-engine-server.exe -engine ./stockfish.exe -skill 10
```

### 4. Maia 1500 (Human-like Club Player)
```bash
./lichess-engine-server.exe \
  -engine ./lc0.exe \
  -type maia \
  -weights ./maia-1500.pb.gz \
  -backend eigen
```

### 5. Maia 1900 (Human-like Expert)
```bash
./lichess-engine-server.exe \
  -engine ./lc0.exe \
  -type maia \
  -weights ./maia-1900.pb.gz \
  -backend eigen
```

### 6. Rodent Fighter Personality
```bash
./lichess-engine-server.exe \
  -engine ./RodentIV. exe \
  -type rodent \
  -personality ./personalities/fighter.txt
```

### 7. Rodent Defensive Personality
```bash
./lichess-engine-server.exe \
  -engine ./RodentIV. exe \
  -type rodent \
  -personality ./personalities/defender.txt
```

### 8. Rodent with Inline Personality Settings
```bash
./lichess-engine-server.exe \
  -engine ./RodentIV. exe \
  -uciargs "setoption name Contempt value 100;setoption name Mobility Weight value 140;setoption name King Safety Weight value 80;setoption name Own Attack value 150"
```

### 9. Leela Chess Zero (GPU)
```bash
./lichess-engine-server.exe \
  -engine ./lc0.exe \
  -type lc0 \
  -weights ./BT4-1740.pb.gz \
  -backend cuda-fp16
```

### 10. Stockfish with Book and Tablebase
```bash
./lichess-engine-server.exe \
  -engine ./stockfish.exe \
  -threads 8 \
  -hash 2048 \
  -book ./performance. bin \
  -booktype polyglot \
  -syzygy ./tablebases/syzygy
```

### 11. Aggressive Personality via Intelligence
```bash
./lichess-engine-server.exe \
  -engine ./stockfish.exe \
  -intelligence \
  -aggressive 1.8 \
  -passive 0.5 \
  -capture 1.5 \
  -castle 0.7
```

### 12. Remote Secure Server
```bash
./lichess-engine-server. exe \
  -addr 0.0.0.0:8443 \
  -tls \
  -cert ./server.crt \
  -key ./server.key \
  -engine ./stockfish.exe \
  -threads 16 \
  -authwrite \
  -authread \
  -passkey "my_secret_key"
```

---

## Personality Files Quick Reference

### Create `./personalities/fighter.txt`:
```ini
; Aggressive Fighter
Contempt 50
Mobility Weight 130
King Safety Weight 80
Own Attack 150
Opp Attack 70
Eval Blur 0
```

### Create `./personalities/defender. txt`:
```ini
; Solid Defender
Contempt -10
Mobility Weight 90
King Safety Weight 140
Own Attack 80
Opp Attack 120
Pawn Structure Weight 120
```

### Create `./personalities/gambit.txt`:
```ini
; Gambit King
Contempt 100
Mobility Weight 140
King Safety Weight 60
Own Attack 160
Material Weight 70
Eval Blur 10
```

### Create `./personalities/positional.txt`:
```ini
; Positional Master
Contempt 0
Mobility Weight 110
Pawn Structure Weight 140
King Safety Weight 100
Positional Weight 140
Own Attack 90
```

---

## Summary Table

| Engine | Strength Control | Command |
|--------|------------------|---------|
| **Stockfish** | Skill Level | `-skill 10` |
| **Stockfish** | ELO Limit | `-limitelo -elo 1500` |
| **Maia** | Weight File | `-weights maia-1500.pb.gz` |
| **Rodent** | Personality File | `-personality fighter.txt` |
| **Rodent** | Inline UCI | `-uciargs "setoption name..."` |
| **Lc0** | Weight File | `-weights BT4.pb.gz` |
| **Any** | Intelligence | `-intelligence -aggressive 1.5` |

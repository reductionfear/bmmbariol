# 📘 INSTRUCTIONS. md
## Lichess External Engine WebSocket Server - Complete Configuration Guide

---

## Table of Contents

1. [Quick Start](#1-quick-start)
2. [Configuration File Reference](#2-configuration-file-reference)
3. [Engine-Specific Configurations](#3-engine-specific-configurations)
   - [Stockfish](#31-stockfish)
   - [Leela Chess Zero (Lc0)](#32-leela-chess-zero-lc0)
   - [Komodo](#33-komodo)
   - [Houdini](#34-houdini)
   - [Fritz](#35-fritz)
   - [Rybka](#36-rybka)
   - [Shredder](#37-shredder)
   - [Rodent IV](#38-rodent-iv)
   - [Maia Chess](#39-maia-chess)
   - [Ethereal](#310-ethereal)
   - [Berserk](#311-berserk)
   - [Koivisto](#312-koivisto)
   - [RubiChess](#313-rubichess)
   - [Igel](#314-igel)
   - [Arasan](#315-arasan)
   - [Crafty](#316-crafty)
   - [Toga II](#317-toga-ii)
   - [Fruit](#318-fruit)
   - [Spike](#319-spike)
   - [SugaR](#320-sugar)
   - [Brainfish](#321-brainfish)
   - [Fat Fritz](#322-fat-fritz)
   - [Dragon by Komodo](#323-dragon-by-komodo)
4. [Opening Books Configuration](#4-opening-books-configuration)
5. [Endgame Tablebases](#5-endgame-tablebases)
6. [Personalities & Playing Styles](#6-personalities--playing-styles)
7. [Intelligence System](#7-intelligence-system)
8. [Multi-Engine Setup](#8-multi-engine-setup)
9. [Network & Security](#9-network--security)
10. [Troubleshooting](#10-troubleshooting)
11. [Complete Example Configurations](#11-complete-example-configurations)

---

## 1. Quick Start

### Basic Usage

```bash
# Windows
./lichess-engine-server.exe -config config.json

# Linux/macOS
./lichess-engine-server -config config.json
```

### Minimal Configuration

Create `config.json`:

```json
{
  "address": "localhost:8080",
  "engine_path": "./engines/stockfish.exe"
}
```

### Command Line Flags (Override Config)

```bash
./lichess-engine-server. exe \
  -config config.json \
  -addr localhost:9000 \
  -engine ./stockfish.exe \
  -multipv 5 \
  -threads 8
```

| Flag | Description | Default |
|------|-------------|---------|
| `-config` | Configuration file path | `config.json` |
| `-addr` | Server address | `localhost:8080` |
| `-engine` | Engine executable path | `./stockfish` |
| `-multipv` | Multi-PV lines | `3` |
| `-threads` | Engine threads | `1` |
| `-hash` | Hash table size (MB) | `128` |
| `-authwrite` | Require auth for writes | `true` |
| `-authread` | Require auth for reads | `false` |
| `-localhost` | Bypass auth for localhost | `true` |
| `-tls` | Enable TLS/WSS | `false` |
| `-cert` | TLS certificate file | `` |
| `-key` | TLS private key file | `` |
| `-intelligence` | Enable intelligence system | `false` |
| `-debug` | Enable debug logging | `false` |

---

## 2. Configuration File Reference

### Complete Schema

```json
{
  "_comment": "Lichess External Engine Server Configuration",
  "_version": "1.0.0",

  "server": {
    "address": "localhost: 8080",
    "tls": false,
    "cert_file": "",
    "key_file":  "",
    "read_timeout": 30,
    "write_timeout": 30,
    "max_connections": 100
  },

  "auth": {
    "require_for_write": true,
    "require_for_read": false,
    "localhost_bypass": true,
    "passkey": "",
    "max_failed_attempts": 3
  },

  "engine": {
    "type": "stockfish",
    "path":  "./engines/stockfish.exe",
    "working_dir": "./engines",
    "startup_timeout": 10,
    "analysis_timeout": 300
  },

  "uci_options": {
    "Threads": 4,
    "Hash": 512,
    "MultiPV": 5,
    "Contempt": 0
  },

  "book": {
    "enabled": false,
    "type": "polyglot",
    "path": "./books/performance. bin",
    "best_move_only": false,
    "max_moves": 20
  },

  "tablebase": {
    "enabled": false,
    "type":  "syzygy",
    "path": "./tablebases/syzygy",
    "probe_depth": 1,
    "probe_limit": 7
  },

  "personality": {
    "enabled": false,
    "file": "",
    "style": "default"
  },

  "intelligence": {
    "enabled": false,
    "avoid_low_intelligence":  false,
    "low_intelligence_threshold":  -1.5,
    "aggressiveness_contempt":  1.0,
    "passiveness_contempt": 1.0,
    "trading_preference": 0.0,
    "capture_preference": 1.0,
    "castle_preference": 1.0,
    "en_passant_preference": 1.0,
    "promotion_preference": 1.0,
    "prefer_early_castling": false,
    "prefer_pins": false,
    "prefer_side_castle": false,
    "castle_side":  null,
    "pawn_preference": 1.0,
    "knight_preference": 1.0,
    "bishop_preference":  1.0,
    "rook_preference": 1.0,
    "queen_preference":  1.0,
    "king_preference": 1.0,
    "stay_equal": false,
    "stalemate_probability": 0.0,
    "always_promote_queen": false,
    "checkmate_immediately": true
  },

  "logging": {
    "level": "info",
    "file":  "./logs/server.log",
    "max_size_mb": 100,
    "max_backups": 5
  }
}
```

---

## 3. Engine-Specific Configurations

### 3.1 Stockfish

**Download:** https://stockfishchess.org/download/

**Versions:** Stockfish 16, 17+

```json
{
  "engine": {
    "type": "stockfish",
    "path": "./engines/stockfish/stockfish-windows-x86-64-avx2.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 1024,
    "MultiPV": 5,
    "Contempt": 24,
    "Analysis Contempt": "Both",
    "Move Overhead": 10,
    "Slow Mover": 100,
    "UCI_Chess960": false,
    "UCI_AnalyseMode": true,
    "UCI_LimitStrength": false,
    "UCI_Elo": 3200,
    "Skill Level": 20,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "Syzygy50MoveRule": true,
    "SyzygyProbeLimit": 7
  },
  "book": {
    "enabled": true,
    "type": "polyglot",
    "path": "./books/Titans.bin"
  },
  "tablebase": {
    "enabled": true,
    "type":  "syzygy",
    "path":  "./tablebases/syzygy"
  }
}
```

#### Stockfish Skill Levels

| Level | Approx ELO | Description |
|-------|------------|-------------|
| 0 | 1350 | Beginner |
| 5 | 1750 | Club player |
| 10 | 2100 | Strong amateur |
| 15 | 2500 | Master level |
| 20 | 3500+ | Full strength |

#### Stockfish UCI Options Reference

| Option | Type | Range | Description |
|--------|------|-------|-------------|
| `Threads` | int | 1-1024 | CPU threads to use |
| `Hash` | int | 1-33554432 | Hash table size in MB |
| `MultiPV` | int | 1-500 | Number of principal variations |
| `Contempt` | int | -100 to 100 | Positive = prefer winning, negative = prefer draws |
| `Skill Level` | int | 0-20 | Playing strength (20 = maximum) |
| `UCI_LimitStrength` | bool | | Enable ELO limiting |
| `UCI_Elo` | int | 1320-3190 | Target ELO when limited |
| `Move Overhead` | int | 0-5000 | Time buffer in ms |
| `Slow Mover` | int | 10-1000 | Time management (100 = normal) |
| `SyzygyPath` | string | | Path to Syzygy tablebases |
| `SyzygyProbeDepth` | int | 1-100 | Minimum depth to probe |
| `SyzygyProbeLimit` | int | 0-7 | Max pieces to probe |

---

### 3.2 Leela Chess Zero (Lc0)

**Download:** https://lczero.org/play/download/

**Weights:** https://lczero.org/play/networks/bestnets/

```json
{
  "engine": {
    "type": "lc0",
    "path": "./engines/lc0/lc0.exe"
  },
  "uci_options": {
    "WeightsFile": "./engines/lc0/weights/BT4-1740. pb. gz",
    "Backend": "cuda-fp16",
    "Threads": 2,
    "NNCacheSize": 2000000,
    "MinibatchSize": 256,
    "MaxPrefetch": 32,
    "CPuct": 1.745,
    "CPuctBase": 38739,
    "CPuctFactor": 3.894,
    "FpuValue": 0.443,
    "FpuStrategy": "reduction",
    "PolicyTemperature": 1.359,
    "MultiPV": 3,
    "ScoreType": "centipawn_with_drawscore",
    "Contempt": 0,
    "WDLCalibrationElo": 0,
    "WDLContemptAttenuation": 1.0,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyFastPlay": true,
    "VerboseMoveStats": true
  }
}
```

#### Lc0 Backend Options

| Backend | Description | Requirements |
|---------|-------------|--------------|
| `cuda-fp16` | NVIDIA GPU (fastest) | CUDA 11.x, cuDNN |
| `cuda` | NVIDIA GPU (FP32) | CUDA 11.x |
| `cudnn-fp16` | NVIDIA cuDNN | cuDNN 8.x |
| `opencl` | AMD/Intel GPU | OpenCL 1.2 |
| `dx12` | DirectX 12 GPU | Windows 10+ |
| `eigen` | CPU AVX2 | AVX2 support |
| `blas` | CPU with OpenBLAS | OpenBLAS |
| `random` | Random moves | Testing only |

#### Recommended Weight Files

| Network | Style | ELO | Size |
|---------|-------|-----|------|
| `BT4-1740.pb.gz` | Strongest | 3550+ | 75MB |
| `T80-768x15.pb.gz` | Strong/Fast | 3450+ | 45MB |
| `maia-1900.pb.gz` | Human-like | 1900 | 25MB |
| `maia-1500.pb.gz` | Intermediate | 1500 | 25MB |
| `maia-1100.pb.gz` | Beginner | 1100 | 25MB |

---

### 3.3 Komodo

**Download:** https://komodochess.com/ (Commercial)

```json
{
  "engine":  {
    "type": "komodo",
    "path":  "./engines/komodo/komodo-14.1-64bit. exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 1024,
    "MultiPV": 5,
    "Contempt": 0,
    "Null Move Pruning": "Always",
    "Selectivity": 120,
    "King Safety": 100,
    "Drawishness": 0,
    "Time Usage": "Normal",
    "Analysis Mode": true,
    "SyzygyPath":  "./tablebases/syzygy",
    "Overhead": 0,
    "Threads per GPU": 0
  },
  "personality": {
    "enabled": true,
    "style": "aggressive"
  }
}
```

#### Komodo Playing Styles

Configure via `Selectivity`, `King Safety`, and `Drawishness`:

| Style | Selectivity | King Safety | Drawishness |
|-------|-------------|-------------|-------------|
| Aggressive | 150 | 80 | -20 |
| Solid | 100 | 120 | 20 |
| Active | 120 | 100 | 0 |
| Positional | 80 | 110 | 10 |
| Tactical | 140 | 90 | -10 |

---

### 3.4 Houdini

**Download:** http://www.cruxis.com/chess/houdini. htm (Commercial)

```json
{
  "engine": {
    "type": "houdini",
    "path": "./engines/houdini/Houdini_6_x64.exe"
  },
  "uci_options":  {
    "Threads": 8,
    "Hash":  1024,
    "MultiPV": 5,
    "Contempt": 0,
    "Analysis Contempt": false,
    "Tactical Mode": false,
    "UCI_Chess960": false,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "SyzygyProbeLimit": 6,
    "OwnBook": false,
    "Book File": "./books/houdini.ctg"
  }
}
```

---

### 3.5 Fritz

**Download:** https://shop.chessbase.com/ (Commercial)

```json
{
  "engine": {
    "type": "fritz",
    "path": "./engines/fritz/Fritz18.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 1024,
    "MultiPV":  5,
    "Contempt":  0,
    "Selectivity": 4,
    "Style": 2,
    "SyzygyPath":  "./tablebases/syzygy",
    "Own Book": true,
    "Book File": "./books/fritz.ctg"
  },
  "book": {
    "enabled": true,
    "type": "ctg",
    "path": "./books/Fritz Powerbook 2023.ctg"
  }
}
```

#### Fritz Style Settings

| Style Value | Description |
|-------------|-------------|
| 0 | Cowardly |
| 1 | Passive |
| 2 | Normal |
| 3 | Active |
| 4 | Aggressive |

---

### 3.6 Rybka

**Download:** http://rybkachess.com/ (Commercial)

```json
{
  "engine": {
    "type": "rybka",
    "path": "./engines/rybka/Rybka4.1.exe"
  },
  "uci_options": {
    "Max CPUs": 8,
    "Hash":  1024,
    "MultiPV": 5,
    "Contempt": 0,
    "Ponder": false,
    "Tablebases": "./tablebases/syzygy",
    "Use Tablebases": true,
    "Tablebase Piece Limit": 6,
    "OwnBook": true,
    "BookFile": "./books/rybka.ctg",
    "Time Buffer": 1000
  }
}
```

---

### 3.7 Shredder

**Download:** https://www.shredderchess.com/ (Commercial)

```json
{
  "engine": {
    "type": "shredder",
    "path": "./engines/shredder/Shredder13-UCI.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 1024,
    "MultiPV":  5,
    "UCI_LimitStrength": false,
    "UCI_Elo":  2800,
    "Selectivity": 3,
    "Randomness": 0,
    "SyzygyPath": "./tablebases/syzygy"
  },
  "personality": {
    "enabled": true,
    "style": "normal"
  }
}
```

#### Shredder Strength Limiting

| UCI_Elo | Skill Level | Description |
|---------|-------------|-------------|
| 1000 | Beginner | Very weak |
| 1400 | Casual | Patzer |
| 1800 | Club | Average club player |
| 2200 | Expert | Strong amateur |
| 2600 | Master | Professional level |
| 2800 | Maximum | Full strength |

---

### 3.8 Rodent IV

**Download:** https://github.com/nescitus/Rodent_IV/releases

**Key Feature:** Extensive personality system! 

```json
{
  "engine":  {
    "type": "rodent",
    "path":  "./engines/rodent/RodentIV. exe"
  },
  "uci_options": {
    "Threads":  4,
    "Hash": 256,
    "MultiPV": 3,
    "Contempt": 0,
    "Eval Blur": 0,
    "Mobility Weight": 100,
    "Pawn Structure Weight": 100,
    "Passed Pawns Weight": 100,
    "King Safety Weight": 100,
    "Material Weight": 100,
    "Positional Weight": 100,
    "Outpost Weight": 100,
    "Lines Weight": 100,
    "Own Mobility":  100,
    "Opp Mobility": 100,
    "Own Attack": 100,
    "Opp Attack": 100,
    "Book File": "./books/rodent.bin",
    "Use Book": true,
    "Verbose Book": false,
    "PersonalityFile": "./personalities/fighter.txt"
  },
  "personality": {
    "enabled": true,
    "file": "./personalities/fighter.txt"
  }
}
```

#### Rodent Built-in Personalities

Create personality files in `./personalities/`:

**fighter.txt** (Aggressive tactical player)
```ini
; Rodent IV Personality:  Fighter
Contempt 50
Eval Blur 0
Mobility Weight 130
Pawn Structure Weight 80
Passed Pawns Weight 100
King Safety Weight 120
Own Attack 130
Opp Attack 90
Material Weight 100
Positional Weight 80
Style aggressive
```

**defender.txt** (Solid defensive player)
```ini
; Rodent IV Personality:  Defender
Contempt -10
Eval Blur 0
Mobility Weight 90
Pawn Structure Weight 120
Passed Pawns Weight 110
King Safety Weight 140
Own Attack 80
Opp Attack 120
Material Weight 110
Positional Weight 120
Style defensive
```

**gambit.txt** (Gambit/sacrifice lover)
```ini
; Rodent IV Personality: Gambit King
Contempt 100
Eval Blur 10
Mobility Weight 140
Pawn Structure Weight 60
Passed Pawns Weight 80
King Safety Weight 80
Own Attack 150
Opp Attack 70
Material Weight 80
Positional Weight 70
Style gambit
```

**positional.txt** (Karpov-style positional)
```ini
; Rodent IV Personality: Positional Master
Contempt 0
Eval Blur 0
Mobility Weight 110
Pawn Structure Weight 130
Passed Pawns Weight 120
King Safety Weight 100
Own Attack 90
Opp Attack 100
Material Weight 100
Positional Weight 140
Style positional
```

**random.txt** (Unpredictable/fun)
```ini
; Rodent IV Personality:  Chaos
Contempt 0
Eval Blur 50
Mobility Weight 100
Pawn Structure Weight 100
Passed Pawns Weight 100
King Safety Weight 100
Own Attack 100
Opp Attack 100
Material Weight 100
Positional Weight 100
Style random
```

---

### 3.9 Maia Chess

**Download:** Uses Lc0 binary with Maia weights

**Weights:** https://github.com/CSSLab/maia-chess/releases

```json
{
  "engine": {
    "type": "maia",
    "path": "./engines/lc0/lc0.exe"
  },
  "uci_options":  {
    "WeightsFile": "./engines/maia/weights/maia-1900.pb.gz",
    "Backend":  "eigen",
    "Threads": 1,
    "MinibatchSize": 1,
    "MaxPrefetch": 0,
    "NNCacheSize": 200000,
    "SmartPruningFactor": 0,
    "VerboseMoveStats": false,
    "MultiPV": 1,
    "Temperature": 0,
    "TempDecayMoves": 0
  }
}
```

#### Maia Rating Levels

| Weight File | Target ELO | Playstyle |
|-------------|------------|-----------|
| `maia-1100.pb.gz` | 1100 | Beginner, frequent blunders |
| `maia-1200.pb.gz` | 1200 | Casual player |
| `maia-1300.pb.gz` | 1300 | Improving beginner |
| `maia-1400.pb.gz` | 1400 | Low club level |
| `maia-1500.pb.gz` | 1500 | Average club player |
| `maia-1600.pb.gz` | 1600 | Solid club player |
| `maia-1700.pb.gz` | 1700 | Strong club player |
| `maia-1800.pb.gz` | 1800 | Expert level |
| `maia-1900.pb.gz` | 1900 | High expert/low master |

---

### 3.10 Ethereal

**Download:** https://github.com/AndyGrant/Ethereal/releases

```json
{
  "engine": {
    "type": "ethereal",
    "path": "./engines/ethereal/Ethereal-14.25.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 512,
    "MultiPV": 5,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "MoveOverhead": 100,
    "UCI_Chess960": false
  }
}
```

---

### 3.11 Berserk

**Download:** https://github.com/jhonnold/berserk/releases

```json
{
  "engine": {
    "type":  "berserk",
    "path":  "./engines/berserk/berserk-12-x64-avx2.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 512,
    "MultiPV": 3,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "SyzygyProbeLimit": 7,
    "MoveOverhead": 50
  }
}
```

---

### 3.12 Koivisto

**Download:** https://github.com/Luecx/Koivisto/releases

```json
{
  "engine":  {
    "type": "koivisto",
    "path": "./engines/koivisto/Koivisto-9.0-bmi2.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 512,
    "MultiPV": 3,
    "SyzygyPath":  "./tablebases/syzygy",
    "Move Overhead": 20
  }
}
```

---

### 3.13 RubiChess

**Download:** https://github.com/Matthies/RubiChess/releases

```json
{
  "engine": {
    "type": "rubichess",
    "path": "./engines/rubichess/RubiChess-2.3.1-x64-avx2.exe"
  },
  "uci_options": {
    "Threads":  8,
    "Hash": 512,
    "MultiPV": 3,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 8,
    "Move Overhead": 50,
    "NNUEPath": "./engines/rubichess/net. nnue"
  }
}
```

---

### 3.14 Igel

**Download:** https://github.com/vshcherbyna/igel/releases

```json
{
  "engine":  {
    "type": "igel",
    "path":  "./engines/igel/igel-3.4.0-x64-avx2.exe"
  },
  "uci_options":  {
    "Threads": 8,
    "Hash":  512,
    "MultiPV":  3,
    "SyzygyPath": "./tablebases/syzygy"
  }
}
```

---

### 3.15 Arasan

**Download:** https://www.arasanchess.org/

```json
{
  "engine": {
    "type":  "arasan",
    "path": "./engines/arasan/arasanx-64.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash": 512,
    "MultiPV": 3,
    "SyzygyPath": "./tablebases/syzygy",
    "SyzygyProbeDepth": 6,
    "Book": true,
    "Book File": "./books/arasan.bin",
    "Learning":  false,
    "Move Overhead": 0,
    "Position Learning": false
  }
}
```

---

### 3.16 Crafty

**Download:** http://www.craftychess.com/

```json
{
  "engine":  {
    "type": "crafty",
    "path":  "./engines/crafty/crafty-25.2.exe"
  },
  "uci_options": {
    "Hash": 512,
    "Threads": 8,
    "Log": false,
    "Ponder": false,
    "Learning": false,
    "Personality": "./personalities/crafty_normal.rc",
    "SyzygyPath":  "./tablebases/syzygy",
    "Syzygy Use": true,
    "Syzygy Probe Depth": 1,
    "BookPath": "./books/crafty.bin"
  }
}
```

---

### 3.17 Toga II

**Download:** https://sourceforge.net/projects/togaii/

```json
{
  "engine": {
    "type": "toga",
    "path":  "./engines/toga/TogaII_4.01.exe"
  },
  "uci_options": {
    "Hash": 256,
    "Threads": 4,
    "MultiPV": 3,
    "Ponder": false,
    "OwnBook": true,
    "Book File": "./books/toga.bin",
    "Nalimov Tablebases": true,
    "Nalimov Path": "./tablebases/nalimov",
    "Nalimov Cache": 32
  }
}
```

---

### 3.18 Fruit

**Download:** http://www.fruitchess.com/

```json
{
  "engine": {
    "type": "fruit",
    "path":  "./engines/fruit/Fruit-2.3.1.exe"
  },
  "uci_options": {
    "Hash":  256,
    "MultiPV": 3,
    "Ponder": false,
    "OwnBook": true,
    "Book File": "./books/fruit.bin",
    "Nalimov Tablebases": true,
    "Nalimov Path": "./tablebases/nalimov",
    "Nalimov Cache": 32
  }
}
```

---

### 3.19 Spike

**Download:** http://spike.lazypics.de/

```json
{
  "engine": {
    "type": "spike",
    "path":  "./engines/spike/Spike1.4.exe"
  },
  "uci_options": {
    "Hash":  256,
    "Threads": 4,
    "MultiPV": 3,
    "Ponder": false,
    "OwnBook": true,
    "Book File": "./books/spike.bin"
  }
}
```

---

### 3.20 SugaR

**Download:** https://github.com/nicola1968/SugaR-AI-ICCF/releases

(Stockfish derivative with anti-computer chess features)

```json
{
  "engine": {
    "type": "sugar",
    "path": "./engines/sugar/SugaR-AI-ICCF. exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash":  1024,
    "MultiPV": 5,
    "Contempt": 24,
    "Analysis Contempt": "Both",
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "SyzygyProbeLimit": 7,
    "Variety": 0,
    "Persisted learning": false,
    "Read only learning": false
  }
}
```

---

### 3.21 Brainfish

**Download:** https://github.com/nicola1968/brainfish/releases

(Stockfish with PolyGlot book support)

```json
{
  "engine": {
    "type": "brainfish",
    "path": "./engines/brainfish/Brainfish. exe"
  },
  "uci_options": {
    "Threads":  8,
    "Hash": 1024,
    "MultiPV": 5,
    "Contempt": 24,
    "SyzygyPath":  "./tablebases/syzygy",
    "BookFile": "./books/brainfish.bin",
    "OwnBook": true,
    "Best Book Move": false,
    "Book Depth": 40
  },
  "book": {
    "enabled": true,
    "type": "polyglot",
    "path": "./books/brainfish. bin",
    "best_move_only": false
  }
}
```

---

### 3.22 Fat Fritz

**Download:** https://shop.chessbase. com/ (Commercial, uses Lc0)

```json
{
  "engine":  {
    "type": "fatfritz",
    "path": "./engines/fatfritz/lc0.exe"
  },
  "uci_options": {
    "WeightsFile": "./engines/fatfritz/FatFritz2. pb.gz",
    "Backend": "cuda-fp16",
    "Threads":  2,
    "NNCacheSize": 2000000,
    "MinibatchSize": 256,
    "MaxPrefetch": 32,
    "MultiPV": 3,
    "SyzygyPath":  "./tablebases/syzygy"
  }
}
```

---

### 3.23 Dragon by Komodo

**Download:** https://komodochess.com/ (Commercial, NNUE-based Komodo)

```json
{
  "engine": {
    "type": "dragon",
    "path": "./engines/dragon/dragon-64bit.exe"
  },
  "uci_options": {
    "Threads": 8,
    "Hash":  1024,
    "MultiPV": 5,
    "Contempt": 0,
    "Analysis Mode": true,
    "SyzygyPath":  "./tablebases/syzygy",
    "Dragon Scale": 100,
    "King Safety Scale": 100
  }
}
```

---

## 4. Opening Books Configuration

### 4.1 Polyglot Books (. bin)

```json
{
  "book": {
    "enabled": true,
    "type": "polyglot",
    "path": "./books/performance.bin",
    "best_move_only": false,
    "book_depth": 30,
    "random_factor": 0
  }
}
```

#### Popular Polyglot Books

| Book | Size | Style | Description |
|------|------|-------|-------------|
| `performance.bin` | 50MB | Balanced | Strong all-around |
| `gm2600.bin` | 15MB | GM games | Grandmaster games only |
| `varied.bin` | 30MB | Diverse | Maximum opening variety |
| `human.bin` | 10MB | Human | Human-like opening choices |
| `komodo.bin` | 20MB | Aggressive | Komodo-tuned openings |
| `cerebellum.bin` | 100MB+ | Deep | Very deep analysis book |

### 4.2 CTG Books (ChessBase format)

```json
{
  "book": {
    "enabled": true,
    "type": "ctg",
    "path": "./books/Mega Database.ctg"
  }
}
```

### 4.3 ABK Books (Arena format)

```json
{
  "book": {
    "enabled": true,
    "type": "abk",
    "path": "./books/arena.abk"
  }
}
```

---

## 5. Endgame Tablebases

### 5.1 Syzygy Tablebases

**Download:** http://tablebase.sesse.net/

```json
{
  "tablebase": {
    "enabled": true,
    "type":  "syzygy",
    "path":  "./tablebases/syzygy",
    "probe_depth": 1,
    "probe_limit": 7,
    "use_50_move_rule":  true
  },
  "uci_options": {
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "Syzygy50MoveRule": true,
    "SyzygyProbeLimit": 7
  }
}
```

#### Syzygy Download Sizes

| Pieces | WDL Size | DTZ Size | Total |
|--------|----------|----------|-------|
| 3-4 | 7 MB | 81 MB | 88 MB |
| 5 | 200 MB | 1.6 GB | 1.8 GB |
| 6 | 68 GB | 150 GB | 218 GB |
| 7 | 8.5 TB | 17 TB | 25.5 TB |

### 5.2 Nalimov Tablebases (Legacy)

```json
{
  "tablebase": {
    "enabled":  true,
    "type": "nalimov",
    "path": "./tablebases/nalimov",
    "cache_size": 32
  },
  "uci_options": {
    "NalimovPath":  "./tablebases/nalimov",
    "NalimovCache": 32
  }
}
```

### 5.3 Gaviota Tablebases

```json
{
  "tablebase": {
    "enabled": true,
    "type":  "gaviota",
    "path":  "./tablebases/gaviota"
  },
  "uci_options": {
    "GaviotaTbPath": "./tablebases/gaviota",
    "GaviotaTbCache": 32
  }
}
```

---

## 6. Personalities & Playing Styles

### 6.1 Using Intelligence System for Personality

The server's intelligence system allows creating custom playing styles:

```json
{
  "intelligence": {
    "enabled": true,
    "aggressiveness_contempt": 1. 5,
    "passiveness_contempt": 0.7,
    "capture_preference": 1.3,
    "castle_preference": 0.8,
    "trading_preference": 3.0,
    "prefer_early_castling": false,
    "knight_preference": 1.2,
    "bishop_preference":  1.1,
    "rook_preference": 0.9,
    "queen_preference":  1.0,
    "checkmate_immediately": true
  }
}
```

### 6.2 Pre-defined Personality Profiles

#### Aggressive Attacker
```json
{
  "intelligence": {
    "enabled": true,
    "aggressiveness_contempt": 1.8,
    "passiveness_contempt":  0.5,
    "capture_preference": 1.5,
    "castle_preference":  0.7,
    "prefer_early_castling": false,
    "king_preference": 0.8,
    "queen_preference": 1.3,
    "checkmate_immediately": true
  }
}
```

#### Solid Defender
```json
{
  "intelligence": {
    "enabled": true,
    "aggressiveness_contempt":  0.7,
    "passiveness_contempt":  1.3,
    "capture_preference": 0.9,
    "castle_preference": 1.5,
    "prefer_early_castling": true,
    "pawn_preference": 1.2,
    "king_preference": 1.3,
    "stay_equal":  true
  }
}
```

#### Positional Grinder
```json
{
  "intelligence":  {
    "enabled": true,
    "aggressiveness_contempt": 0.9,
    "passiveness_contempt":  1.1,
    "trading_preference": 5.0,
    "capture_preference":  1.0,
    "pawn_preference": 1.3,
    "bishop_preference": 1.2,
    "stay_equal": false
  }
}
```

#### Tactical Wizard
```json
{
  "intelligence":  {
    "enabled": true,
    "aggressiveness_contempt": 1.4,
    "passiveness_contempt":  0.8,
    "capture_preference": 1.4,
    "prefer_pins": true,
    "knight_preference": 1.3,
    "queen_preference": 1.2,
    "checkmate_immediately": true
  }
}
```

#### Romantic Player (Gambit Lover)
```json
{
  "intelligence": {
    "enabled": true,
    "aggressiveness_contempt": 2.0,
    "passiveness_contempt": 0.4,
    "capture_preference": 1.2,
    "castle_preference": 0.5,
    "pawn_preference": 0.8,
    "knight_preference": 1.4,
    "bishop_preference":  1.3,
    "checkmate_immediately": true,
    "stay_equal": false
  }
}
```

---

## 7. Intelligence System

### 7.1 Complete Intelligence Configuration

```json
{
  "intelligence":  {
    "enabled": true,
    
    "avoid_low_intelligence": true,
    "low_intelligence_threshold":  -1.5,
    
    "aggressiveness_contempt": 1.0,
    "passiveness_contempt":  1.0,
    "trading_preference": 0.0,
    "capture_preference": 1.0,
    "castle_preference": 1.0,
    "en_passant_preference": 1.0,
    "promotion_preference": 1.0,
    
    "prefer_early_castling": false,
    "prefer_pins": false,
    "prefer_side_castle": false,
    "castle_side":  null,
    
    "pawn_preference": 1.0,
    "knight_preference": 1.0,
    "bishop_preference": 1.0,
    "rook_preference": 1.0,
    "queen_preference": 1.0,
    "king_preference":  1.0,
    
    "stay_equal": false,
    "stalemate_probability": 0.0,
    "always_promote_queen": true,
    "checkmate_immediately": true
  }
}
```

### 7.2 Multiplier Reference

| Setting | Range | Description |
|---------|-------|-------------|
| `aggressiveness_contempt` | 0.1 - 3.0 | Boost aggressive moves (attacks, advances) |
| `passiveness_contempt` | 0.1 - 3.0 | Boost passive moves (retreats, quiet) |
| `trading_preference` | 0.0 - 9.0 | Minimum piece value to accept trades |
| `capture_preference` | 0.1 - 3.0 | Boost capturing moves |
| `castle_preference` | 0.1 - 3.0 | Boost castling moves |
| `en_passant_preference` | 0.1 - 3.0 | Boost en passant captures |
| `promotion_preference` | 0.1 - 3.0 | Boost pawn promotions |
| `pawn_preference` | 0.1 - 3.0 | Boost pawn moves |
| `knight_preference` | 0.1 - 3.0 | Boost knight moves |
| `bishop_preference` | 0.1 - 3.0 | Boost bishop moves |
| `rook_preference` | 0.1 - 3.0 | Boost rook moves |
| `queen_preference` | 0.1 - 3.0 | Boost queen moves |
| `king_preference` | 0.1 - 3.0 | Boost king moves |
| `low_intelligence_threshold` | -3.0 to -1.0 | Avoid moves below this evaluation |

### 7.3 Boolean Settings

| Setting | Description |
|---------|-------------|
| `prefer_early_castling` | Boost castling in moves 1-15 |
| `prefer_pins` | Boost moves that create pins |
| `prefer_side_castle` | Prefer specific castle side |
| `stay_equal` | Avoid winning too much (for handicap) |
| `always_promote_queen` | Never underpromote |
| `checkmate_immediately` | Always play mate when found |

---

## 8. Multi-Engine Setup

### 8.1 Engine Pool Configuration

```json
{
  "engines": [
    {
      "id": "stockfish",
      "name": "Stockfish 17",
      "type": "stockfish",
      "path":  "./engines/stockfish/stockfish. exe",
      "default": true,
      "uci_options": {
        "Threads": 8,
        "Hash": 1024
      }
    },
    {
      "id": "lc0",
      "name": "Leela Chess Zero",
      "type":  "lc0",
      "path": "./engines/lc0/lc0.exe",
      "uci_options": {
        "WeightsFile": "./engines/lc0/weights/BT4.pb.gz",
        "Backend":  "cuda-fp16"
      }
    },
    {
      "id": "maia1900",
      "name": "Maia 1900",
      "type": "maia",
      "path": "./engines/lc0/lc0.exe",
      "uci_options": {
        "WeightsFile": "./engines/maia/maia-1900.pb.gz",
        "Backend": "eigen"
      }
    },
    {
      "id":  "rodent",
      "name": "Rodent Fighter",
      "type": "rodent",
      "path":  "./engines/rodent/RodentIV.exe",
      "uci_options": {
        "PersonalityFile": "./personalities/fighter. txt"
      }
    }
  ],
  "active_engine": "stockfish"
}
```

### 8.2 Switching Engines via WebSocket

```
# List available engines
listengines

# Response
engines stockfish,lc0,maia1900,rodent

# Switch engine
setengine lc0

# Response  
engineset lc0
```

---

## 9. Network & Security

### 9.1 TLS/WSS Configuration

```json
{
  "server": {
    "address": "0.0.0.0:8443",
    "tls":  true,
    "cert_file": "./certs/server.crt",
    "key_file": "./certs/server. key"
  }
}
```

Generate self-signed certificates: 
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout server.key -out server.crt \
  -subj "/CN=localhost"
```

### 9.2 Authentication

```json
{
  "auth": {
    "require_for_write":  true,
    "require_for_read": false,
    "localhost_bypass": true,
    "passkey": "my_secret_key",
    "max_failed_attempts": 3
  }
}
```

### 9.3 Remote Access

For remote access (not localhost):

```json
{
  "server":  {
    "address": "0.0.0.0:8080",
    "tls": true,
    "cert_file": "./certs/fullchain.pem",
    "key_file": "./certs/privkey.pem"
  },
  "auth": {
    "require_for_write": true,
    "require_for_read": true,
    "localhost_bypass": false,
    "passkey": "strong_random_passkey_here"
  }
}
```

---

## 10.  Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| Engine not starting | Check path, permissions, and dependencies |
| "Backend not found" (Lc0) | Install CUDA/cuDNN or use `eigen` backend |
| Book not loading | Verify book format matches `type` setting |
| Tablebase errors | Check path and sufficient disk space |
| Connection refused | Verify address and firewall settings |
| Authentication failed | Check passkey or enable localhost bypass |
| Slow analysis | Increase Threads and Hash size |
| Missing weights (Lc0/Maia) | Download network file to specified path |

### Debug Mode

```bash
./lichess-engine-server. exe -config config.json -debug
```

Enables verbose logging of: 
- All UCI commands sent/received
- WebSocket messages
- Intelligence calculations
- Book/tablebase probes

### Log Files

```json
{
  "logging": {
    "level": "debug",
    "file": "./logs/server.log",
    "max_size_mb": 100,
    "max_backups": 5
  }
}
```

---

## 11. Complete Example Configurations

### 11.1 Maximum Strength Analysis

```json
{
  "server": {
    "address": "localhost:8080"
  },
  "engine": {
    "type": "stockfish",
    "path": "./engines/stockfish/stockfish-windows-x86-64-avx2.exe"
  },
  "uci_options": {
    "Threads": 16,
    "Hash": 4096,
    "MultiPV": 5,
    "SyzygyPath":  "./tablebases/syzygy",
    "SyzygyProbeDepth": 1,
    "SyzygyProbeLimit": 7
  },
  "book": {
    "enabled": true,
    "type": "polyglot",
    "path": "./books/cerebellum.bin"
  },
  "tablebase": {
    "enabled": true,
    "type":  "syzygy",
    "path":  "./tablebases/syzygy"
  }
}
```

### 11.2 Human-Like Play (Maia 1500)

```json
{
  "server": {
    "address": "localhost: 8080"
  },
  "engine": {
    "type": "maia",
    "path": "./engines/lc0/lc0.exe"
  },
  "uci_options":  {
    "WeightsFile": "./engines/maia/weights/maia-1500.pb.gz",
    "Backend": "eigen",
    "Threads": 1,
    "MinibatchSize": 1,
    "MaxPrefetch":  0
  },
  "intelligence": {
    "enabled": false
  }
}
```

### 11.3 Aggressive Tactical Personality

```json
{
  "server":  {
    "address": "localhost:8080"
  },
  "engine":  {
    "type": "stockfish",
    "path":  "./engines/stockfish/stockfish. exe"
  },
  "uci_options": {
    "Threads":  8,
    "Hash": 512,
    "MultiPV": 5,
    "Contempt": 50
  },
  "book": {
    "enabled": true,
    "type": "polyglot",
    "path": "./books/gambit.bin"
  },
  "intelligence": {
    "enabled": true,
    "aggressiveness_contempt": 1.8,
    "passiveness_contempt":  0.5,
    "capture_preference": 1.5,
    "castle_preference": 0.7,
    "prefer_early_castling": false,
    "knight_preference": 1.3,
    "checkmate_immediately": true
  }
}
```

### 11.4 Rodent with Custom Personality

```json
{
  "server": {
    "address": "localhost:8080"
  },
  "engine": {
    "type": "rodent",
    "path": "./engines/rodent/RodentIV.exe"
  },
  "uci_options": {
    "Threads": 4,
    "Hash": 256,
    "MultiPV": 3,
    "PersonalityFile": "./personalities/fighter.txt",
    "Use Book": true,
    "Book File": "./books/rodent.bin"
  },
  "personality": {
    "enabled": true,
    "file": "./personalities/fighter.txt"
  }
}
```

### 11.5 Neural Network Analysis (Lc0 + GPU)

```json
{
  "server": {
    "address": "localhost: 8080"
  },
  "engine": {
    "type": "lc0",
    "path": "./engines/lc0/lc0.exe"
  },
  "uci_options": {
    "WeightsFile":  "./engines/lc0/weights/BT4-1740.pb.gz",
    "Backend":  "cuda-fp16",
    "Threads": 2,
    "NNCacheSize": 2000000,
    "MinibatchSize": 256,
    "MaxPrefetch": 32,
    "MultiPV": 3,
    "SyzygyPath":  "./tablebases/syzygy"
  },
  "tablebase": {
    "enabled": true,
    "type":  "syzygy",
    "path":  "./tablebases/syzygy"
  }
}
```

### 11.6 Remote Secure Server

```json
{
  "server":  {
    "address": "0.0.0.0:8443",
    "tls": true,
    "cert_file": "/etc/letsencrypt/live/myserver.com/fullchain.pem",
    "key_file": "/etc/letsencrypt/live/myserver. com/privkey. pem",
    "max_connections": 50
  },
  "auth": {
    "require_for_write": true,
    "require_for_read": true,
    "localhost_bypass": false,
    "passkey": "super_secret_passkey_12345",
    "max_failed_attempts": 3
  },
  "engine": {
    "type": "stockfish",
    "path": "/opt/engines/stockfish"
  },
  "uci_options": {
    "Threads": 32,
    "Hash": 8192,
    "MultiPV": 5,
    "SyzygyPath":  "/opt/tablebases/syzygy"
  },
  "logging": {
    "level": "info",
    "file": "/var/log/lichess-engine/server.log"
  }
}
```

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────┐
│                    LICHESS ENGINE SERVER                         │
├─────────────────────────────────────────────────────────────────┤
│ START:      ./lichess-engine-server.exe -config config.json      │
│ DEBUG:     ./lichess-engine-server. exe -config config.json -debug│
│ CONNECT:   ws://localhost:8080/ws                                │
│ SECURE:    wss://localhost:8443/ws                               │
├─────────────────────────────────────────────────────────────────┤
│ ENGINES:   stockfish, lc0, maia, komodo, rodent, houdini, etc.  │
│ BOOKS:     . bin (Polyglot), .ctg (ChessBase), .abk (Arena)      │
│ TABLEBASES:  Syzygy (7-piece), Nalimov, Gaviota                  │
├─────────────────────────────────────────────────────────────────┤
│ COMMANDS:  whoareyou, whatengine, sub, unsub, lock, unlock      │
│ UCI:       position fen <fen>, go depth 20, stop, quit          │
└─────────────────────────────────────────────────────────────────┘
```

---

**Document Version:** 1.0.0  
**Last Updated:** 2026-01-16  
**Repository:** `reductionfear/bmmbariol`

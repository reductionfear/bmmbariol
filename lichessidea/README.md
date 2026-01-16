# Lichess External Engine Chrome Extension

This Chrome extension provides chess analysis with external engine support for Lichess.org. It integrates the entire Lichess Funnies userscript directly into the extension and uses the extension's privileged context to bypass CSP restrictions.

## Features

- **Built-in Stockfish WASM**: Local chess engine for immediate analysis
- **External Engine Support**: Connect to remote/more powerful engines via WebSocket
- **Multiple Game Modes**: Auto-play, hint mode, piece selection mode
- **Human-like Timing**: Configurable delays and move timing patterns
- **Varied Play**: Multi-PV analysis with weighted move selection
- **Arrow Visualization**: Color-coded arrows showing multiple lines
- **Floating Control Dock**: Easy-to-use UI for toggling features

## Installation

1. Open Chrome and navigate to `chrome://extensions/`
2. Enable **Developer mode** (toggle in the top-right corner)
3. Click **Load unpacked**
4. Select the `extension` folder from this repository
5. The extension should now appear in your extensions list
6. Navigate to Lichess.org to use

**Note**: If you're updating the extension, click the refresh icon (🔄) on the extension card in `chrome://extensions/` after pulling the latest changes, then reload any open Lichess.org tabs.

**No Tampermonkey needed**: This extension includes all functionality directly - you don't need a userscript manager.

## Controls

### UI Buttons

The extension adds a floating control dock with these buttons:

- **Hint** - Show engine analysis arrows for current position
- **Auto-ON/OFF** - Toggle auto-move mode
- **Cfg: X** - Cycle through timing presets (7.5s, 15s, 30s)
- **Arr-ON/OFF** - Toggle arrow display
- **Piece-ON/OFF** - Toggle piece selection mode
- **Human-ON/OFF** - Toggle human-like timing
- **Vary-ON/OFF** - Toggle varied move selection
- **Ext-OFF/✓/⟳** - Toggle external engine (green=connected, orange=reconnecting)
- **⚙️** - Configure external engine WebSocket URL

### Keyboard Shortcuts

- **w** - Toggle hint + auto mode
- **p** - Toggle piece selection mode
- **h** - Toggle human mode
- **v** - Toggle varied mode
- **e** - Toggle external engine

## External Engine Setup

The extension supports external engines via the [chesshook-intermediary](https://github.com/0mlml/chesshook-intermediary) WebSocket protocol.

### 1. Setup the intermediary server

```bash
# Clone and build chesshook-intermediary
git clone https://github.com/0mlml/chesshook-intermediary
cd chesshook-intermediary
go build

# Run with your engine
./chesshook-intermediary -engine /path/to/stockfish
```

### 2. Enable External Engine

1. Click the **Ext-OFF** button in the floating dock to toggle external engine mode
2. Click the **⚙️** button to configure the WebSocket URL (default: `ws://localhost:8080/ws`)
3. The button will show:
   - **Ext-✓** (green) when connected
   - **Ext-⟳** (orange) when reconnecting
   - **Ext-OFF** (gray) when disabled

The extension handles CSP bypassing automatically - no additional setup required!

## Configuration Presets

### 7.5s Mode (Blunder Mode)
- Very fast engine analysis (12ms)
- High blunder chance (45%)
- Allows large CP losses (up to 900)
- Minimum delays

### 15s Mode (Balanced)
- Moderate engine analysis (20ms)
- Moderate blunder chance (16%)
- Allows CP losses up to 300
- Balanced timing

### 30s Mode (Careful)
- Longer engine analysis (60ms)
- Low blunder chance (8%)
- Strict CP loss limit (200)
- Longer thinking times

## Advanced Features

### Multi-PV Analysis
The engine analyzes the top 4 moves simultaneously, with color-coded arrows:
- **Green** - Top move (PV1)
- **Orange** - Second best (PV2)
- **Blue** - Third best (PV3)
- **Purple** - Fourth best (PV4)

### Varied Move Selection
When enabled, the script uses weighted random selection across multiple PVs to simulate human play.

### Human-like Timing
When enabled, the script adds realistic delays with panic mode in time trouble and occasional "tanks" (long thinks).


## Verification

To verify the extension is working:

1. Open Lichess.org
2. Open Chrome DevTools (F12)
3. Go to the **Console** tab
4. You should see messages like `[Extension] CSP meta tag removed`
5. The floating control dock should appear on game pages
6. Try clicking the "Hint" button to test engine functionality

## Troubleshooting

### Extension not loading

1. Ensure Developer mode is enabled in `chrome://extensions/`
2. Check that the extension is enabled (toggle switch is on)
3. Verify the extension has host permissions for `*://lichess.org/*`
4. Reload the extension: click the refresh icon on the extension card
5. Hard refresh Lichess.org (Ctrl+Shift+R or Cmd+Shift+R)

### Control dock not appearing

1. Make sure you're on a Lichess game page (not the home page)
2. Check the browser console for JavaScript errors
3. Try refreshing the page

### External engine not connecting

1. Verify chesshook-intermediary is running: `ps aux | grep chesshook`
2. Check the WebSocket URL is correct (default: `ws://localhost:8080/ws`)
3. Check browser console for connection errors
4. Verify no firewall is blocking WebSocket connections

## Technical Details

This extension uses a **background service worker architecture** to bypass CSP restrictions for external engine WebSocket connections:

### Architecture

```
[Lichess Page] <--postMessage--> [Content Script] <--chrome.runtime--> [Background Worker] <--WebSocket--> [External Engine]
```

### Components

1. **Background Service Worker** (background.js):
   - Runs in a privileged context with **NO CSP restrictions**
   - Manages WebSocket connections to external engines
   - Relays engine commands/responses between content script and external engine
   - Maintains connection state and handles reconnection logic

2. **Content Script** (content.js):
   - Runs at `document_start` in an isolated context
   - Bridges communication between the page and background worker
   - Uses `chrome.runtime.connect()` to establish a persistent connection to background
   - Uses `window.postMessage()` to communicate with the injected script
   - Injects chess.js, stockfish.js, and inject.js into the page context

3. **Injected Script** (inject.js):
   - Runs in the page context (can interact with Lichess's JavaScript)
   - Modified version of the main userscript with external engine code replaced
   - Uses `window.postMessage()` instead of direct WebSocket connections
   - All engine commands are relayed through the background service worker

The extension targets only Lichess.org domains (`lichess.org` and `www.lichess.org`).

### Required Permissions

- `storage` - To store user preferences
- `host_permissions` for:
  - `*://lichess.org/*` - To access Lichess.org pages
  - `*://www.lichess.org/*` - To access www.lichess.org pages
  - `ws://localhost/*` - To connect to local WebSocket servers
  - `ws://127.0.0.1/*` - To connect to local WebSocket servers

### Why This Approach?

Even with `web_accessible_resources`, injected scripts still run under the page's Content Security Policy (CSP). Lichess's CSP blocks WebSocket connections to external hosts. By moving the WebSocket connection to the background service worker (which has no CSP restrictions), we completely bypass these limitations while maintaining full functionality.

## Uninstallation

To remove the extension:

1. Go to `chrome://extensions/`
2. Find "Lichess External Engine"
3. Click **Remove**
4. Confirm the removal

## License

This extension is provided as-is for educational and personal use. Use responsibly and in accordance with Lichess.org's Terms of Service.

// Lichess External Engine Chrome Extension
// Injected script for chess analysis with external engine support
// Dependencies: chess.js, stockfish.js (injected separately)

// NOTE: Lichess is a single-page app (SPA). Games often start without a full page reload,
// so we must NOT exit early on the homepage. We simply wait for game DOM nodes to appear.

// --- socket wrapper ---
let webSocketWrapper = null;
let currentAck = 0;
const TrueNativeWebSocket = window.WebSocket; // Native WebSocket in page context
const NativeWebSocket = window.WebSocket; // Save reference to native WebSocket
const webSocketProxy = new Proxy(window.WebSocket, {
  construct: function(target, args) {
    let ws = new target(...args);
    webSocketWrapper = ws;
    ws.addEventListener("message", function(event) {
      try {
        let msg = JSON.parse(event.data);
        if (msg.t === 'move' && msg.d && typeof msg.d.ply !== 'undefined') {
          currentAck = msg.d.ply;
        }
      } catch (e) {}
    });
    return ws;
  }
});
window.WebSocket = webSocketProxy;

window.lichess = window.site;
window.game = new Chess();

// --- Settings State ---
var autoRun = localStorage.getItem('autorun') ?? "0";
var showArrows = localStorage.getItem('showArrows') !== "0";
var autoHint = autoRun == "1";
var pieceSelectMode = localStorage.getItem('pieceSelectMode') === "1";
var humanMode = localStorage.getItem('humanMode') === "1";
var variedMode = localStorage.getItem('variedMode') !== "0";
var configMode = localStorage.getItem('configMode') || "15s";

// --- External Engine Settings ---
var useExternalEngine = localStorage.getItem('useExternalEngine') === "1";
var externalEngineUrl = localStorage.getItem('externalEngineUrl') || "ws://localhost:8080/ws";
var externalEngineWs = null;
var externalEngineReady = false;
var externalEngineConnected = false;
var externalEngineReconnectTimer = null;

// --- CONFIG PRESETS ---
const PRESETS = {
  '7.5s': {
    engineMs: 12,
    varied: {
      maxCpLoss: 900,          // MODIFIED: Allows huge blunders (hanging queen/mate)
      weights: [8, 40, 28, 24],
      maxBlundersPerGame: 50,   // MODIFIED: High limit
      blunderThreshold: 100,
      blunderChance: 0.45,      // MODIFIED: 45% chance to play the bad move
    },
    human: {
      baseDelayMs: 180,
      maxDelayMs: 600,
      premoveDelayMs: 0,
      premoveMaxMs: 10,
      lowPieceDelayMs: 25,
      lowPieceMaxMs: 120,
      premovePieceThreshold: 12,
      lowPieceThreshold: 22,
      quickMoveChance: 0.35,
      quickMoveMs: 0,
      tankChance: 0.008,
      tankMinMs: 250,
      tankMaxMs: 500,
      randomVariance: 0.25,
    }
  },
  '15s': {
    engineMs: 20,
    varied: {
      maxCpLoss: 300,           // Normal safety
      weights: [10, 45, 23, 22],
      maxBlundersPerGame: 10,
      blunderThreshold: 100,
      blunderChance: 0.16,
    },
    human: {
      baseDelayMs: 250,
      maxDelayMs: 800,
      premoveDelayMs: 0,
      premoveMaxMs: 20,
      lowPieceDelayMs: 30,
      lowPieceMaxMs: 150,
      premovePieceThreshold: 10,
      lowPieceThreshold: 20,
      quickMoveChance: 0.25,
      quickMoveMs: 0,
      tankChance: 0.01,
      tankMinMs: 400,
      tankMaxMs: 600,
      randomVariance: 0.27,
    }
  },
  '30s': {
    engineMs: 60,
    varied: {
      maxCpLoss: 200,           // Strict safety
      weights: [30, 55, 10, 5],
      maxBlundersPerGame: 5,
      blunderThreshold: 100,
      blunderChance: 0.08,
    },
    human: {
      baseDelayMs: 500,
      maxDelayMs: 1200,
      premoveDelayMs: 50,
      premoveMaxMs: 150,
      lowPieceDelayMs: 100,
      lowPieceMaxMs: 500,
      premovePieceThreshold: 8,
      lowPieceThreshold: 16,
      quickMoveChance: 0.20,
      quickMoveMs: 60,
      tankChance: 0.05,
      tankMinMs: 1000,
      tankMaxMs: 2000,
      randomVariance: 0.37,
    }
  }
};

// --- Active Config Globals ---
if (!PRESETS[configMode]) configMode = '15s';
let activeHuman = PRESETS[configMode].human;
let activeVaried = PRESETS[configMode].varied;
let activeEngineMs = PRESETS[configMode].engineMs;

function applyConfig(mode) {
  configMode = mode;
  localStorage.setItem('configMode', mode);
  activeHuman = PRESETS[mode].human;
  activeVaried = PRESETS[mode].varied;
  activeEngineMs = PRESETS[mode].engineMs;
  console.log(`[Config] Applied preset: ${mode}`);
}

// --- Stats & State Variables ---
let humanTimingStats = { totalMoves: 0, totalTimeMs: 0, engineTimeMs: 0 };
let varietyStats = { pv1: 0, pv2: 0, pv3: 0, pv4: 0, blunders: 0 };
let gameBlunderCount = 0;
let cachedPVs = null;
let cachedPVsFen = null;
let cachedPieceCount = null;
let cachedFen = null;
let engineReady = false;
let pendingMove = false;
let isProcessing = false;
let lastMoveSent = null;
let lastMoveSentTime = 0;

// Dynamic Panic Threshold (Defaults to 1.5, changes per turn)
let panicThreshold = 1.5;

// --- Helpers ---
function waitForElement(sel) {
  return new Promise(res => {
    const el = document.querySelector(sel);
    if (el) { res(el); return; }
    const obs = new MutationObserver(() => {
      const el = document.querySelector(sel);
      if (el) { obs.disconnect(); res(el); }
    });
    obs.observe(document.documentElement, { childList: true, subtree: true });
  });
}

// --- Clock Parser ---
function getClockSeconds() {
  let clockEl = document.querySelector('.rclock-bottom .time');

  if (!clockEl) {
    const cgWrap = document.querySelector('.cg-wrap');
    if (cgWrap) {
      const isWhite = cgWrap.classList.contains('orientation-white');
      const colorClass = isWhite ? 'rclock-white' : 'rclock-black';
      clockEl = document.querySelector(`.rclock-bottom.${colorClass} .time`);
    }
  }

  if (!clockEl) {
    clockEl = document.querySelector('.rclock.rclock-bottom .time');
  }

  if (!clockEl) {
    return 999;
  }

  const text = clockEl.textContent || '';
  const match = text.match(/(\d+):(\d+)(?:\.(\d))?/);
  if (!match) {
    return 999;
  }

  const mins = parseInt(match[1], 10);
  const secs = parseInt(match[2], 10);
  const tenths = match[3] ? parseInt(match[3], 10) / 10 : 0;

  return mins * 60 + secs + tenths;
}

function getArrowCoords(sq, color) {
  const f = sq[0].toLowerCase(), r = sq[1];
  let x = { a: -3.5, b: -2.5, c: -1.5, d: -0.5, e: 0.5, f: 1.5, g: 2.5, h: 3.5 }[f];
  let y = { 1: 3.5, 2: 2.5, 3: 1.5, 4: 0.5, 5: -0.5, 6: -1.5, 7: -2.5, 8: -3.5 }[r];
  if (color == "black") { x = -x; y = -y; }
  return [x, y];
}

function coordsToSquare(x, y, board) {
  const rect = board.getBoundingClientRect();
  const sz = rect.width / 8;
  const isWhite = board.classList.contains('orientation-white');
  let fi = Math.floor((x - rect.left) / sz);
  let ri = Math.floor((y - rect.top) / sz);
  if (isWhite) ri = 7 - ri; else fi = 7 - fi;
  if (fi < 0 || fi > 7 || ri < 0 || ri > 7) return null;
  return 'abcdefgh'[fi] + (ri + 1);
}

function countPieces() {
  const fen = game.fen();
  if (fen === cachedFen && cachedPieceCount !== null) return cachedPieceCount;
  cachedFen = fen;
  cachedPieceCount = (fen.split(' ')[0].match(/[pnbrqkPNBRQK]/g) || []).length;
  return cachedPieceCount;
}

// --- VARIED SELECTION (MODIFIED FOR 7.5S BLUNDERS) ---
function selectVariedMove(pvs) {
  if (!pvs || pvs.length === 0) return null;

  const valid = [];
  for (let i = 0; i < pvs.length && i < 4; i++) {
    if (pvs[i] && pvs[i].firstMove) valid.push({ ...pvs[i], idx: i });
  }

  if (valid.length === 0) return null;

  const cfg = activeVaried;
  const topEval = valid[0].evalCp || 0;

  let allowBlunder = false;
  if (gameBlunderCount < cfg.maxBlundersPerGame && topEval > -100 && Math.random() < cfg.blunderChance) {
    allowBlunder = true;
    console.log('[Vary] 🎲 Blunder allowed!');
  }

  const candidates = [];

  for (const pv of valid) {
    const cpLoss = topEval - (pv.evalCp || 0);
    const isBlunder = cpLoss >= cfg.blunderThreshold;

    if (pv.evalType === 'mate' && pv.mateVal !== null && pv.mateVal < 0 && pv.mateVal >= -3) continue;

    if (cpLoss > cfg.maxCpLoss) {
      // 7.5s Override: If maxCpLoss is huge (1200), we allow huge blunders
      // 15s/30s: Stays strict (maxCpLoss is ~300)
      if (!allowBlunder) continue;
    }

    let weight = cfg.weights[pv.idx] || 5;

    // Only penalize weight if we are NOT in '7.5s' high-blunder mode
    // If maxCpLoss is > 1000, we treat blunders as equally valid moves to trick evaluations
    if (cfg.maxCpLoss < 1000) {
        weight = weight - (cpLoss * 0.1);
    }

    weight = Math.max(weight, 3);

    candidates.push({ ...pv, weight, cpLoss, isBlunder });
  }

  if (candidates.length === 0) {
    varietyStats.pv1++;
    return { ...valid[0], move: valid[0].firstMove };
  }

  const totalWeight = candidates.reduce((s, c) => s + c.weight, 0);
  let rand = Math.random() * totalWeight;
  let selected = candidates[0];

  for (const c of candidates) {
    rand -= c.weight;
    if (rand <= 0) { selected = c; break; }
  }

  if (selected.idx === 0) varietyStats.pv1++;
  else if (selected.idx === 1) varietyStats.pv2++;
  else if (selected.idx === 2) varietyStats.pv3++;
  else varietyStats.pv4++;

  if (selected.isBlunder) {
    gameBlunderCount++;
    varietyStats.blunders++;
    console.log(`[Vary] ⚠️ BLUNDER! (${gameBlunderCount}/${cfg.maxBlundersPerGame})`);
  }

  return { ...selected, move: selected.firstMove };
}

// --- TIMING ---
function calculateHumanDelay(uci) {
  const cfg = activeHuman;
  const clockSecs = getClockSeconds();

  // LOW TIME - instant (Uses variable threshold)
  if (clockSecs < panicThreshold) {
    return 0;
  }

  // CAPTURES - always instant
  if (uci && uci.length >= 4) {
    const targetSquare = uci.substring(2, 4);
    const targetPiece = game.get(targetSquare);
    if (targetPiece) {
      return 0;
    }
  }

  const pc = countPieces();

  // Premove mode
  if (pc <= cfg.premovePieceThreshold) {
    const delay = cfg.premoveDelayMs + Math.random() * (cfg.premoveMaxMs - cfg.premoveDelayMs);
    return Math.max(0, Math.round(delay));
  }

  // Low piece mode
  if (pc <= cfg.lowPieceThreshold) {
    const delay = cfg.lowPieceDelayMs + Math.random() * (cfg.lowPieceMaxMs - cfg.lowPieceDelayMs);
    return Math.max(0, Math.round(delay));
  }

  // Normal mode
  let delay = cfg.baseDelayMs;
  delay *= (1 + (Math.random() * 2 - 1) * cfg.randomVariance);

  const roll = Math.random();
  if (roll < cfg.quickMoveChance) {
    delay = cfg.quickMoveMs + Math.random() * 50;
  } else if (roll < cfg.quickMoveChance + cfg.tankChance) {
    delay = cfg.tankMinMs + Math.random() * (cfg.tankMaxMs - cfg.tankMinMs);
  }

  delay = Math.max(0, Math.min(delay, cfg.maxDelayMs));

  if (humanTimingStats.totalMoves > 5) {
    const avg = (humanTimingStats.totalTimeMs + humanTimingStats.engineTimeMs) / humanTimingStats.totalMoves;
    if (avg > 580) delay *= Math.max(0.5, 580 / avg);
  }

  return Math.round(delay);
}

function updateTimingStats(delayMs, engineMs = 0) {
  humanTimingStats.totalMoves++;
  humanTimingStats.totalTimeMs += delayMs;
  humanTimingStats.engineTimeMs += engineMs;
}

function resetStats() {
  humanTimingStats = { totalMoves: 0, totalTimeMs: 0, engineTimeMs: 0 };
  varietyStats = { pv1: 0, pv2: 0, pv3: 0, pv4: 0, blunders: 0 };
  gameBlunderCount = 0;
  cachedPieceCount = null;
  cachedFen = null;
  cachedPVsFen = null;
  pendingMove = false;
  isProcessing = false;
  lastMoveSent = null;
  lastMoveSentTime = 0;
  console.log('[Stats] Reset');
}

// --- External Engine WebSocket ---
function connectExternalEngine() {
  if (externalEngineWs && (externalEngineWs.readyState === TrueNativeWebSocket.CONNECTING || externalEngineWs.readyState === TrueNativeWebSocket.OPEN)) {
    console.log('[ExtEngine] Already connected or connecting');
    return;
  }

  console.log('[ExtEngine] Connecting to', externalEngineUrl);
  
  try {
    externalEngineWs = new TrueNativeWebSocket(externalEngineUrl);
    
    externalEngineWs.onopen = () => {
      console.log('[ExtEngine] ✅ Connected');
      externalEngineConnected = true;
      externalEngineReady = false;
      
      // Subscribe to engine output
      externalEngineWs.send('sub');
      
      // Query engine info
      externalEngineWs.send('whoareyou');
      externalEngineWs.send('whatengine');
      
      // Configure engine
      setTimeout(() => {
        sendToEngine('uci');
        sendToEngine('setoption name Threads value 1');
        sendToEngine('setoption name Contempt value 20');
        sendToEngine(`setoption name MultiPV value ${SF_THREADS}`);
        sendToEngine('isready');
      }, 100);
    };
    
    externalEngineWs.onmessage = (e) => {
      const data = e.data;
      console.log('[ExtEngine]', data);
      
      if (data.startsWith('iam ')) {
        console.log('[ExtEngine] Server:', data);
      } else if (data.startsWith('engine ')) {
        console.log('[ExtEngine] Engine:', data);
      } else if (data === 'authok') {
        console.log('[ExtEngine] ✅ Authenticated');
      } else if (data === 'autherr') {
        console.log('[ExtEngine] ❌ Authentication failed');
      } else if (data === 'readyok') {
        externalEngineReady = true;
        engineReady = true;
        console.log('[ExtEngine] ✅ Ready!');
      } else {
        // Route all other messages (info, bestmove) to listeners
        for (const fn of sfListeners) {
          try { fn({ data }); } catch(x) {}
        }
      }
    };
    
    externalEngineWs.onerror = (err) => {
      console.error('[ExtEngine] ❌ Error:', err);
      externalEngineConnected = false;
      externalEngineReady = false;
    };
    
    externalEngineWs.onclose = () => {
      console.log('[ExtEngine] Disconnected');
      externalEngineConnected = false;
      externalEngineReady = false;
      
      // Auto-reconnect if external engine is enabled
      if (useExternalEngine) {
        console.log('[ExtEngine] Reconnecting in 3s...');
        if (externalEngineReconnectTimer) clearTimeout(externalEngineReconnectTimer);
        externalEngineReconnectTimer = setTimeout(() => {
          connectExternalEngine();
        }, 3000);
      }
    };
  } catch (err) {
    console.error('[ExtEngine] Failed to connect:', err);
    externalEngineConnected = false;
    externalEngineReady = false;
  }
}

function disconnectExternalEngine() {
  if (externalEngineReconnectTimer) {
    clearTimeout(externalEngineReconnectTimer);
    externalEngineReconnectTimer = null;
  }
  
  if (externalEngineWs) {
    externalEngineWs.close();
    externalEngineWs = null;
  }
  
  externalEngineConnected = false;
  externalEngineReady = false;
}

// --- Stockfish ---
const SF_THREADS = 4;
const sfListeners = new Set();

stockfish.onmessage = (e) => {
  const data = String(e.data || '');
  if (data === 'readyok') {
    engineReady = true;
    console.log('[Engine] ✅ Ready!');
  }
  for (const fn of sfListeners) {
    try { fn(e); } catch(x) {}
  }
};

// --- Engine Abstraction Layer ---
function sendToEngine(cmd) {
  if (useExternalEngine && externalEngineConnected && externalEngineWs) {
    externalEngineWs.send(cmd);
  } else {
    stockfish.postMessage(cmd);
  }
}

function configureEngine() {
  return new Promise((resolve) => {
    console.log('[Engine] Configuring...');
    
    if (useExternalEngine) {
      // Connect to external engine
      connectExternalEngine();
      
      // Wait for external engine to be ready
      const checkReady = setInterval(() => {
        if (externalEngineReady) {
          clearInterval(checkReady);
          resolve();
        }
      }, 50);

      setTimeout(() => {
        clearInterval(checkReady);
        if (!externalEngineReady) {
          console.log('[Engine] External engine timeout, falling back to local');
          useExternalEngine = false;
          localStorage.setItem('useExternalEngine', '0');
          configureLocalEngine().then(resolve);
        }
      }, 5000);
    } else {
      configureLocalEngine().then(resolve);
    }
  });
}

function configureLocalEngine() {
  return new Promise((resolve) => {
    console.log('[Engine] Configuring local Stockfish...');
    stockfish.postMessage('uci');
    stockfish.postMessage('setoption name Threads value 1');
    stockfish.postMessage('setoption name Contempt value 20');
    stockfish.postMessage(`setoption name MultiPV value ${SF_THREADS}`);
    stockfish.postMessage('isready');

    const checkReady = setInterval(() => {
      if (engineReady) {
        clearInterval(checkReady);
        resolve();
      }
    }, 50);

    setTimeout(() => {
      clearInterval(checkReady);
      engineReady = true;
      resolve();
    }, 3000);
  });
}

function parseInfoLine(text) {
  if (!text.startsWith('info ')) return null;
  const mpv = text.match(/multipv (\d+)/);
  const cp = text.match(/score cp (-?\d+)/);
  const mate = text.match(/score mate (-?\d+)/);
  const pv = text.match(/ pv (.+)$/);
  if (!pv) return null;

  let evalCp = null, evalType = 'cp', mateVal = null;
  if (cp) {
    evalCp = parseInt(cp[1], 10);
  } else if (mate) {
    mateVal = parseInt(mate[1], 10);
    evalCp = (mateVal > 0 ? 100000 : -100000) + mateVal;
    evalType = 'mate';
  } else {
    return null;
  }

  return {
    multipv: mpv ? parseInt(mpv[1], 10) : 1,
    evalType, evalCp, mateVal,
    pv: pv[1].trim(),
    firstMove: pv[1].trim().split(' ')[0]
  };
}

// MODIFIED FOR VARIABLE PANIC MODE
function getMultiPV(fen, retryCount = 0) {
  return new Promise((resolve) => {
    if (!engineReady) {
      setTimeout(() => getMultiPV(fen, retryCount).then(resolve), 100);
      return;
    }

    const pvs = new Map();
    let resolved = false;

    // Check for low time (Variable threshold)
    const clockSecs = getClockSeconds();
    const isLowTime = clockSecs < panicThreshold;
    const engineTime = isLowTime ? 1 : activeEngineMs;

    const handler = (e) => {
      if (resolved) return;
      const txt = String(e.data || '');

      if (txt.startsWith('info ')) {
        const p = parseInfoLine(txt);
        if (p && p.firstMove) pvs.set(p.multipv, p);
      }

      if (txt.startsWith('bestmove')) {
        resolved = true;
        sfListeners.delete(handler);
        const arr = [...pvs.entries()].sort((a, b) => a[0] - b[0]).map(([, v]) => v);

        if (arr.length === 0 && retryCount < 3) {
          setTimeout(() => getMultiPV(fen, retryCount + 1).then(resolve), 150);
          return;
        }

        // Cache with FEN
        cachedPVs = arr;
        cachedPVsFen = fen;

        resolve(arr);
      }
    };

    // Failsafe: If engine hangs for more than 400ms in low time, resolve anyway
    setTimeout(() => {
      if (!resolved) {
        resolved = true;
        sfListeners.delete(handler);
        const arr = [...pvs.entries()].sort((a, b) => a[0] - b[0]).map(([, v]) => v);
        cachedPVs = arr;
        cachedPVsFen = fen;
        resolve(arr);
      }
    }, isLowTime ? 400 : 2000);

    sfListeners.add(handler);
    sendToEngine('stop');
    sendToEngine('position fen ' + fen);

    if (isLowTime) {
        // PANIC MODE: Depth 1 is instant. Movetime 1 still has overhead.
        sendToEngine('go depth 1');
    } else {
        sendToEngine(`go movetime ${engineTime}`);
    }
  });
}

// --- Drawing ---
const PV_COLORS = [
  { name: 'pv1', hex: '#15781B' },
  { name: 'pv2', hex: '#D35400' },
  { name: 'pv3', hex: '#2980B9' },
  { name: 'pv4', hex: '#8E44AD' }
];

function ensureMarkers() {
  const defs = $('svg.cg-shapes defs')[0];
  if (!defs) return;
  for (const { name, hex } of PV_COLORS) {
    if (!document.getElementById(`arrowhead-${name}`)) {
      defs.innerHTML += `<marker id="arrowhead-${name}" orient="auto" markerWidth="4" markerHeight="8" refX="2.05" refY="2"><path d="M0,0 V4 L3,2 Z" fill="${hex}"></path></marker>`;
    }
  }
}

function drawArrows(pvs) {
  if (!showArrows || !pvs || !pvs.length) return;
  ensureMarkers();
  const layer = $('svg.cg-shapes g')[0];
  if (!layer) return;
  layer.innerHTML = '';

  const seen = new Set();
  const col = $('.cg-wrap')[0].classList.contains('orientation-white') ? 'white' : 'black';
  const topEval = pvs[0]?.evalCp || 0;

  pvs.slice(0, 4).forEach((pv, i) => {
    const m = pv.firstMove;
    if (!m || seen.has(m)) return;
    seen.add(m);

    const pal = PV_COLORS[i];
    const [x1, y1] = getArrowCoords(m.substring(0, 2), col);
    const [x2, y2] = getArrowCoords(m.substring(2, 4), col);

    layer.innerHTML += `<line stroke="${pal.hex}" stroke-width="${0.22 - i*0.015}" stroke-linecap="round" marker-end="url(#arrowhead-${pal.name})" opacity="${1 - i*0.1}" x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}"></line>`;

    const cp = pv.evalCp || 0;
    const cpLoss = topEval - cp;
    let label = pv.evalType === 'mate'
      ? `${pv.mateVal > 0 ? '+' : ''}M${pv.mateVal}`
      : `${cp >= 0 ? '+' : ''}${(cp/100).toFixed(1)}`;
    if (i > 0 && cpLoss > 0) label += ` (-${(cpLoss/100).toFixed(1)})`;

    const w = label.length * 0.13 + 0.5;
    layer.innerHTML += `<rect x="${x2 - w/2}" y="${y2 - 0.4}" width="${w}" height="0.34" rx="0.06" fill="#FFF" opacity="0.9" stroke="${pal.hex}" stroke-width="0.02"></rect>`;
    layer.innerHTML += `<text x="${x2}" y="${y2 - 0.18}" fill="${pal.hex}" text-anchor="middle" font-size="0.24" font-weight="bold">${label}</text>`;
  });
}

// --- Execute Move ---
function executeMove(uci) {
  if (!uci) return false;
  if (!webSocketWrapper || webSocketWrapper.readyState !== 1) {
    setTimeout(() => executeMove(uci), 200);
    return false;
  }

  const cgWrap = $('.cg-wrap')[0];
  if (!cgWrap) return false;

  const myCol = cgWrap.classList.contains('orientation-white') ? 'w' : 'b';
  if (game.turn() !== myCol) return false;

  const now = Date.now();
  const clockSecs = getClockSeconds();

  // PREVENT DUPLICATE MOVES (Only if we have time)
  // If we have < 2.0 seconds, allow spamming (disable duplicate check)
  if (clockSecs > 2.0 && uci === lastMoveSent && (now - lastMoveSentTime) < 500) {
    console.log(`[Exec] ⚠️ Duplicate blocked: ${uci}`);
    return false;
  }

  lastMoveSent = uci;
  lastMoveSentTime = now;

  console.log(`[Exec] ✅ Sending: ${uci}`);
  webSocketWrapper.send(JSON.stringify({ t: "move", d: { u: uci, a: currentAck, b: 0, l: 10000 } }));
  pendingMove = false;
  isProcessing = false;
  return true;
}

function executeMoveHumanized(uci, engineMs = 0) {
  if (!uci) { isProcessing = false; return; }

  const clockSecs = getClockSeconds();

  // LOW TIME - EVERYTHING IS INSTANT (Uses variable threshold)
  if (clockSecs < panicThreshold) {
    console.log(`[⚡ LOW TIME] ${uci} | Clock: ${clockSecs.toFixed(1)}s (Panic: ${panicThreshold.toFixed(1)}s) - INSTANT`);
    executeMove(uci);
    return;
  }

  // Instant capture
  if (uci.length >= 4) {
    const targetSquare = uci.substring(2, 4);
    const targetPiece = game.get(targetSquare);
    if (targetPiece) {
      console.log(`[⚡ CAPTURE] ${uci} | Clock: ${clockSecs.toFixed(1)}s | Engine: ${engineMs}ms`);
      executeMove(uci);
      return;
    }
  }

  const delay = calculateHumanDelay(uci);
  updateTimingStats(delay, engineMs);
  console.log(`[Human] ${uci} | Clock: ${clockSecs.toFixed(1)}s | Delay: ${delay}ms | Engine: ${engineMs}ms`);

  if (delay <= 0) executeMove(uci);
  else setTimeout(() => executeMove(uci), delay);
}

// --- Main Selection ---
function selectBestMove(pvs) {
  if (!pvs || pvs.length === 0) return null;

  let result = null;
  if (variedMode) result = selectVariedMove(pvs);

  if (!result || !result.move) {
    if (pvs[0] && pvs[0].firstMove) {
      result = { ...pvs[0], move: pvs[0].firstMove, idx: 0 };
      varietyStats.pv1++;
    }
  }
  return result;
}

function actOnHint(data, engineMs = 0) {
  if (!data || !data.move) { isProcessing = false; return; }
  const uci = data.move;

  const colorIdx = data.idx || 0;
  const cgWrap = $('.cg-wrap')[0];
  if (cgWrap && showArrows) {
    const col = cgWrap.classList.contains('orientation-white') ? 'white' : 'black';
    const [x1, y1] = getArrowCoords(uci.substring(0, 2), col);
    const [x2, y2] = getArrowCoords(uci.substring(2, 4), col);
    const layer = $('svg.cg-shapes g')[0];
    if (layer) layer.innerHTML += `<line stroke="${PV_COLORS[colorIdx].hex}" stroke-width="0.3" stroke-linecap="round" marker-end="url(#arrowhead-pv1)" opacity="1" x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}"></line>`;
  }

  if (!autoHint) { isProcessing = false; return; }
  if (!webSocketWrapper || !cgWrap) { isProcessing = false; return; }

  const myCol = cgWrap.classList.contains('orientation-white') ? 'w' : 'b';
  if (game.turn() !== myCol) { isProcessing = false; return; }

  if (humanMode) executeMoveHumanized(uci, engineMs);
  else executeMove(uci);
}

// --- Process Turn ---
async function processTurn() {
  if (isProcessing) return;
  const cgWrap = $('.cg-wrap')[0];
  if (!cgWrap) return;

  const myCol = cgWrap.classList.contains('orientation-white') ? 'w' : 'b';
  if (game.turn() !== myCol) return;

  if (!autoHint) return;

  isProcessing = true;
  pendingMove = true;

  // NEW: Randomize Panic Threshold at start of turn (1.0 to 1.7s)
  panicThreshold = 1.0 + Math.random() * 0.7;

  const clockSecs = getClockSeconds();
  const currentFen = game.fen();
  const t0 = performance.now();

  try {
    let pvs;
    let engineTime;

    // LOW TIME MODE - use cached PVs ONLY if they match current position AND we are under panic threshold
    if (clockSecs < panicThreshold && cachedPVs && cachedPVs.length > 0 && cachedPVsFen === currentFen) {
      pvs = cachedPVs;
      engineTime = 0;
      console.log(`[Engine] ⚡ LOW TIME - using cached PVs (valid)`);
    } else {
      // Must compute - but use minimal time if low on clock
      pvs = await getMultiPV(currentFen);
      engineTime = Math.round(performance.now() - t0);
      console.log(`[Engine] Clock: ${clockSecs.toFixed(1)}s | Took: ${engineTime}ms`);
    }

    if (!pvs || pvs.length === 0) {
      isProcessing = false;
      pendingMove = false;
      setTimeout(processTurn, 500);
      return;
    }

    if (!pieceSelectMode) drawArrows(pvs);

    if (!pieceSelectMode) {
      const chosen = selectBestMove(pvs);
      if (chosen && chosen.move) actOnHint(chosen, engineTime);
      else { isProcessing = false; pendingMove = false; }
    } else {
      isProcessing = false;
      pendingMove = false;
    }
  } catch (err) {
    isProcessing = false;
    pendingMove = false;
  }
}

// --- Sync game state ---
function syncGameState() {
  try {
    game = new Chess();
    const moves = $('kwdb, u8t');
    for (let i = 0; i < moves.length; i++) {
      const moveText = moves[i].textContent.replace('✓', '').trim();
      if (moveText) try { game.move(moveText); } catch(e) {}
    }
  } catch(e) {}
}

// --- Piece Select Mode ---
function setupPieceSelectMode() {
  const board = $('.cg-wrap')[0];
  if (!board) return;

  board.addEventListener('click', async (e) => {
    if (!pieceSelectMode) return;
    const sq = coordsToSquare(e.clientX, e.clientY, board);
    if (!sq) return;
    const myCol = board.classList.contains('orientation-white') ? 'w' : 'b';
    if (game.turn() !== myCol) return;
    const piece = game.get(sq);
    if (!piece || piece.color !== myCol) return;

    if (!cachedPVs) cachedPVs = await getMultiPV(game.fen());
    const moves = cachedPVs.filter(p => p.firstMove?.substring(0, 2) === sq);
    if (moves.length > 0) {
      const best = moves[0].firstMove;
      setTimeout(() => humanMode ? executeMoveHumanized(best) : executeMove(best), 50);
    }
  });
}

// --- Floating UI Dock ---
function createBottomDock() {
  const existing = document.getElementById('lf-bottom-dock-content');
  if (existing) return existing;

  const root = document.createElement('div');
  root.id = 'lf-bottom-dock';
  root.style.cssText = [
    'position: fixed',
    'left:50%',
    'bottom:8px',
    'transform:translateX(-50%)',
    'z-index:999999',
    'pointer-events:none'
  ].join(';');

  const bar = document.createElement('div');
  bar.style.cssText = [
    'pointer-events:auto',
    'display:flex',
    'align-items: center',
    'gap:6px',
    'padding:6px 8px',
    'background: rgba(0,0,0,0.35)',
    'border-radius:12px',
    'backdrop-filter:blur(4px)',
    'max-width:calc(100vw - 16px)',
    'overflow-x:auto'
  ].join(';');

  const toggle = document.createElement('button');
  toggle.classList.add('fbt');
  toggle.style.fontSize = '10px';
  toggle.style.padding = '2px 6px';
  toggle.style.minWidth = '28px';

  const content = document.createElement('div');
  content.id = 'lf-bottom-dock-content';
  content.style.cssText = 'display:flex;align-items:center;gap:6px;';

  bar.appendChild(toggle);
  bar.appendChild(content);
  root.appendChild(bar);
  document.body.appendChild(root);

  let collapsed = localStorage.getItem('lfDockCollapsed') === '1';
  function render() {
    content.style.display = collapsed ? 'none' : 'flex';
    toggle.textContent = collapsed ? '▲' : '▼';
    toggle.title = collapsed ? 'Expand' : 'Minimize';
    localStorage.setItem('lfDockCollapsed', collapsed ? '1' : '0');
  }
  toggle.onclick = () => { collapsed = !collapsed; render(); };
  render();

  return content;
}

// --- Main ---
async function run() {
  console.log('[Init] Starting...');
  await configureEngine();
  setupPieceSelectMode();
  syncGameState();

  // Initial Turn Check
  const cgWrap = $('.cg-wrap')[0];
  if (cgWrap) {
    const myCol = cgWrap.classList.contains('orientation-white') ? 'w' : 'b';
    if (game.turn() === myCol && autoHint) setTimeout(processTurn, 500);
  }

  // Observers
  const moveObs = new MutationObserver((muts) => {
    for (const mut of muts) {
      if (mut.addedNodes.length === 0) continue;
      if (mut.addedNodes[0].tagName === "I5Z") continue;

      const lastEl = $('l4x')[0]?.lastChild;
      if (!lastEl) continue;

      try { game.move(lastEl.textContent); } catch (e) {}

      // CLEAR ALL CACHES - position changed
      cachedPVs = null;
      cachedPVsFen = null;
      cachedPieceCount = null;
      cachedFen = null;
      isProcessing = false;
      pendingMove = false;

      setTimeout(processTurn, 100);
    }
  });

  waitForElement('rm6').then((el) => {
    moveObs.observe(el, { childList: true, subtree: true });
    syncGameState();
    setTimeout(processTurn, 500);
  });

  const endObs = new MutationObserver(() => {
    if ($('div.rcontrols')[0]?.textContent.includes("Rematch")) {
      resetStats();
      if (window.lichess?.socket?.send) window.lichess.socket.send("rematch-yes");
      setTimeout(() => { try { $('a.fbt[href^="/?hook_like"]')[0]?.click(); } catch(e) {} }, 1000);
      endObs.disconnect();
    }
  });
  if ($('div.rcontrols')[0]) endObs.observe($('div.rcontrols')[0], { childList: true, subtree: true });

  setInterval(() => {
    if (!isProcessing && autoHint && !pieceSelectMode) {
      const cg = $('.cg-wrap')[0];
      if (cg) {
        const myCol = cg.classList.contains('orientation-white') ? 'w' : 'b';
        if (game.turn() === myCol) processTurn();
      }
    }
  }, 3000);

  // --- UI BUTTONS ---
  const btnCont = createBottomDock() || $('div.ricons')[0];

  // 1. Hint
  const hintBtn = document.createElement('button');
  hintBtn.innerText = 'Hint';
  hintBtn.classList.add('fbt');
  hintBtn.onclick = () => {
    getMultiPV(game.fen()).then(pvs => { if (pvs.length) { cachedPVs = pvs; drawArrows(pvs); } });
  };
  if (btnCont) btnCont.appendChild(hintBtn);

  // 2. Auto Toggle
  const autoBtn = document.createElement('button');
  autoBtn.innerText = autoHint ? 'Auto-ON' : 'Auto-OFF';
  autoBtn.classList.add('fbt');
  autoBtn.style.backgroundColor = autoHint ? "green" : "";
  autoBtn.onclick = () => {
    autoHint = !autoHint;
    localStorage.setItem('autorun', autoHint ? "1" : "0");
    autoBtn.innerText = autoHint ? 'Auto-ON' : 'Auto-OFF';
    autoBtn.style.backgroundColor = autoHint ? "green" : "";
    if (autoHint) { isProcessing = false; processTurn(); }
  };
  if (btnCont) btnCont.appendChild(autoBtn);

  // 3. Config Toggle
  const CFG_ORDER = ['7.5s', '15s', '30s'];
  const CFG_COLORS = { '7.5s': "#8E44AD", '15s': "#A93226", '30s': "#229954" };

  const configBtn = document.createElement('button');
  configBtn.innerText = `Cfg: ${configMode}`;
  configBtn.classList.add('fbt');
  configBtn.style.backgroundColor = CFG_COLORS[configMode] || "#229954";
  configBtn.style.fontSize = "10px";
  configBtn.onclick = () => {
    const idx = Math.max(0, CFG_ORDER.indexOf(configMode));
    const newMode = CFG_ORDER[(idx + 1) % CFG_ORDER.length];
    applyConfig(newMode);
    configBtn.innerText = `Cfg: ${newMode}`;
    configBtn.style.backgroundColor = CFG_COLORS[newMode] || "#229954";
  };
  if (btnCont) btnCont.appendChild(configBtn);

  // 4. Arrow Toggle
  const arrowBtn = document.createElement('button');
  arrowBtn.innerText = showArrows ? 'Arr-ON' : 'Arr-OFF';
  arrowBtn.classList.add('fbt');
  arrowBtn.style.fontSize = "10px";
  arrowBtn.style.opacity = showArrows ? "1" : "0.5";
  arrowBtn.onclick = () => {
    showArrows = !showArrows;
    localStorage.setItem('showArrows', showArrows ? "1" : "0");
    arrowBtn.innerText = showArrows ? 'Arr-ON' : 'Arr-OFF';
    arrowBtn.style.opacity = showArrows ? "1" : "0.5";
    if (!showArrows) { const l = $('svg.cg-shapes g')[0]; if (l) l.innerHTML = ''; }
  };
  if (btnCont) btnCont.appendChild(arrowBtn);

  // 5. Piece Mode
  const pieceBtn = document.createElement('button');
  pieceBtn.innerText = pieceSelectMode ? 'Piece-ON' : 'Piece-OFF';
  pieceBtn.classList.add('fbt');
  pieceBtn.style.fontSize = "9px";
  pieceBtn.style.backgroundColor = pieceSelectMode ? "#2980B9" : "";
  pieceBtn.onclick = () => {
    pieceSelectMode = !pieceSelectMode;
    localStorage.setItem('pieceSelectMode', pieceSelectMode ? "1" : "0");
    pieceBtn.innerText = pieceSelectMode ? 'Piece-ON' : 'Piece-OFF';
    pieceBtn.style.backgroundColor = pieceSelectMode ? "#2980B9" : "";
  };
  if (btnCont) btnCont.appendChild(pieceBtn);

  // 6. Human Mode
  const humanBtn = document.createElement('button');
  humanBtn.innerText = humanMode ? 'Human-ON' : 'Human-OFF';
  humanBtn.classList.add('fbt');
  humanBtn.style.fontSize = "9px";
  humanBtn.style.backgroundColor = humanMode ? "#E74C3C" : "";
  humanBtn.onclick = () => {
    humanMode = !humanMode;
    localStorage.setItem('humanMode', humanMode ? "1" : "0");
    humanBtn.innerText = humanMode ? 'Human-ON' : 'Human-OFF';
    humanBtn.style.backgroundColor = humanMode ? "#E74C3C" : "";
    if (humanMode) resetStats();
  };
  if (btnCont) btnCont.appendChild(humanBtn);

  // 7. Vary Mode
  const varyBtn = document.createElement('button');
  varyBtn.innerText = variedMode ? 'Vary-ON' : 'Vary-OFF';
  varyBtn.classList.add('fbt');
  varyBtn.style.fontSize = "9px";
  varyBtn.style.backgroundColor = variedMode ? "#9B59B6" : "";
  varyBtn.onclick = () => {
    variedMode = !variedMode;
    localStorage.setItem('variedMode', variedMode ? "1" : "0");
    varyBtn.innerText = variedMode ? 'Vary-ON' : 'Vary-OFF';
    varyBtn.style.backgroundColor = variedMode ? "#9B59B6" : "";
  };
  if (btnCont) btnCont.appendChild(varyBtn);

  // 8. External Engine Toggle
  const extEngineBtn = document.createElement('button');
  
  function updateExtEngineBtn() {
    if (useExternalEngine) {
      if (externalEngineConnected) {
        extEngineBtn.innerText = 'Ext-✓';
        extEngineBtn.style.backgroundColor = "#27AE60";
      } else {
        extEngineBtn.innerText = 'Ext-⟳';
        extEngineBtn.style.backgroundColor = "#F39C12";
      }
    } else {
      extEngineBtn.innerText = 'Ext-OFF';
      extEngineBtn.style.backgroundColor = "";
    }
  }
  
  extEngineBtn.classList.add('fbt');
  extEngineBtn.style.fontSize = "9px";
  extEngineBtn.title = "Toggle External Engine";
  updateExtEngineBtn();
  
  extEngineBtn.onclick = () => {
    useExternalEngine = !useExternalEngine;
    localStorage.setItem('useExternalEngine', useExternalEngine ? "1" : "0");
    
    if (useExternalEngine) {
      // Switch to external engine
      engineReady = false;
      connectExternalEngine();
    } else {
      // Switch back to local engine
      disconnectExternalEngine();
      configureLocalEngine();
    }
    
    updateExtEngineBtn();
  };
  
  // Update button status periodically
  setInterval(updateExtEngineBtn, 1000);
  
  if (btnCont) btnCont.appendChild(extEngineBtn);
  
  // 9. External Engine Config
  const extConfigBtn = document.createElement('button');
  extConfigBtn.innerText = '⚙️';
  extConfigBtn.classList.add('fbt');
  extConfigBtn.style.fontSize = "11px";
  extConfigBtn.style.padding = "2px 6px";
  extConfigBtn.title = "Configure External Engine URL";
  extConfigBtn.onclick = () => {
    const newUrl = prompt('External Engine WebSocket URL:', externalEngineUrl);
    if (newUrl && newUrl.trim()) {
      const url = newUrl.trim();
      // Validate WebSocket URL format
      if (!url.startsWith('ws://') && !url.startsWith('wss://')) {
        alert('Invalid WebSocket URL. Must start with ws:// or wss://');
        return;
      }
      
      externalEngineUrl = url;
      localStorage.setItem('externalEngineUrl', externalEngineUrl);
      console.log('[ExtEngine] URL updated to:', externalEngineUrl);
      
      // Reconnect if currently using external engine
      if (useExternalEngine) {
        disconnectExternalEngine();
        setTimeout(() => connectExternalEngine(), 500);
      }
    }
  };
  if (btnCont) btnCont.appendChild(extConfigBtn);

  // Stats Display
  const stats = document.createElement('span');
  stats.id = 'human-stats';
  stats.style.cssText = 'font-size: 9px;color:#888;margin-left:5px;display:none;';
  if (btnCont) btnCont.appendChild(stats);

  setInterval(() => {
    const el = document.getElementById('human-stats');
    if (!el) return;
    if (humanMode && humanTimingStats.totalMoves > 0) {
      const avg = Math.round((humanTimingStats.totalTimeMs + humanTimingStats.engineTimeMs) / humanTimingStats.totalMoves);
      const proj = Math.round(30000 / avg);
      const tot = varietyStats.pv1 + varietyStats.pv2 + varietyStats.pv3 + varietyStats.pv4;
      const pv1p = tot > 0 ? Math.round(varietyStats.pv1 / tot * 100) : 0;
      el.textContent = variedMode
        ? `${avg}ms|${configMode}|PV1: ${pv1p}%|⚠️${gameBlunderCount}`
        : `${avg}ms|${configMode}`;
      el.style.display = 'inline';
    } else {
      el.style.display = 'none';
    }
  }, 1000);

  $('.fbt').on('mousedown', function() {
    this.style.border = '6px solid blue';
    setTimeout((a) => { a.style.border = ''; }, 500, this);
  });

  $(document).on("keydown", (e) => {
    if (e.key === "w") { hintBtn.click(); autoBtn.click(); }
    if (e.key === "p") pieceBtn.click();
    if (e.key === "h") humanBtn.click();
    if (e.key === "v") varyBtn.click();
    if (e.key === "e") extEngineBtn.click();
  });
}

waitForElement('rm6').then(() => run());

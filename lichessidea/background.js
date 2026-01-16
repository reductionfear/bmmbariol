// External engine WebSocket connection (runs in privileged context - no CSP!)
let engineWs = null;
let engineConnected = false;
let engineUrl = 'ws://localhost:8080/ws';
let connectedPorts = new Set();

// Handle connections from content scripts
chrome.runtime.onConnect.addListener((port) => {
  if (port.name !== 'engine-bridge') return;
  
  connectedPorts.add(port);
  console.log('[Background] Content script connected');
  
  // Send current status
  port.postMessage({ type: 'status', connected: engineConnected });
  
  port.onMessage.addListener((msg) => {
    console.log('[Background] Received:', msg);
    
    if (msg.type === 'connect') {
      engineUrl = msg.url || engineUrl;
      connectEngine();
    } else if (msg.type === 'disconnect') {
      disconnectEngine();
    } else if (msg.type === 'send' && engineWs && engineConnected) {
      engineWs.send(msg.data);
    }
  });
  
  port.onDisconnect.addListener(() => {
    connectedPorts.delete(port);
    console.log('[Background] Content script disconnected');
  });
});

function connectEngine() {
  if (engineWs && (engineWs.readyState === WebSocket.CONNECTING || engineWs.readyState === WebSocket.OPEN)) {
    console.log('[Background] Already connected');
    return;
  }
  
  console.log('[Background] Connecting to', engineUrl);
  
  try {
    engineWs = new WebSocket(engineUrl);
    
    engineWs.onopen = () => {
      console.log('[Background] ✅ Connected to engine');
      engineConnected = true;
      broadcast({ type: 'status', connected: true });
      broadcast({ type: 'open' });
    };
    
    engineWs.onmessage = (e) => {
      broadcast({ type: 'message', data: e.data });
    };
    
    engineWs.onerror = (err) => {
      console.error('[Background] ❌ WebSocket error');
      engineConnected = false;
      broadcast({ type: 'status', connected: false });
      broadcast({ type: 'error' });
    };
    
    engineWs.onclose = () => {
      console.log('[Background] WebSocket closed');
      engineConnected = false;
      broadcast({ type: 'status', connected: false });
      broadcast({ type: 'close' });
    };
  } catch (err) {
    console.error('[Background] Failed to connect:', err);
    engineConnected = false;
    broadcast({ type: 'status', connected: false });
  }
}

function disconnectEngine() {
  if (engineWs) {
    engineWs.close();
    engineWs = null;
  }
  engineConnected = false;
  broadcast({ type: 'status', connected: false });
}

function broadcast(msg) {
  for (const port of connectedPorts) {
    try {
      port.postMessage(msg);
    } catch (e) {
      console.warn('[Background] Failed to send message to port, removing:', e.message);
      connectedPorts.delete(port);
    }
  }
}

console.log('[Background] Service worker started');

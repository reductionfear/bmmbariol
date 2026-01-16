// Bridge between page and background service worker
let port = null;

function connectToBackground() {
  port = chrome.runtime.connect({ name: 'engine-bridge' });
  
  port.onMessage.addListener((msg) => {
    // Safety check: ensure msg is an object before spreading
    const data = (typeof msg === 'object' && msg !== null) ? msg : { data: msg };
    window.postMessage({ source: 'lichess-ext-bg', ...data }, '*');
  });
  
  port.onDisconnect.addListener(() => {
    console.log('[Content] Disconnected from background, reconnecting...');
    setTimeout(connectToBackground, 1000);
  });
}

window.addEventListener('message', (e) => {
  // Fixed spacing in 'e.data.source'
  if (e.source !== window || !e.data || e.data.source !== 'lichess-ext-page') return;
  if (port) {
    port.postMessage(e.data);
  }
});

// Inject scripts SEQUENTIALLY - wait for each to load
function injectScriptAndWait(file) {
  return new Promise((resolve) => {
    // Fixed spacing in 'document.createElement'
    const script = document.createElement('script');
    script.src = chrome.runtime.getURL(file);
    script.onload = () => {
      console.log('[Content] Loaded:', file);
      resolve();
    };
    script.onerror = () => {
      // Fixed spacing in 'console.error'
      console.error('[Content] Failed to load:', file);
      resolve(); // Resolve anyway to let the next script try loading
    };
    (document.head || document.documentElement).appendChild(script);
  });
}

async function injectAll() {
  await injectScriptAndWait('stockfish.js');
  await injectScriptAndWait('chess.js');
  await injectScriptAndWait('inject.js');
  console.log('[Content] All scripts loaded');
}

connectToBackground();

if (document.head) {
  injectAll();
} else {
  document.addEventListener('DOMContentLoaded', injectAll);
}

console.log('[Content] Content script loaded');
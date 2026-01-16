// Inject scripts into the page context
function injectScript(file) {
  const script = document.createElement('script');
  script.src = chrome.runtime.getURL(file);
  script.onload = function() {
    this.remove();
  };
  (document.head || document.documentElement).appendChild(script);
}

// Remove CSP meta tag
function removeCSP() {
  const observer = new MutationObserver((mutations, obs) => {
    const csp = document.querySelector('meta[http-equiv="Content-Security-Policy"]');
    if (csp) {
      csp.remove();
      console.log('[Extension] CSP meta tag removed');
    }
  });
  observer.observe(document.documentElement, { childList: true, subtree: true });
  
  // Also try immediately
  const csp = document.querySelector('meta[http-equiv="Content-Security-Policy"]');
  if (csp) {
    csp.remove();
    console.log('[Extension] CSP meta tag removed immediately');
  }
}

// Wait for document to be ready enough, then inject
removeCSP();

// Inject required libraries first, then main script
const scripts = ['chess.js', 'stockfish.js', 'lj.js'];
let loadedCount = 0;

function injectNext() {
  if (loadedCount < scripts.length) {
    const script = document.createElement('script');
    script.src = chrome.runtime.getURL(scripts[loadedCount]);
    script.onload = () => {
      loadedCount++;
      injectNext();
    };
    (document.head || document.documentElement).appendChild(script);
  }
}

// Start injection when DOM is available
if (document.head) {
  injectNext();
} else {
  document.addEventListener('DOMContentLoaded', injectNext);
}

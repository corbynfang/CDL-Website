import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

// Register service worker for cache management
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
          navigator.serviceWorker.register('/sw.js').then(
        () => {
          // console.log('SW registered: ', registration);
        },
        () => {
          // console.log('SW registration failed: ', registrationError);
        }
      );
  });
}

// Clear any existing caches on page load
if ('caches' in window) {
  caches.keys().then((names) => {
    names.forEach((name) => {
      caches.delete(name);
    });
  });
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
// Force rebuild Tue Aug  5 00:08:35 CDT 2025

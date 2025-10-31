import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'

// Unregister all service workers and clear caches on load
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.getRegistrations().then((registrations) => {
    registrations.forEach((registration) => {
      registration.unregister();
    });
  });
}

// Clear all caches
if ('caches' in window) {
  caches.keys().then((names) => {
    names.forEach((name) => {
      caches.delete(name);
    });
  });
}

// Register new service worker after clearing old ones
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    // Small delay to ensure old service workers are unregistered
    setTimeout(() => {
      navigator.serviceWorker.register('/sw.js');
    }, 100);
  });
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

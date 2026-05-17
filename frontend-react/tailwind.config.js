/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // davis7.sh inspired dark theme
        'dark': '#0a0a0a',
        'darker': '#050505',
        'card': '#111111',
        'card-hover': '#161616',
        'border': '#1a1a1a',
        'border-light': '#2a2a2a',
        'text-muted': '#a3a3a3',
        'text-dim': '#737373',
        'text-primary': '#f5f5f5',
        'win': '#22c55e',
        'loss': '#ef4444',
      },
      fontFamily: {
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'sans-serif'],
        grotesk: ['Space Grotesk', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
    },
  },
  plugins: [],
} 
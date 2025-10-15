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
      },
      fontFamily: {
        sans: ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
      },
    },
  },
  plugins: [],
} 
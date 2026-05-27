import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  test: {
    globals: true,
    environment: 'happy-dom',        // 3-5× faster than jsdom
    setupFiles: ['./src/test/setup.ts'],
    environmentMatchGlobs: [
      // Pure utility / service tests — no DOM needed, run in lightweight node env
      ['src/utils/**', 'node'],
      ['src/services/**', 'node'],
    ],
  },
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'https://d2ifjmrbn1xida.cloudfront.net',
        changeOrigin: true,
        secure: true,
      },
    },
  },
  build: {
    rollupOptions: {
      input: {
        main: 'index.html',
      },
      output: {
        manualChunks: {
          'vendor-react': ['react', 'react-dom', 'react-router-dom'],
          'vendor-http':  ['axios'],
        },
      },
    },
  },
})

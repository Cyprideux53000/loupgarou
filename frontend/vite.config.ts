import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/game': 'http://localhost:8080',
      '/status': 'http://localhost:8080',
      '/step': 'http://localhost:8080',
    }
  },
  build: {
    outDir: '../web',
    emptyOutDir: true,
  }
})

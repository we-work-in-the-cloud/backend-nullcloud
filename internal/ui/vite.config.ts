import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  base: '/ui/',
  plugins: [react()],
  build: {
    outDir: 'build',
    minify: 'terser',
    rollupOptions: {
      output: {
        entryFileNames: 'app.js',
        chunkFileNames: 'app.js',
        assetFileNames: '[name][extname]',
        format: 'iife',
      },
    },
    sourcemap: false,
    target: 'es2020',
  },
  server: {
    port: 5173,
  },
})

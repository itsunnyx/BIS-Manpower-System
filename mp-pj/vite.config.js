import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
// export default defineConfig({
//   plugins: [react(), tailwindcss()],
// })

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173, // หรือพอร์ตที่คุณอยากให้ Vite dev ใช้
    proxy: {
      '/api': {
        target: 'http://localhost:8080',  // พอร์ต backend
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
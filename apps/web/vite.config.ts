import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { quasar, transformAssetUrls } from '@quasar/vite-plugin'

// 開發時把 /api 轉到核心 API、/guide 轉到 AI 組長服務（server/README 有啟動方式）
export default defineConfig({
  plugins: [vue({ template: { transformAssetUrls } }), quasar()],
  server: {
    proxy: {
      '/api/': 'http://localhost:8081',
      '/guide/': 'http://localhost:8082',
    },
  },
  // AI 組長頁的 3D 場景（three.js + VRM）約 900 kB，只在進入該頁時才載入
  build: { chunkSizeWarningLimit: 1000 },
  test: { environment: 'jsdom' },
})

import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "vue-i18n": "vue-i18n/dist/vue-i18n.cjs.js",
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  base: "/",
  build: {
    rollupOptions: {
      output: {
        //将所有的未使用资源进行移除
        manualChunks: undefined,
      },
    },
    chunkSizeWarningLimit: 3000,
  },
  server: {
    open: true,
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/blog': {
        target: 'http://localhost:8889',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          const newPath = path.replace(/^\/blog/, '');
          console.log(`[Proxy] Rewriting path: ${path} -> ${newPath}`);
          return newPath;
        },
        ws: true, // 支持 WebSocket
        configure: (proxy, _options) => {
          proxy.on('error', (err, _req, _res) => {
            console.error('[Proxy Error]', err);
          });
          proxy.on('proxyReq', (proxyReq, req, _res) => {
            console.log(`[Proxy] Request: ${req.method} ${req.url} -> http://localhost:8889${proxyReq.path}`);
          });
          proxy.on('proxyRes', (proxyRes, req, _res) => {
            console.log(`[Proxy] Response: ${proxyRes.statusCode} for ${req.url}`);
          });
        },
      },
      //为后端本地图片资源重定向
      '/public': {
        target: 'http://localhost:8889',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => {
          const newPath = path.replace(/^\/public/, '');
          console.log(`[Proxy] Rewriting path: ${path} -> ${newPath}`);
          return newPath;
        },
      },
    },
  }
});

import { fileURLToPath, URL } from "node:url";
import { cpSync } from "node:fs";
import { resolve } from "node:path";
import type { Plugin } from "vite";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";

// 编辑器运行时资源与应用同域，避免外部 CDN 不可达时无法初始化。
function localVditor(): Plugin {
  const root = fileURLToPath(new URL("./node_modules/vditor", import.meta.url));
  let output = "";
  return {
    name: "local-vditor-assets",
    configResolved(config) { output = resolve(config.root, config.build.outDir); },
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        if (req.url?.startsWith("/vendor/vditor/")) req.url = `/@fs/${root}${req.url.slice("/vendor/vditor".length)}`;
        next();
      });
    },
    closeBundle() { cpSync(resolve(root, "dist"), resolve(output, "vendor/vditor/dist"), { recursive: true, filter: path => !path.endsWith(".ts") && !path.includes("/dist/ts") && !path.includes("/dist/types") }); },
  };
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), localVditor(), Components({ dts: false, resolvers: [ElementPlusResolver({ importStyle: "css" })] })],
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

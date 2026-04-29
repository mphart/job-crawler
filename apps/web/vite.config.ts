import { defineConfig } from "vite";

export default defineConfig({
  esbuild: {
    jsxInject: 'import React from "react"',
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/healthz": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});

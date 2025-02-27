import { defineConfig } from "vite";
import liveReload from "vite-plugin-live-reload";
import dotenv from "dotenv";
dotenv.config();

const vitePort = process.env.VITE_PORT || 5173;
const distPath = process.env.APP_DIST_PATH || "static/dist";

export default defineConfig({
  root: ".",
  base: "./",
  plugins: [
    liveReload([
      "./assets/css/**/*.css",
      "./assets/js/**/*.js",
      "./assets/js/**/*.ts",
      "./cmd/**/*.go",
      "./internal/**/*.go",
      "./templates/**/*.html",
      "./templates/**/*.tmpl",
    ]),
  ],
  build: {
    outDir: `./${distPath}`,
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      input: {
        main: "./assets/js/entry.ts", // Adjusted input path relative to the root
        style: "./assets/css/entry.css",
      },
      output: {
        entryFileNames: "js/[name].[hash].js",
        chunkFileNames: "js/[name].[hash].js",
        assetFileNames: ({ name }) => {
          if (/\.(css)$/.test(name ?? "")) {
            return "css/[name].[hash].[ext]";
          }
          return "assets/[name].[hash].[ext]";
        },
      },
    },
  },
  server: {
    proxy: {
      "/static": `http://localhost:${vitePort}`,
    },
    hmr: {
      delay: 400, // Wait 400ms for gohst dev restarts
    },
  },
});

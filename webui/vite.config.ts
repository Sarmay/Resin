import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv, type Plugin } from "vite";

const webuiDir = path.dirname(fileURLToPath(import.meta.url));
const guideFiles = ["user-guide.zh-CN.md", "user-guide.en.md"];

function publishGuides(): Plugin {
  const guideDir = path.join(webuiDir, "src/features/docs");
  const serve = (url: string | undefined, res: { setHeader: (name: string, value: string) => void; end: (body: string) => void }, next: () => void) => {
    const pathname = (url ?? "").split("?")[0];
    const name = pathname.replace(/^\/ui\//, "").replace(/^\//, "");
    if (!guideFiles.includes(name)) {
      next();
      return;
    }
    const file = path.join(guideDir, name);
    res.setHeader("Content-Type", "text/markdown; charset=utf-8");
    res.end(fs.readFileSync(file, "utf8"));
  };

  return {
    name: "resin-publish-guides",
    configureServer(server) {
      server.middlewares.use((req, res, next) => serve(req.url, res, next));
    },
    generateBundle() {
      for (const name of guideFiles) {
        this.emitFile({
          type: "asset",
          fileName: name,
          source: fs.readFileSync(path.join(guideDir, name)),
        });
      }
    },
  };
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const apiTarget = env.VITE_DEV_API_TARGET || "http://127.0.0.1:2260";

  return {
    base: "/ui/",
    plugins: [react(), publishGuides()],
    build: {
      chunkSizeWarningLimit: 1200,
    },
    server: {
      host: "0.0.0.0",
      port: 5173,
      proxy: {
        "/api": {
          target: apiTarget,
          changeOrigin: true,
        },
        "/healthz": {
          target: apiTarget,
          changeOrigin: true,
        },
      },
    },
  };
});

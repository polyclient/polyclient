import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { fileURLToPath, URL } from "url";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: [
      {
        find: "@bindings",
        replacement: fileURLToPath(
          new URL(
            "./bindings/github.com/polyclient/polyclient/pkg/app",
            import.meta.url,
          ),
        ),
      },
    ],
  },
});

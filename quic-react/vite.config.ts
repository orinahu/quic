import { defineConfig } from "vite";
import fs from "fs";
import path from "path";
import react from "@vitejs/plugin-react";
import mkcert from 'vite-plugin-mkcert'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), mkcert()],
  server: {
    https: {
      key: fs.readFileSync(path.resolve(__dirname, "localhost+2-key.pem")),
      cert: fs.readFileSync(path.resolve(__dirname, "localhost+2.pem")),
    },
    port: 4444, // Change to an available port
  },
});
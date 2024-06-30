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
      key: fs.readFileSync(path.resolve(__dirname, "key.pem")),
      cert: fs.readFileSync(path.resolve(__dirname, "cert.pem")),
    },
    port: 4444, // Change to an available port
  },
});
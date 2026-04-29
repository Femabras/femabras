// femabras/frontend/next.config.ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,

  // Allow Next.js dev server to accept HMR websocket connections from the
  // Docker bridge network IP (172.18.0.x). Without this, hot-module reload
  // is blocked and the browser console shows a cross-origin warning.
  allowedDevOrigins: ["172.18.0.0/16"],
};

export default nextConfig;

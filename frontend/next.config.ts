// femabras/frontend/next.config.ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,

  // Allow Next.js dev server to accept HMR websocket connections from the
  // Docker bridge network IP (172.18.0.x). Without this, hot-module reload
  // is blocked and the browser console shows a cross-origin warning.
  allowedDevOrigins: ["172.18.0.0/16"],
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          {
            key: "Content-Security-Policy",
            value:
              "default-src 'self'; " +
              "script-src 'self' 'unsafe-inline' 'unsafe-eval' " +
              "https://securepubads.g.doubleclick.net " +
              "https://pagead2.googlesyndication.com; " +
              "style-src 'self' 'unsafe-inline'; " +
              "img-src 'self' data: blob: https://lh3.googleusercontent.com; " +
              "connect-src 'self' http://localhost:8080 https://api.femabras.com; " +
              "frame-src https://*.doubleclick.net; " +
              "object-src 'none'; frame-ancestors 'none';",
          },
        ],
      },
    ];
  },
};

export default nextConfig;

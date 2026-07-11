import type { NextConfig } from "next";

// In production this app runs behind Nginx, which proxies /api/v1/* to the
// backend on the same origin — no rewrite needed there. For local `next dev`
// (frontend on :3000, backend on :8080, no Nginx in front), this rewrite lets
// the browser keep calling the same relative "/api/v1/..." path so cookies
// set by the backend are still treated as same-origin by the browser.
const backendOrigin = process.env.BACKEND_INTERNAL_URL || "http://localhost:8080";

const nextConfig: NextConfig = {
  output: "standalone",
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${backendOrigin}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;

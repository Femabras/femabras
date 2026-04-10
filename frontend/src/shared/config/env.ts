//femabras/frontend/src/shared/config/env.ts
export const env = {
  // 🟢 Smart Network Routing:
  // If typeof window === "undefined", we are running on the Server (inside Docker), so we use the internal Docker network URL.
  // If window exists, we are running on the Client (Browser), so we use the public localhost URL.
  apiUrl:
    typeof window === "undefined"
      ? process.env.API_INTERNAL_URL ||
        process.env.NEXT_PUBLIC_API_URL ||
        "http://localhost:8080"
      : process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
};

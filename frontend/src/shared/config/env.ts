//frontend/src/shared/config/env.ts
export const env = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
} as const;

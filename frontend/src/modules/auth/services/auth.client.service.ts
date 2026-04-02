//femabras/frontend/src/modules/auth/services/auth.client.service.ts
import { APIError } from "@/shared/lib/errors";
import { env } from "@/shared/config/env";

export const authClientService = {
  async login(identifier: string, password: string) {
    const res = await fetch(`${env.apiUrl}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ identifier, password }),
      credentials: "include",
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new APIError(err?.error || "Login failed", res.status);
    }
    return res.json();
  },

  async register(email: string, password: string) {
    const res = await fetch(`${env.apiUrl}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new APIError(err?.error || "Registration failed", res.status);
    }
    return res.json();
  },

  async verifyOTP(userId: string, otp: string) {
    const res = await fetch(`${env.apiUrl}/verify-otp`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user_id: userId, otp }),
      credentials: "include",
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new APIError(err?.error || "Invalid OTP", res.status);
    }
    return res.json();
  },
};

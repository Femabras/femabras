// femabras/frontend/src/modules/auth/services/auth.client.service.ts
//
// Migrated to apiFetch which automatically attaches the CSRF token from the
// double-submit cookie. Without this, every POST returns 403 from the new
// CSRF middleware on the backend.

import { APIError } from "@/shared/lib/errors";
import { env } from "@/shared/config/env";
import { apiFetch } from "@/shared/lib/api.client";

export const authClientService = {
  async login(identifier: string, password: string) {
    const res = await apiFetch(`${env.apiUrl}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ identifier, password }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new APIError(err?.error || "Login failed", res.status);
    }
    return res.json();
  },

  async logout() {
    await apiFetch(`${env.apiUrl}/logout`, {
      method: "POST",
    });
  },

  async register(name: string, email: string, password: string) {
    const res = await apiFetch(`${env.apiUrl}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, email, password }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new APIError(err?.error || "Registration failed", res.status);
    }
    return res.json();
  },

  async verifyOTP(userId: string, otp: string) {
    const res = await apiFetch(`${env.apiUrl}/verify-otp`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user_id: userId, otp }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => null);
      throw new APIError(err?.error || "Invalid OTP", res.status);
    }
    return res.json();
  },
};

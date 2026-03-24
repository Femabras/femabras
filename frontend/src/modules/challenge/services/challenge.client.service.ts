//frontend/src/modules/challenge/services/challenge.client.service.ts
import { APIError } from "@/shared/lib/errors";
import { env } from "@/shared/config/env";
import { CHALLENGE_CONFIG } from "@/shared/config/constants";
import { GuessResponse } from "../types";

export const challengeClientService = {
  getTodayAttempts(): number {
    if (typeof window === "undefined") return CHALLENGE_CONFIG.MAX_ATTEMPTS;
    const today = new Date().toISOString().split("T")[0];
    const stored = localStorage.getItem(CHALLENGE_CONFIG.STORAGE_KEYS.ATTEMPTS);

    if (!stored) return CHALLENGE_CONFIG.MAX_ATTEMPTS;

    try {
      const data = JSON.parse(stored);
      return data.date === today
        ? data.remaining
        : CHALLENGE_CONFIG.MAX_ATTEMPTS;
    } catch {
      return CHALLENGE_CONFIG.MAX_ATTEMPTS;
    }
  },

  saveTodayAttempts(remaining: number): void {
    if (typeof window === "undefined") return;
    const today = new Date().toISOString().split("T")[0];
    localStorage.setItem(
      CHALLENGE_CONFIG.STORAGE_KEYS.ATTEMPTS,
      JSON.stringify({ date: today, remaining }),
    );
  },

  async submitGuess(guess: string): Promise<GuessResponse> {
    const res = await fetch(`${env.apiUrl}/guess`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ guess }),
      credentials: "include", // <--- ADD THIS LINE! It sends the auth_token to Go.
    });

    if (!res.ok) {
      const errorData = await res.json().catch(() => null);
      throw new APIError(
        errorData?.error || "Failed to submit guess",
        res.status,
      );
    }

    return res.json() as Promise<GuessResponse>;
  },

  async fetchLiveAttempts(): Promise<number> {
    const res = await fetch(`${env.apiUrl}/challenge/attempts`, {
      method: "GET",
      credentials: "include",
    });
    if (!res.ok) return 0;
    const data = await res.json();
    this.saveTodayAttempts(data.remaining_attempts);
    return data.remaining_attempts;
  },
};

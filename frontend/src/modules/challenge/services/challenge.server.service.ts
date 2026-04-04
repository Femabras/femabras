//femabras/frontend/src/modules/challenge/services/challenge.server.service.ts
import { env } from "@/shared/config/env";
import { cookies } from "next/headers";
import { DailyChallengeResponse } from "../types";
import { APIError } from "@/shared/lib/errors";

export const challengeServerService = {
  async getDailyChallenge(): Promise<DailyChallengeResponse> {
    try {
      const res = await fetch(`${env.apiUrl}/challenge`, {
        cache: "no-store",
      });

      if (!res.ok) {
        const errData = await res.json().catch(() => null);
        throw new APIError(
          "Failed to fetch daily challenge",
          res.status,
          errData,
        );
      }

      return (await res.json()) as DailyChallengeResponse;
    } catch (error) {
      if (error instanceof APIError) throw error;
      throw new APIError("Failed to connect to backend", 500);
    }
  },

  async getMyStatus(): Promise<{
    is_winner: boolean;
    payout_status: "unclaimed" | "pending" | "paid" | "rejected";
  }> {
    const cookieStore = await cookies();
    const token = cookieStore.get("auth_token")?.value;

    if (!token) return { is_winner: false, payout_status: "unclaimed" };

    const res = await fetch(
      `${env.apiUrl}/challenge/my-status?t=${Date.now()}`,
      {
        headers: { Cookie: `auth_token=${token}` },
        cache: "no-store",
      },
    );

    if (!res.ok) return { is_winner: false, payout_status: "unclaimed" };
    return res.json();
  },
};

//femabras/frontend/src/modules/challenge/services/challenge.server.service.ts
import { env } from "@/shared/config/env";
import { cookies } from "next/headers";
import { DailyChallengeResponse } from "../types";
import { APIError } from "@/shared/lib/errors";
import { logger } from "@/shared/lib/logger";

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
      logger.error("Failed to fetch daily challenge from backend", error, {
        service: "ChallengeServerService",
        method: "getDailyChallenge",
      });

      if (error instanceof APIError) throw error;
      throw new APIError("Failed to connect to backend", 500);
    }
  },

  async getMyStatus(): Promise<{
    is_winner: boolean;
    payout_status: "unclaimed" | "pending" | "paid" | "rejected";
  }> {
    const cookieStore = await cookies();
    const token = cookieStore.get("access_token")?.value;

    if (!token) return { is_winner: false, payout_status: "unclaimed" };

    try {
      const res = await fetch(
        `${env.apiUrl}/challenge/my-status?t=${Date.now()}`,
        {
          headers: { Cookie: `access_token=${token}` },
          cache: "no-store",
        },
      );

      if (!res.ok) {
        logger.warn("Non-200 response fetching user status", {
          status: res.status,
          service: "ChallengeServerService",
          method: "getMyStatus",
        });
        return { is_winner: false, payout_status: "unclaimed" };
      }

      return await res.json();
    } catch (error) {
      logger.error("Network failure fetching user status", error, {
        service: "ChallengeServerService",
        method: "getMyStatus",
        fallbackValue: "unclaimed",
      });

      return { is_winner: false, payout_status: "unclaimed" };
    }
  },

  async getAttempts(): Promise<number> {
    const cookieStore = await cookies();
    const token = cookieStore.get("access_token")?.value;

    if (!token) return 5;

    try {
      const res = await fetch(
        `${env.apiUrl}/challenge/attempts?t=${Date.now()}`,
        {
          headers: { Cookie: `access_token=${token}` },
          cache: "no-store",
        },
      );

      if (!res.ok) {
        logger.warn("Non-200 response fetching user attempts", {
          status: res.status,
          service: "ChallengeServerService",
          method: "getAttempts",
        });
        return 5;
      }

      const data = await res.json();
      return typeof data.remaining_attempts === "number"
        ? data.remaining_attempts
        : 5;
    } catch (error) {
      logger.error("Network failure fetching user attempts", error, {
        service: "ChallengeServerService",
        method: "getAttempts",
        fallbackValue: 5,
      });

      return 5;
    }
  },
};

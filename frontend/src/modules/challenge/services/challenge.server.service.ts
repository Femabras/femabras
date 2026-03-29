//femabras/frontend/src/modules/challenge/services/challenge.server.service.ts

import { DailyChallengeResponse } from "../types";
import { APIError } from "@/shared/lib/errors";

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export const challengeServerService = {
  async getDailyChallenge(): Promise<DailyChallengeResponse> {
    try {
      const res = await fetch(`${BASE_URL}/challenge`, {
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
};

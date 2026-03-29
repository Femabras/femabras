//femabras/src/modules/challenge/components/ad-button.client.tsx
"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import { GAME } from "@/shared/config/gameStyles";
import { challengeClientService } from "../services/challenge.client.service";

interface AdButtonProps {
  label: string;
  onRewardGranted: (newAttempts: number) => void;
}

export function AdButton({ label, onRewardGranted }: AdButtonProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const [isVerifying, setIsVerifying] = useState(false);

  const handleWatchAd = async () => {
    setIsPlaying(true);

    try {
      await new Promise((resolve) => setTimeout(resolve, 3000));

      setIsPlaying(false);
      setIsVerifying(true);

      setTimeout(async () => {
        const liveAttempts = await challengeClientService.fetchLiveAttempts();
        if (liveAttempts > 0) {
          onRewardGranted(liveAttempts);
        } else {
          alert("Reward verification delayed. Please check back in a moment.");
        }
        setIsVerifying(false);
      }, 1500);
    } catch (error) {
      console.error("Ad Engine Error:", error);
      setIsPlaying(false);
      setIsVerifying(false);
      alert("Ad failed to load. Please disable adblockers.");
    }
  };

  return (
    <button
      onClick={handleWatchAd}
      disabled={isPlaying || isVerifying}
      className={cn(
        GAME.adjustBtn,
        "flex items-center justify-center gap-2 px-8! disabled:opacity-50",
      )}>
      {isPlaying
        ? "📺 Playing Ad..."
        : isVerifying
          ? "🔄 Verifying..."
          : `📺 ${label}`}
    </button>
  );
}

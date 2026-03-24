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
      // 1. Call Google H5 Ads SDK (Mocked here until you paste the actual AdSense snippet)
      // window.adBreak({ type: 'reward', name: 'extra_attempt', ... })
      await new Promise((resolve) => setTimeout(resolve, 3000)); // Simulate watching a 3s ad

      // 2. The Ad finished. Now we show "Verifying..." while we wait for Google
      // to ping our Go server Server-to-Server.
      setIsPlaying(false);
      setIsVerifying(true);

      // Give Google's servers 1.5 seconds to ping our Go server, then check our balance
      setTimeout(async () => {
        const liveAttempts = await challengeClientService.fetchLiveAttempts();
        if (liveAttempts > 0) {
          onRewardGranted(liveAttempts); // This updates the UI and unlocks the board!
        } else {
          alert("Reward verification delayed. Please check back in a moment.");
        }
        setIsVerifying(false);
      }, 1500);
    } catch (error) {
      console.error("Ad Engine Error:", error); // <-- FIXED: Variable is now used
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

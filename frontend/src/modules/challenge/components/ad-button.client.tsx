// femabras/frontend/src/modules/challenge/components/ad-button.client.tsx
"use client";

/**
 * Google Ad Manager (GAM) Rewarded Ad Button
 *
 * SETUP REQUIRED before this works in production:
 * 1. Create a Google Ad Manager account at admanager.google.com
 * 2. Create a "Rewarded" ad unit and copy its full path
 *    e.g. /12345678/femabras_rewarded
 * 3. Add to your .env:
 *      NEXT_PUBLIC_GAM_AD_UNIT_PATH=/12345678/femabras_rewarded
 * 4. Configure Server-Side Verification (SSV) in GAM pointing to:
 *      GET https://yourdomain.com/webhooks/ad-reward
 *
 * AdMob NOTE: AdMob is a native mobile SDK (iOS/Android only).
 * This web implementation uses Google Ad Manager (GAM) — the correct
 * product for rewarded web ads.
 */

import { useState, useEffect, useCallback } from "react";
import { cn } from "@/shared/lib/utils";
import { GAME } from "@/shared/config/gameStyles";
import { challengeClientService } from "../services/challenge.client.service";
import type { Dictionary } from "@/i18n/get-dictionary";

declare global {
  interface Window {
    googletag: {
      cmd: Array<() => void>;
      defineOutOfPageSlot: (
        adUnitPath: string,
        format: unknown,
      ) => {
        addService: (service: unknown) => unknown;
      };
      pubads: () => {
        addEventListener: (
          event: string,
          handler: (e: unknown) => void,
        ) => void;
        enableSingleRequest: () => void;
        refresh: (slots: unknown[]) => void;
      };
      enableServices: () => void;
      display: (slot: unknown) => void;
      OutOfPageFormat?: { REWARDED: unknown };
      destroySlots: (slots: unknown[]) => void;
    };
  }
}

interface AdButtonProps {
  label: string;
  dict: Dictionary["challenge"];
  onRewardGranted: (newAttempts: number) => void;
}

type AdState = "idle" | "loading" | "playing" | "verifying" | "error";

const AD_UNIT_PATH = process.env.NEXT_PUBLIC_GAM_AD_UNIT_PATH || "";

export function AdButton({ label, dict, onRewardGranted }: AdButtonProps) {
  // Lazy initialiser — reads window on first render, avoids setState-in-effect
  const [gptReady, setGptReady] = useState<boolean>(() => {
    if (typeof window === "undefined") return false;
    return !!window.googletag?.cmd;
  });

  const [adState, setAdState] = useState<AdState>("idle");
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  useEffect(() => {
    if (gptReady) return;

    window.googletag = window.googletag || { cmd: [] };

    const script = document.createElement("script");
    script.src = "https://securepubads.g.doubleclick.net/tag/js/gpt.js";
    script.async = true;

    script.onload = () => {
      window.googletag.cmd.push(() => {
        window.googletag.enableServices();
      });
      setGptReady(true);
    };

    script.onerror = () => {
      setAdState("error");
      setErrorMsg(dict.adErrorBlocker);
    };

    document.head.appendChild(script);
  }, [gptReady, dict.adErrorBlocker]);

  const handleRewardGranted = useCallback(async () => {
    setAdState("verifying");
    try {
      const liveAttempts = await challengeClientService.fetchLiveAttempts();
      if (liveAttempts > 0) {
        onRewardGranted(liveAttempts);
        setAdState("idle");
      } else {
        setAdState("error");
        setErrorMsg(dict.adErrorDelayed);
      }
    } catch {
      setAdState("error");
      setErrorMsg(dict.adErrorVerify);
    }
  }, [onRewardGranted, dict.adErrorDelayed, dict.adErrorVerify]);

  const handleWatchAd = useCallback(() => {
    if (!gptReady) {
      setAdState("error");
      setErrorMsg(dict.adErrorNotReady);
      return;
    }

    if (!AD_UNIT_PATH) {
      if (process.env.NODE_ENV === "development") {
        setAdState("playing");
        setTimeout(() => handleRewardGranted(), 2000);
        return;
      }
      setAdState("error");
      setErrorMsg(dict.adErrorFailed);
      return;
    }

    setAdState("loading");
    setErrorMsg(null);

    window.googletag.cmd.push(() => {
      const googletag = window.googletag;

      const rewardedSlot = googletag
        .defineOutOfPageSlot(AD_UNIT_PATH, googletag.OutOfPageFormat?.REWARDED)
        ?.addService(googletag.pubads());

      if (!rewardedSlot) {
        setAdState("error");
        setErrorMsg(dict.adErrorFailed);
        return;
      }

      googletag
        .pubads()
        .addEventListener("slotRenderEnded", (event: unknown) => {
          const e = event as { isEmpty: boolean };
          if (e.isEmpty) {
            googletag.destroySlots([rewardedSlot]);
            setAdState("error");
            setErrorMsg(dict.adErrorNoAds);
          } else {
            setAdState("playing");
          }
        });

      googletag.pubads().addEventListener("rewardedSlotGranted", () => {
        googletag.destroySlots([rewardedSlot]);
        handleRewardGranted();
      });

      googletag.pubads().addEventListener("rewardedSlotClosed", () => {
        googletag.destroySlots([rewardedSlot]);
        setAdState("idle");
      });

      googletag.pubads().refresh([rewardedSlot]);
    });
  }, [gptReady, handleRewardGranted, dict]);

  const isDisabled = adState !== "idle" && adState !== "error";

  const buttonLabel: Record<AdState, string> = {
    idle: `📺 ${label}`,
    loading: `⏳ ${dict.adLoading}`,
    playing: `📺 ${dict.adPlaying}`,
    verifying: `🔄 ${dict.adVerifying}`,
    error: dict.adTryAgain,
  };

  return (
    <div className="flex flex-col items-center gap-2 w-full">
      <button
        onClick={
          adState === "error"
            ? () => {
                setAdState("idle");
                setErrorMsg(null);
              }
            : handleWatchAd
        }
        disabled={isDisabled}
        className={cn(
          GAME.adjustBtn,
          "flex items-center justify-center gap-2 px-8! disabled:opacity-50",
        )}>
        {buttonLabel[adState]}
      </button>

      {errorMsg && (
        <p className="text-xs text-red-400 text-center max-w-xs animate-in fade-in">
          {errorMsg}
        </p>
      )}
    </div>
  );
}

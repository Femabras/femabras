// femabras/frontend/src/modules/challenge/components/winner-podium.client.tsx
"use client";

import Image from "next/image";
import Link from "next/link"; // 🟢 ADDED THIS
import type { DailyChallengeResponse } from "../types";
import type { Dictionary } from "@/i18n/get-dictionary";

interface WinnerPodiumProps {
  challenge: Extract<DailyChallengeResponse, { status: "solved" }>;
  dict: Dictionary["challenge"];
  isAuthenticated: boolean; // 🟢 ADDED THIS
}

export function WinnerPodium({
  challenge,
  dict,
  isAuthenticated,
}: WinnerPodiumProps) {
  const winnerName = challenge.winner.name || dict.podiumFallbackName;
  const initial = winnerName.charAt(0).toUpperCase();

  return (
    <div className="flex flex-col items-center justify-center w-full max-w-md mx-auto p-8 bg-white/5 border border-white/10 rounded-3xl shadow-2xl backdrop-blur-xl animate-in fade-in zoom-in-95 duration-700 relative overflow-hidden">
      {/* Background Glow */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-1/2 bg-yellow-500/20 blur-[60px] pointer-events-none rounded-t-full" />

      {/* Trophy / Avatar Container */}
      <div className="relative z-10 mb-6 group">
        <div className="absolute -inset-4 bg-yellow-500/30 rounded-full blur-xl group-hover:bg-yellow-500/50 transition-colors duration-500" />

        {challenge.winner.picture ? (
          <Image
            src={challenge.winner.picture}
            alt={winnerName}
            width={160}
            height={160}
            unoptimized
            className="relative w-32 h-32 sm:w-40 sm:h-40 rounded-full border-4 border-yellow-400 shadow-2xl object-cover"
            referrerPolicy="no-referrer"
          />
        ) : (
          <div className="relative flex items-center justify-center w-32 h-32 sm:w-40 sm:h-40 rounded-full border-4 border-yellow-400 bg-background shadow-2xl text-5xl font-black text-yellow-400">
            {initial}
          </div>
        )}

        {/* Little Floating Trophy Emoji */}
        <div className="absolute -bottom-2 -right-2 text-4xl sm:text-5xl drop-shadow-xl animate-bounce">
          🏆
        </div>
      </div>

      <div className="text-center z-10 w-full">
        <h2 className="text-2xl sm:text-3xl font-black tracking-tight mb-2 text-yellow-400">
          {dict.podiumTitle}
        </h2>
        <p className="text-sm sm:text-base text-foreground/70 mb-1">
          {dict.podiumSubtitle}
        </p>
        <p className="text-3xl sm:text-4xl font-black text-foreground drop-shadow-md tracking-tighter mb-8">
          {winnerName}
        </p>

        {/* 🟢 CONDITIONAL RENDER: Text for logged-in users, Button for guests */}
        {isAuthenticated ? (
          <div className="inline-block px-4 py-2 bg-white/5 border border-white/10 rounded-xl">
            <p className="text-xs uppercase tracking-widest text-foreground/50">
              {dict.podiumNext}
            </p>
          </div>
        ) : (
          <Link
            href="/register"
            className="inline-flex items-center justify-center w-full px-6 py-4 bg-yellow-500 text-black text-sm font-black uppercase tracking-widest rounded-2xl shadow-xl shadow-yellow-500/20 hover:scale-105 transition-all">
            {dict.podiumRegisterBtn}
          </Link>
        )}
      </div>
    </div>
  );
}

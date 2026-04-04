//femabras/frontend/scr/modules/challenge/components/game-board.client.tsx
"use client";

import { useRef } from "react";
import {
  DndContext,
  PointerSensor,
  TouchSensor,
  useSensor,
  useSensors,
  DragOverlay,
  defaultDropAnimationSideEffects,
} from "@dnd-kit/core";

import { ChallengeSlot } from "./challenge-slot.client";
import { DraggableNumber } from "./draggable-number.client";
import { AuthOverlay } from "./auth-overlay.client";
import { useGameEngine } from "../hooks/useGameEngine";
import type { GameBoardProps } from "../types";
import { GAME } from "@/shared/config/gameStyles";
import { UI } from "../utils/styles";
import { cn } from "@/shared/lib/utils";
import { authClientService } from "@/modules/auth/services/auth.client.service";
import { ClaimPrizeForm } from "./claim-prize-form.client";
import { AdButton } from "./ad-button.client";

export function GameBoard({
  challenge,
  isAuthenticated,
  dict,
}: GameBoardProps) {
  const { state, actions } = useGameEngine(
    challenge?.slots || 0,
    challenge?.digits || [],
    dict || {},
    isAuthenticated,
  );

  const inputRefs = useRef<(HTMLInputElement | null)[]>(
    new Array(challenge?.slots || 0).fill(null),
  );

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 1 } }),
    useSensor(TouchSensor, {
      activationConstraint: { delay: 200, tolerance: 5 },
    }),
  );

  if (!dict || !challenge) {
    return (
      <div className="flex flex-col items-center justify-center min-h-100px">
        <div className="animate-pulse flex flex-col items-center gap-4">
          <div className="h-12 w-48 bg-white/5 rounded-xl" />
          <div className="h-24 w-64 bg-white/5 rounded-2xl" />
        </div>
      </div>
    );
  }

  const activeItem = state.trayDigits.find((d) => d.id === state.activeDragId);
  const showTray =
    !state.isComplete && !state.hasWon && !state.authPrompt.isActive;

  const handleLogout = async () => {
    try {
      await authClientService.logout();
      sessionStorage.removeItem("femabras_saved_guess");
      localStorage.removeItem("femabras_attempts");

      const currentPath = window.location.pathname;
      const locale = currentPath.split("/")[1] || "en";
      window.location.href = `/${locale}`;
    } catch (error) {
      console.error("Logout failed:", error);
      window.location.href = "/";
    }
  };

  return (
    <div className="flex flex-col items-center justify-center touch-none w-full relative min-h-25">
      {isAuthenticated && (
        <div className="fixed top-4 left-0 w-full px-4 sm:px-8 flex justify-between items-center z-40 pointer-events-none">
          {/* Left Side: Attempts */}
          <div className="pointer-events-auto flex items-center gap-3 px-4 py-2 bg-white/5 border border-white/10 rounded-2xl backdrop-blur-md shadow-lg">
            <span className="flex h-6 w-6 sm:h-8 sm:w-8 items-center justify-center rounded-full bg-yellow-500/20 text-yellow-500 text-xs sm:text-sm font-black">
              {state.attempts}
            </span>
            <span className="text-[10px] sm:text-xs font-bold uppercase tracking-widest text-foreground/70">
              {dict.attemptsLeft}
            </span>
          </div>
          {/* Right Side: Prize & Logout */}
          <div className="flex items-center gap-4 sm:gap-6 pointer-events-auto">
            <button
              onClick={handleLogout}
              className="pointer-events-auto px-4 py-2 sm:px-5 sm:py-3 text-[10px] sm:text-xs font-bold uppercase tracking-widest text-foreground/50 hover:text-red-400 hover:bg-red-500/10 rounded-xl transition-all border border-transparent hover:border-red-500/20">
              {dict.logout}
            </button>
          </div>
        </div>
      )}

      {state.toastMsg && (
        <div className="absolute top-4 sm:top-10 z-100 animate-in slide-in-from-top-4 fade-in duration-300">
          <div className={UI.toast}>
            <span>⚠️</span> {state.toastMsg}
          </div>
        </div>
      )}

      {state.authPrompt.isActive && (
        <AuthOverlay
          countdown={state.authPrompt.countdown}
          onConfirm={actions.confirmAuthRedirect}
          onCancel={() =>
            actions.setAuthPrompt({ isActive: false, countdown: 5 })
          }
          dict={dict}
        />
      )}

      <DndContext
        id="femabras-dnd-context"
        sensors={sensors}
        onDragStart={actions.handleDragStart}
        onDragOver={actions.handleDragOver}
        onDragEnd={actions.handleDragEnd}>
        <div
          className="text-center mb-8 sm:mb-10 px-2 transition-opacity"
          style={{ opacity: state.authPrompt.isActive ? 0.1 : 1 }}>
          <h2 className="game-title flex flex-col items-center gap-1 sm:gap-2">
            <span>{state.title}</span>

            {/* Only show the prize if the default title is showing */}
            {state.title === dict.titleDefault && challenge?.prize > 0 && (
              <span className="text-lg sm:text-2xl text-foreground/80 font-bold tracking-normal">
                {dict.earnPrize}{" "}
                <span className="text-yellow-400 font-black drop-shadow-[0_0_15px_rgba(250,204,21,0.6)] tracking-tight">
                  {/* 🟢 THE FIX: Native Locale Currency Formatting */}
                  {new Intl.NumberFormat(
                    typeof window !== "undefined"
                      ? window.location.pathname.split("/")[1] || "en"
                      : "en",
                    {
                      style: "decimal",
                      minimumFractionDigits: 2,
                      maximumFractionDigits: 2,
                    },
                  ).format(challenge.prize)}{" "}
                  AOA
                </span>
              </span>
            )}
          </h2>

          {!state.hasWon &&
            !state.isOutOfAttempts &&
            !state.authPrompt.isActive &&
            (!state.isComplete ? (
              <p className="text-[10px] sm:text-sm opacity-60 italic animate-in fade-in">
                {dict.hint}
              </p>
            ) : !isAuthenticated ? (
              <div className="flex justify-center animate-in fade-in duration-500 mt-4">
                <button
                  onClick={actions.confirmAuthRedirect}
                  className={`${UI.guestBanner} hover:scale-105 hover:bg-white/10 hover:border-white/20 transition-all cursor-pointer group`}>
                  <p className="text-xs sm:text-sm text-foreground/80 font-medium flex items-center gap-2 group-hover:text-foreground transition-colors">
                    <span className="text-yellow-400">✨</span>{" "}
                    {dict.guestBanner}
                  </p>
                </button>
              </div>
            ) : null)}
        </div>

        <div
          className="flex flex-col items-center mb-10 sm:mb-12 w-full transition-opacity"
          style={{ opacity: state.authPrompt.isActive ? 0.1 : 1 }}>
          {state.isOutOfAttempts ? (
            <div className="flex flex-col items-center justify-center space-y-6 py-6 animate-in zoom-in duration-500">
              <div className="text-6xl drop-shadow-2xl">🔒</div>
              <p className="text-foreground/70 text-center max-w-sm text-sm sm:text-base leading-relaxed">
                {dict.outOfAttemptsMsg}
              </p>

              <div className="flex flex-col sm:flex-row gap-4 w-full mt-4">
                <AdButton
                  label={dict.btnWatchAd}
                  onRewardGranted={(newAttempts: number) => {
                    actions.setAttempts(newAttempts);
                  }}
                />

                <button
                  onClick={() =>
                    alert("Payment Provider Integration Pending...")
                  }
                  className={cn(
                    GAME.submitBtn,
                    "flex items-center justify-center gap-2 bg-yellow-500! text-black! shadow-yellow-500/20! hover:scale-105!",
                  )}>
                  💎 {dict.btnBuyAttempts}
                </button>
              </div>
            </div>
          ) : (
            <div className="flex gap-2 sm:gap-5 justify-center max-w-full">
              {state.slots.map((slotObj, i) => (
                <div
                  key={i}
                  className={cn(
                    state.isShaking && GAME.animations.shake,
                    "scale-90 sm:scale-100",
                  )}
                  onClick={() => actions.removeDigit(i)}>
                  <ChallengeSlot
                    indexOrder={i}
                    orderValue={slotObj?.val || ""}
                    inputRef={(el) => (inputRefs.current[i] = el)}
                    onPlace={(idx, val) => {
                      const nextIdx = actions.handleManualInput(idx, val);
                      if (nextIdx !== -1 && nextIdx !== idx)
                        inputRefs.current[nextIdx]?.focus();
                      else if (nextIdx === -1) inputRefs.current[idx]?.blur();
                    }}
                    onKeyDown={(e) => {
                      if (e.key === "Backspace") {
                        if (!state.slots[i] && i > 0) {
                          e.preventDefault();
                          inputRefs.current[i - 1]?.focus();
                          actions.removeDigit(i - 1);
                        } else {
                          actions.removeDigit(i);
                        }
                      }
                    }}
                    onFocus={() => actions.setFocusedIndex(i)}
                    onBlur={() =>
                      setTimeout(() => actions.setFocusedIndex(null), 100)
                    }
                    isShaking={state.isShaking}
                  />
                </div>
              ))}
            </div>
          )}
        </div>

        <div
          className="min-h-25 flex flex-col items-center justify-center transition-opacity w-full"
          style={{ opacity: state.authPrompt.isActive ? 0.1 : 1 }}>
          {state.hasWon ? (
            <ClaimPrizeForm prizeAmount={challenge.prize} dict={dict} />
          ) : state.isOutOfAttempts ? null : (
            <>
              {showTray ? (
                <div className="token-tray animate-in fade-in slide-in-from-bottom-4 duration-700">
                  {state.trayDigits.map((digit) => (
                    <div
                      key={digit.id}
                      className="scale-75 sm:scale-100 origin-center transition-transform hover:scale-110">
                      <DraggableNumber
                        id={digit.id}
                        value={digit.val}
                        onClick={() => actions.handleTrayClick(digit)}
                      />
                    </div>
                  ))}
                </div>
              ) : (
                <div className="flex flex-col sm:flex-row gap-4 animate-in zoom-in duration-300 w-full sm:w-auto px-6">
                  <button
                    onClick={actions.submitSequence}
                    disabled={state.isSubmitting}
                    className={GAME.submitBtn}>
                    {state.isSubmitting ? dict.btnProcessing : dict.btnSubmit}
                  </button>
                  <button
                    onClick={() => actions.removeDigit(challenge.slots - 1)}
                    className={GAME.adjustBtn}>
                    {dict.btnRevise}
                  </button>
                </div>
              )}
            </>
          )}
        </div>

        <DragOverlay
          dropAnimation={
            state.isOverSlot
              ? null
              : {
                  sideEffects: defaultDropAnimationSideEffects({
                    styles: { active: { opacity: "0.5" } },
                  }),
                }
          }>
          {state.activeDragId && activeItem ? (
            <div className="h-14 w-14 sm:h-20 sm:w-20 bg-foreground text-background flex items-center justify-center rounded-2xl text-2xl sm:text-4xl font-black shadow-2xl shadow-white/20 scale-110 cursor-grabbing border border-white/20">
              {activeItem.val}
            </div>
          ) : null}
        </DragOverlay>
      </DndContext>
    </div>
  );
}

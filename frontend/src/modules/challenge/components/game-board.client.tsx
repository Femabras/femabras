// femabras/frontend/src/modules/challenge/components/game-board.client.tsx
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
      <div className="flex flex-col items-center justify-center min-h-25">
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

  const formattedPrize = new Intl.NumberFormat("en", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(challenge.prize);

  const displayTitle =
    state.title === dict.titleDefault
      ? dict.titleDefault.replace("{prize}", formattedPrize)
      : state.title;

  return (
    <div className="flex flex-1 flex-col items-center justify-center touch-none w-full relative px-4 py-8">
      {state.toastMsg && (
        <div className="fixed bottom-10 left-1/2 -translate-x-1/2 z-100 animate-in slide-in-from-bottom-4 fade-in duration-300">
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

      <div className="w-full max-w-4xl flex flex-col items-center gap-8 sm:gap-12">
        <DndContext
          id="femabras-dnd"
          sensors={sensors}
          onDragStart={actions.handleDragStart}
          onDragOver={actions.handleDragOver}
          onDragEnd={actions.handleDragEnd}>
          <div
            className="flex flex-col items-center w-full"
            style={{ opacity: state.authPrompt.isActive ? 0.1 : 1 }}>
            <h1 className="mb-8 sm:mb-12 text-[clamp(1.5rem,5vw,3.5rem)] font-black tracking-tight text-foreground text-center leading-tight sm:leading-snug">
              {displayTitle.split(formattedPrize).map((part, i, arr) => (
                <span key={i}>
                  {part}
                  {i < arr.length - 1 && (
                    <span className="text-brand-gold drop-shadow-[0_0_15px_var(--color-brand-gold-glow)]">
                      {formattedPrize}
                    </span>
                  )}
                </span>
              ))}
            </h1>

            <div className="w-full flex flex-col items-center gap-6">
              {state.isOutOfAttempts ? (
                <div className="flex flex-col items-center justify-center gap-6 py-6 animate-in zoom-in duration-500">
                  <div className="text-6xl drop-shadow-2xl">🔒</div>
                  <p className="text-foreground/70 text-center max-w-sm">
                    {dict.outOfAttemptsMsg}
                  </p>
                  <div className="flex flex-col sm:flex-row gap-4 w-full mt-2">
                    {/* dict prop added — all ad button strings are now i18n */}
                    <AdButton
                      label={dict.btnWatchAd}
                      dict={dict}
                      onRewardGranted={(newAttempts: number) =>
                        actions.setAttempts(newAttempts)
                      }
                    />
                    <button
                      className={cn(
                        GAME.submitBtn,
                        "bg-brand-gold! text-black! shadow-[0_0_15px_var(--color-brand-gold-glow)] hover:scale-105!",
                      )}>
                      💎 {dict.btnBuyAttempts}
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <div className="flex flex-nowrap justify-center items-center gap-2 sm:gap-4 mx-auto w-full px-2">
                    {state.slots.map((slotObj, i) => (
                      <div
                        key={i}
                        className={cn(state.isShaking && GAME.animations.shake)}
                        onClick={() => actions.removeDigit(i)}>
                        <ChallengeSlot
                          indexOrder={i}
                          orderValue={slotObj?.val || ""}
                          inputRef={(el) => (inputRefs.current[i] = el)}
                          onPlace={(idx, val) => {
                            const nextIdx = actions.handleManualInput(idx, val);
                            if (nextIdx !== -1 && nextIdx !== idx)
                              inputRefs.current[nextIdx]?.focus();
                            else if (nextIdx === -1)
                              inputRefs.current[idx]?.blur();
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

                  {!state.hasWon &&
                    !state.authPrompt.isActive &&
                    (!state.isComplete ? (
                      <p className="text-[10px] sm:text-xs text-foreground/40 italic">
                        {dict.hint}
                      </p>
                    ) : !isAuthenticated ? (
                      <button
                        onClick={actions.confirmAuthRedirect}
                        className={UI.guestBanner}>
                        <p className="text-xs sm:text-sm text-foreground/80 font-medium flex items-center gap-2 group-hover:text-foreground transition-colors">
                          <span className="text-brand-gold">✨</span>{" "}
                          {dict.guestBanner}
                        </p>
                      </button>
                    ) : null)}
                </>
              )}
            </div>
          </div>

          <div
            className="flex flex-col items-center justify-center w-full min-h-20 mt-4"
            style={{ opacity: state.authPrompt.isActive ? 0.1 : 1 }}>
            {state.hasWon ? (
              <ClaimPrizeForm prizeAmount={challenge.prize} dict={dict} />
            ) : state.isOutOfAttempts ? null : (
              <>
                {showTray ? (
                  <div className="token-tray flex flex-wrap justify-center gap-3 sm:gap-4 animate-in fade-in slide-in-from-bottom-4 duration-700 w-full max-w-70 sm:max-w-md mx-auto">
                    {state.trayDigits.map((digit) => (
                      <div
                        key={digit.id}
                        className="shrink-0 origin-center transition-transform hover:scale-110">
                        <DraggableNumber
                          id={digit.id}
                          value={digit.val}
                          onClick={() => actions.handleTrayClick(digit)}
                        />
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
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
              <div className="h-16 w-16 sm:h-20 sm:w-20 bg-foreground text-background flex items-center justify-center rounded-full text-3xl sm:text-4xl font-black shadow-2xl scale-110 cursor-grabbing border border-white/20">
                {activeItem.val}
              </div>
            ) : null}
          </DragOverlay>
        </DndContext>
      </div>
    </div>
  );
}

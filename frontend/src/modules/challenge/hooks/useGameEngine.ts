//femabras/frontend/src/modules/challenge/hooks/useGameEngine.ts
import { useState, useEffect, useMemo, useRef } from "react";
import { useRouter } from "next/navigation";
import { DragEndEvent, DragStartEvent, DragOverEvent } from "@dnd-kit/core";
import type { SlotItem } from "../types";
import { challengeClientService } from "../services/challenge.client.service";
import { THEME_CONFIG } from "@/shared/config/theme";
import { APIError } from "@/shared/lib/errors";
import { CHALLENGE_CONFIG } from "@/shared/config/constants";
import { challengeUtils } from "../utils/challenge.utils";
import type { Dictionary } from "@/i18n/get-dictionary";
import { env } from "@/shared/config/env";

export function useGameEngine(
  slotsCount: number,
  digits: string[],
  dict: Dictionary["challenge"],
  isAuthenticated: boolean,
) {
  const router = useRouter();
  const toastTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const [attempts, setAttempts] = useState<number>(
    CHALLENGE_CONFIG.MAX_ATTEMPTS,
  );
  const [hasWon, setHasWon] = useState(false);
  const [isShaking, setIsShaking] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [authPrompt, setAuthPrompt] = useState({
    isActive: false,
    countdown: 5,
  });
  const [slots, setSlots] = useState<SlotItem[]>(Array(slotsCount).fill(null));
  const [activeDragId, setActiveDragId] = useState<string | null>(null);
  const [focusedIndex, setFocusedIndex] = useState<number | null>(null);
  const [isOverSlot, setIsOverSlot] = useState(false);
  const [toastMsg, setToastMsg] = useState<string | null>(null);

  useEffect(() => {
    if (isAuthenticated) return;
    setAttempts(challengeClientService.getTodayAttempts());
  }, [isAuthenticated]);

  useEffect(() => {
    if (!isAuthenticated) return;

    challengeClientService.fetchLiveAttempts().then((liveAttempts) => {
      setAttempts(liveAttempts);
      challengeClientService.saveTodayAttempts(liveAttempts);
    });
  }, [isAuthenticated]);

  // Combined Auth & Restore Sequence logic
  useEffect(() => {
    if (isAuthenticated) {
      setHasWon(false);
      setIsShaking(false);
      setAuthPrompt({ isActive: false, countdown: 5 });

      const saved = sessionStorage.getItem("femabras_saved_guess");
      if (saved) {
        try {
          const parsedSlots = JSON.parse(saved);
          setSlots(parsedSlots);
          sessionStorage.removeItem("femabras_saved_guess"); // Clean up
        } catch {
          setSlots(Array(slotsCount).fill(null));
        }
      } else {
        setSlots(Array(slotsCount).fill(null));
      }
    }
  }, [isAuthenticated, slotsCount]);

  const trayDigits = useMemo(() => challengeUtils.createTray(digits), [digits]);
  const isComplete = slots.length > 0 && slots.every((slot) => slot !== null);
  const isOutOfAttempts = attempts === 0 && !hasWon;
  const title = challengeUtils.getBoardTitle(
    hasWon,
    isOutOfAttempts,
    isComplete,
    authPrompt.isActive,
    dict,
  );

  useEffect(() => {
    if (!authPrompt.isActive) return;
    if (authPrompt.countdown > 0) {
      const timer = setTimeout(
        () =>
          setAuthPrompt((prev) => ({ ...prev, countdown: prev.countdown - 1 })),
        1000,
      );
      return () => clearTimeout(timer);
    } else {
      setAuthPrompt({ isActive: false, countdown: 5 });
    }
  }, [authPrompt]);

  const showToast = (msg: string) => {
    if (toastTimeoutRef.current) clearTimeout(toastTimeoutRef.current);
    setToastMsg(msg);
    toastTimeoutRef.current = setTimeout(() => setToastMsg(null), 3000);
  };

  // The Free-Tier Realtime Pulse
  useEffect(() => {
    if (hasWon || isOutOfAttempts) return;

    const pulse = setInterval(async () => {
      try {
        const res = await fetch(`${env.apiUrl}/challenge`, {
          cache: "no-store",
        });
        if (res.ok) {
          const data = await res.json();
          if (data.status === "solved") {
            router.refresh();
          }
        }
      } catch (e) {
        if (process.env.NODE_ENV === "development") {
          console.warn("Background pulse failed to reach server:", e);
        }
      }
    }, 20000);

    return () => clearInterval(pulse);
  }, [hasWon, isOutOfAttempts, router]);

  const submitSequence = async () => {
    const guess = slots.map((s) => s?.val ?? "").join("");
    setIsSubmitting(true);
    try {
      const res = await challengeClientService.submitGuess(guess);
      setAttempts(res.remaining_attempts);
      challengeClientService.saveTodayAttempts(res.remaining_attempts);

      if (res.status === "success") setHasWon(true);
      else triggerError();
    } catch (error) {
      if (error instanceof APIError && error.status === 401) {
        setAuthPrompt({ isActive: true, countdown: 5 });
      } else if (error instanceof APIError && error.status === 403) {
        setAttempts(0);
        challengeClientService.saveTodayAttempts(0);
        triggerError();
      } else if (error instanceof APIError && error.status === 404) {
        router.refresh(); // Someone else won
      } else {
        triggerError();
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  const confirmAuthRedirect = () => {
    sessionStorage.setItem("femabras_saved_guess", JSON.stringify(slots));
    router.push("/login");
  };

  const triggerError = () => {
    setIsShaking(true);
    setTimeout(() => {
      setIsShaking(false);
      setSlots(Array(slotsCount).fill(null));
    }, THEME_CONFIG.animations.shakeDurationMs);
  };

  const handleDragStart = (e: DragStartEvent) => {
    setIsOverSlot(false);
    setActiveDragId(String(e.active.id));
  };

  const handleDragOver = (e: DragOverEvent) =>
    setIsOverSlot(!!e.over && String(e.over.id).startsWith("slot-"));

  const handleDragEnd = (e: DragEndEvent) => {
    setActiveDragId(null);
    if (
      !e.over ||
      authPrompt.isActive ||
      !String(e.over.id).startsWith("slot-")
    )
      return;
    const slotIndex = parseInt(String(e.over.id).split("-")[1]);
    const draggedDigit = trayDigits.find((d) => d.id === e.active.id);
    if (draggedDigit)
      setSlots((prev) =>
        challengeUtils.placeDigit(prev, slotIndex, draggedDigit),
      );
  };

  const handleTrayClick = (digit: { id: string; val: string }) => {
    if (authPrompt.isActive) return;
    const targetIndex = focusedIndex ?? slots.findIndex((s) => s === null);
    if (targetIndex !== -1)
      setSlots((prev) => challengeUtils.placeDigit(prev, targetIndex, digit));
  };

  const handleManualInput = (index: number, value: string): number => {
    if (authPrompt.isActive) return index;

    if (value === "") {
      removeDigit(index);
      return index;
    }

    const typed = value.slice(-1);

    if (!challengeUtils.isValidDigit(typed, digits)) {
      showToast(dict.toastInvalid.replace("{num}", typed));
      return index;
    }

    const found = { id: `manual-${Date.now()}`, val: typed };
    const nextSlots = challengeUtils.placeDigit(slots, index, found);
    setSlots(nextSlots);

    return challengeUtils.findNextEmpty(nextSlots, index);
  };

  const removeDigit = (index: number) => {
    if (hasWon || authPrompt.isActive) return;
    setSlots((prev) => challengeUtils.placeDigit(prev, index, null));
  };

  return {
    state: {
      attempts,
      hasWon,
      isShaking,
      isSubmitting,
      slots,
      activeDragId,
      authPrompt,
      trayDigits,
      toastMsg,
      isOverSlot,
      title,
      isComplete,
      isOutOfAttempts,
    },
    actions: {
      handleDragStart,
      handleDragOver,
      handleDragEnd,
      handleTrayClick,
      handleManualInput,
      removeDigit,
      submitSequence,
      confirmAuthRedirect,
      setFocusedIndex,
      setAuthPrompt,
      setAttempts,
    },
  };
}

//femabras/frontend/src/modules/challenge/utils/challenge.utils.ts
import type { Dictionary } from "@/i18n/get-dictionary";
import type { SlotItem } from "../types";

export const challengeUtils = {
  placeDigit(slots: SlotItem[], index: number, digit: SlotItem): SlotItem[] {
    const n = [...slots];
    n[index] = digit;
    return n;
  },
  isValidDigit(typed: string, available: string[]): boolean {
    return available.includes(typed);
  },
  findNextEmpty(slots: SlotItem[], afterIndex: number): number {
    const nextIdx = slots.findIndex((s, i) => i > afterIndex && s === null);
    if (nextIdx !== -1) return nextIdx;
    return slots.findIndex((s) => s === null);
  },
  createTray(digits: string[]) {
    return digits.map((val, index) => ({ id: `digit-${index}`, val }));
  },
  getBoardTitle(
    hasWon: boolean,
    outOfAttempts: boolean,
    isComplete: boolean,
    authActive: boolean,
    dict: Dictionary["challenge"],
  ) {
    if (hasWon) return dict.titleWon;
    if (outOfAttempts) return dict.titleNoAttempts;
    if (isComplete && !authActive) return dict.titleReady;
    return dict.titleDefault;
  },
};

//frontend/src/modules/challenge/types.ts
import type { Dictionary } from "@/i18n/get-dictionary";

export interface DailyChallengeResponse {
  slots: number;
  date: string;
  digits: string[];
}

export interface GuessResponse {
  status: "success" | "incorrect";
  remaining_attempts: number;
}

export type SlotItem = { id: string; val: string } | null;

export interface AuthOverlayProps {
  countdown: number;
  onConfirm: () => void;
  onCancel: () => void;
  dict: Dictionary["challenge"];
}

export interface GameBoardProps {
  challenge: DailyChallengeResponse;
  isAuthenticated: boolean;
  dict: Dictionary["challenge"];
}

export interface ChallengeSlotProps {
  indexOrder: number;
  orderValue: string;
  onPlace: (indexOrder: number, orderValue: string) => void;
  onKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void;
  inputRef: (element: HTMLInputElement | null) => void;
  onFocus: () => void;
  onBlur: () => void;
  isShaking?: boolean;
}

export interface DraggableNumberProps {
  id: string;
  value: string;
  onClick?: () => void;
  classNameExtra?: string;
}

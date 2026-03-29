//femabras/frontend/src/modules/challenge/components/challenge-slot.client.tsx
"use client";

import { useDroppable } from "@dnd-kit/core";
import { GAME } from "@/shared/config/gameStyles";

interface SlotProps {
  indexOrder: number;
  orderValue: string;
  onPlace: (indexOrder: number, orderValue: string) => void;
  onKeyDown: (key: React.KeyboardEvent) => void;
  inputRef: (element: HTMLInputElement | null) => void;
  onFocus: () => void;
  onBlur: () => void;
  isShaking?: boolean;
}

export function ChallengeSlot({
  indexOrder,
  orderValue,
  onPlace,
  onKeyDown,
  inputRef,
  onFocus,
  onBlur,
  isShaking = false,
}: SlotProps) {
  const { isOver, setNodeRef } = useDroppable({ id: `slot-${indexOrder}` });
  const hasValue = orderValue !== "";

  return (
    <div ref={setNodeRef} className={GAME.slot(isShaking, isOver, hasValue)}>
      <input
        ref={inputRef}
        type="text"
        inputMode="numeric"
        pattern="[0-9]*"
        maxLength={1}
        value={orderValue}
        onKeyDown={onKeyDown}
        onFocus={(e) => {
          e.target.select();
          onFocus();
        }}
        onBlur={onBlur}
        onClick={(e) => e.stopPropagation()}
        onChange={(e) => onPlace(indexOrder, e.target.value)}
        className={`h-full w-full bg-transparent text-center text-2xl sm:text-4xl font-bold uppercase outline-none animate-pop-in ${
          hasValue && !isShaking ? "text-background" : "text-foreground"
        }`}
        placeholder="?"
      />
      <div className="absolute inset-0 pointer-events-none" />
    </div>
  );
}

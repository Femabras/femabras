//femabras/frontend/src/modules/challenge/components/challenge-slot.client.tsx
"use client";

import { useDroppable } from "@dnd-kit/core";
import { GAME } from "@/shared/config/gameStyles";
import { cn } from "@/shared/lib/utils";

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
    <div
      ref={setNodeRef}
      className={cn(
        GAME.slot(isShaking, isOver, hasValue),
        // 🟢 FIX: Fluid width on mobile (14vw), fixed max width on desktop.
        // This stops them from stretching horizontally while maintaining a perfect rectangle.
        "w-[14vw] sm:w-16 md:w-20 h-auto aspect-[3/4] p-0 flex items-center justify-center relative shrink-0",
      )}>
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
        className={cn(
          "absolute inset-0 h-full w-full bg-transparent text-center font-bold uppercase outline-none animate-pop-in z-10",
          "text-[clamp(1.5rem,6vw,2.5rem)]",
          hasValue && !isShaking ? "text-background" : "text-foreground",
        )}
        placeholder="?"
      />
      <div className="absolute inset-0 pointer-events-none z-0" />
    </div>
  );
}

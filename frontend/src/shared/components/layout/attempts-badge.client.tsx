// femabras/frontend/src/shared/components/layout/attempts-badge.client.tsx
"use client";

import { useState, useEffect, useRef } from "react";
import { cn } from "@/shared/lib/utils";
import { ATTEMPTS_EVENT } from "@/shared/lib/events";

interface AttemptsBadgeProps {
  initialAttempts: number;
  /** i18n label — pass dict.layout.attempts from the server component */
  label: string;
  containerClassName?: string;
  textClassName?: string;
}

export function AttemptsBadge({
  initialAttempts,
  label,
  containerClassName,
  textClassName,
}: AttemptsBadgeProps) {
  const [attempts, setAttempts] = useState(initialAttempts);
  const [isAnimating, setIsAnimating] = useState(false);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const handleAttemptsChange = (e: Event) => {
      const customEvent = e as CustomEvent<number>;
      setAttempts(customEvent.detail);
      setIsAnimating(true);
      if (timerRef.current) clearTimeout(timerRef.current);
      timerRef.current = setTimeout(() => setIsAnimating(false), 300);
    };

    window.addEventListener(ATTEMPTS_EVENT, handleAttemptsChange);
    return () => {
      window.removeEventListener(ATTEMPTS_EVENT, handleAttemptsChange);
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, []);

  return (
    <div
      className={cn(
        "flex items-center px-3 py-1.5 bg-brand-gold/10 border border-brand-gold/20 rounded-full backdrop-blur-md shadow-[0_0_15px_rgba(234,179,8,0.1)]",
        containerClassName,
      )}>
      <span
        className={cn(
          "flex h-5 w-5 items-center justify-center rounded-full bg-brand-gold text-background text-[10px] font-black transition-transform duration-300",
          isAnimating ? "scale-125" : "scale-100",
        )}>
        {attempts}
      </span>
      <span
        className={cn(
          "text-[10px] font-bold uppercase tracking-widest text-brand-gold",
          textClassName,
        )}>
        {label}
      </span>
    </div>
  );
}

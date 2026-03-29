//femabras/frontend/src/modules/challenge/components/auth-overlay.client.tsx
"use client";

import { Button } from "@/shared/components/ui/button";
import type { AuthOverlayProps } from "../types";
import { UI } from "../utils/styles";

export function AuthOverlay({
  countdown,
  onConfirm,
  onCancel,
  dict,
}: AuthOverlayProps) {
  return (
    <div className={UI.authOverlay}>
      <h3 className="text-xl sm:text-2xl font-bold text-yellow-400 mb-2">
        {dict.authTitle}
      </h3>
      <div className="text-7xl sm:text-8xl font-black text-foreground mb-4 animate-pulse tracking-tighter">
        {countdown}
        <span className="text-4xl text-foreground/40">s</span>
      </div>
      <p className="text-center text-foreground/80 mb-8 max-w-sm text-sm sm:text-base leading-relaxed">
        {dict.authBody}
      </p>

      <Button variant="warning" onClick={onConfirm} className="max-w-xs mb-4">
        {dict.authConfirm}
      </Button>

      <button
        onClick={onCancel}
        className="text-xs uppercase tracking-widest opacity-50 hover:opacity-100 transition-opacity">
        {dict.authCancel}
      </button>
    </div>
  );
}

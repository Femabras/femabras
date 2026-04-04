//femabras/frontend/src/modules/challenge/components/winner-status.client.tsx
"use client";

import { cn } from "@/shared/lib/utils";
import type { Dictionary } from "@/i18n/get-dictionary";

export function WinnerStatus({
  prize,
  payoutStatus,
  dict,
}: {
  prize: number;
  payoutStatus: "unclaimed" | "pending" | "paid" | "rejected";
  dict: Dictionary["challenge"];
}) {
  const steps = [
    {
      label: dict.statusCracked || "Code Cracked",
      active: true,
    },
    {
      label: dict.statusClaimed || "Claim Submitted",
      active:
        payoutStatus === "pending" ||
        payoutStatus === "paid" ||
        payoutStatus === "rejected",
    },
    {
      label:
        payoutStatus === "rejected"
          ? dict.statusRejected || "Issue with Payout"
          : dict.statusPaid || "Prize Sent",
      active: payoutStatus === "paid" || payoutStatus === "rejected",
      isError: payoutStatus === "rejected",
    },
  ];

  return (
    <div className="w-full max-w-md mx-auto mt-8 p-8 bg-background border border-yellow-500/30 rounded-3xl shadow-[0_0_40px_rgba(250,204,21,0.1)] relative overflow-hidden animate-in zoom-in-95 duration-700">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-1/2 bg-yellow-500/10 blur-3xl rounded-full pointer-events-none" />

      <div className="text-center relative z-10 mb-8">
        <h2 className="text-3xl font-black text-foreground tracking-tight mb-2">
          {dict.youWonTitle || "You Won Today!"}
        </h2>
        <p className="text-sm text-foreground/70">
          {dict.youWonSubtitle || "Your prize progress is tracked below."}
        </p>

        <div className="mt-6 inline-block bg-yellow-500/10 border border-yellow-500/50 rounded-2xl px-6 py-3">
          <p className="text-[10px] uppercase font-bold text-yellow-500 tracking-widest mb-1">
            {dict.prizeLocked}
          </p>
          <p className="text-3xl font-black text-yellow-400 drop-shadow-[0_0_10px_rgba(250,204,21,0.5)]">
            {new Intl.NumberFormat("en", { minimumFractionDigits: 2 }).format(
              prize,
            )}{" "}
            <span className="text-lg">AOA</span>
          </p>
        </div>
      </div>

      <div className="relative z-10 space-y-6 pl-4">
        {steps.map((step, index) => (
          <div key={index} className="flex items-start gap-4 relative">
            {index < steps.length - 1 && (
              <div
                className={cn(
                  "absolute left-3.5 top-8 w-0.5 h-full -ml-px transition-colors duration-500",
                  steps[index + 1].active ? "bg-yellow-500" : "bg-white/10",
                )}
              />
            )}

            <div
              className={cn(
                "relative z-10 flex h-7 w-7 items-center justify-center rounded-full border-2 transition-all duration-500",
                step.active && !step.isError
                  ? "border-yellow-500 bg-yellow-500 text-background"
                  : step.isError
                    ? "border-red-500 bg-red-500 text-white"
                    : "border-white/20 bg-background",
              )}>
              {step.active && !step.isError ? "✓" : step.isError ? "!" : ""}
            </div>

            <div className="pt-0.5">
              <p
                className={cn(
                  "font-bold transition-colors duration-500",
                  step.active ? "text-foreground" : "text-foreground/40",
                  step.isError && "text-red-400",
                )}>
                {step.label}
              </p>
              {index === 1 && step.active && payoutStatus === "pending" && (
                <p className="text-xs text-yellow-500 mt-1 animate-pulse">
                  {dict.statusReviewing}
                </p>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

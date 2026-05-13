// femabras/frontend/src/modules/challenge/components/claim-prize-form.client.tsx
//
// Updated to use apiFetch which automatically attaches the CSRF token.
// The form also validates input lengths client-side to mirror the backend's
// new strict ClaimRequest validation rules.

"use client";

import { useState } from "react";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { Button } from "@/shared/components/ui/button";
import { env } from "@/shared/config/env";
import { apiFetch } from "@/shared/lib/api.client";
import type { Dictionary } from "@/i18n/get-dictionary";
import { cn } from "@/shared/lib/utils";

export function ClaimPrizeForm({
  prizeAmount,
  dict,
}: {
  prizeAmount: number;
  dict: Dictionary["challenge"];
}) {
  const [method, setMethod] = useState<"Bank" | "ATM" | "Multicaixa">("Bank");
  const [destination, setDestination] = useState("");
  const [accountName, setAccountName] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const atmEligible =
    prizeAmount >= 1000 && prizeAmount <= 30000 && prizeAmount % 1000 === 0;

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    // Client-side validation mirrors backend ClaimRequest binding tags
    if (destination.trim().length < 5 || destination.trim().length > 50) {
      setErrorMsg("Destination must be between 5 and 50 characters.");
      return;
    }
    if (accountName.trim().length < 2 || accountName.trim().length > 100) {
      setErrorMsg("Account name must be between 2 and 100 characters.");
      return;
    }

    setIsSubmitting(true);
    setErrorMsg(null);

    try {
      const res = await apiFetch(`${env.apiUrl}/claim`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          method,
          destination: destination.trim(),
          account_name: accountName.trim(),
        }),
      });

      if (!res.ok) {
        const err = await res.json().catch(() => null);
        throw new Error(err?.error || "Failed to submit claim.");
      }

      setSuccess(true);
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : "An error occurred.");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (success) {
    return (
      <div className="w-full max-w-md mx-auto mt-8 p-8 bg-background border border-yellow-500/30 rounded-3xl text-center animate-in zoom-in-95 duration-500">
        <div className="text-5xl mb-4">🎉</div>
        <h3 className="text-2xl font-black text-foreground mb-2">
          {dict.claimSuccessTitle}
        </h3>
        <p className="text-sm text-foreground/70">{dict.claimSuccessMsg}</p>
      </div>
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="w-full max-w-md mx-auto mt-8 p-6 sm:p-8 bg-white/5 border border-white/10 rounded-3xl backdrop-blur-md flex flex-col gap-5 animate-in fade-in zoom-in-95 duration-500">
      <div className="text-center">
        <h3 className="text-2xl font-black text-foreground mb-1">
          {dict.claimTitle}
        </h3>
        <p className="text-xs text-foreground/60">{dict.claimSubtitle}</p>
        <p className="mt-3 text-3xl font-black text-yellow-400 drop-shadow-[0_0_10px_rgba(250,204,21,0.5)]">
          {new Intl.NumberFormat("en", { minimumFractionDigits: 2 }).format(
            prizeAmount,
          )}{" "}
          <span className="text-base">AOA</span>
        </p>
      </div>

      {errorMsg && (
        <div className="p-3 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200 text-xs text-center">
          {errorMsg}
        </div>
      )}

      <div className="flex flex-col gap-2">
        <Label htmlFor="method">{dict.claimMethodLabel}</Label>
        <div className="grid grid-cols-3 gap-2">
          {(["Bank", "ATM", "Multicaixa"] as const).map((opt) => {
            const disabled = opt === "ATM" && !atmEligible;
            return (
              <button
                key={opt}
                type="button"
                disabled={disabled}
                onClick={() => setMethod(opt)}
                className={cn(
                  "py-3 px-2 rounded-xl text-xs font-bold transition-all",
                  method === opt
                    ? "bg-brand-gold text-background"
                    : "bg-white/5 text-foreground/70 hover:bg-white/10",
                  disabled && "opacity-30 cursor-not-allowed",
                )}>
                {opt === "Bank"
                  ? dict.claimBank
                  : opt === "ATM"
                    ? dict.claimAtm
                    : dict.claimExpress}
              </button>
            );
          })}
        </div>
        {!atmEligible && (
          <p className="text-[10px] text-foreground/40 italic">
            {dict.claimAtmDisabled}
          </p>
        )}
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor="destination">
          {method === "Bank" ? dict.claimIbanLabel : dict.claimPhoneLabel}
        </Label>
        <Input
          id="destination"
          name="destination"
          required
          minLength={5}
          maxLength={50}
          value={destination}
          onChange={(e) => setDestination(e.target.value)}
          disabled={isSubmitting}
        />
      </div>

      <div className="flex flex-col gap-2">
        <Label htmlFor="account_name">{dict.claimAccountName}</Label>
        <Input
          id="account_name"
          name="account_name"
          required
          minLength={2}
          maxLength={100}
          value={accountName}
          onChange={(e) => setAccountName(e.target.value)}
          disabled={isSubmitting}
        />
      </div>

      <Button
        type="submit"
        variant="warning"
        className="w-full mt-2"
        disabled={isSubmitting}>
        {isSubmitting ? dict.btnProcessing : dict.btnClaim}
      </Button>
    </form>
  );
}
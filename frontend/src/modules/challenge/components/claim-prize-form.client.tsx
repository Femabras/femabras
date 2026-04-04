//femabras/frontend/src/modules/challenge/components/claim-prize-form.client.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { env } from "@/shared/config/env";
import type { Dictionary } from "@/i18n/get-dictionary";

export function ClaimPrizeForm({
  prizeAmount,
  dict,
}: {
  prizeAmount: number;
  dict: Dictionary["challenge"];
}) {
  const router = useRouter();

  const isAtmEligible =
    prizeAmount >= 1000 && prizeAmount <= 30000 && prizeAmount % 1000 === 0;

  const [method, setMethod] = useState<string>("Multicaixa Express");
  const [destination, setDestination] = useState("");
  const [accountName, setAccountName] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrorMsg(null);

    try {
      const res = await fetch(`${env.apiUrl}/claim`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          method,
          destination,
          account_name: method === "Bank Transfer" ? accountName : "",
        }),
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Failed to submit claim.");
      }

      router.refresh();

      setSuccess(true);
    } catch (error) {
      setErrorMsg(
        error instanceof Error ? error.message : "An error occurred.",
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  // 🟢 The Glassmorphism Overlay Wrapper
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-in fade-in duration-700">
      {success ? (
        <div className="p-8 bg-background border border-green-500/50 rounded-3xl text-center animate-in zoom-in duration-500 shadow-[0_0_40px_rgba(34,197,94,0.2)] max-w-sm w-full">
          <div className="text-5xl mb-4">💸</div>
          <h3 className="text-2xl font-bold text-green-400 mb-2">
            {dict.claimSuccessTitle}
          </h3>
          <p className="text-sm text-foreground/80">{dict.claimSuccessMsg}</p>
        </div>
      ) : (
        <div className="w-full max-w-md bg-background border border-white/10 rounded-3xl p-6 sm:p-8 animate-in zoom-in-95 duration-500 shadow-2xl relative overflow-hidden">
          {/* Subtle glowing background effect inside the modal */}
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-1/2 bg-yellow-500/10 blur-3xl rounded-full pointer-events-none" />

          <div className="text-center mb-6 relative z-10">
            <h3 className="text-2xl font-black text-yellow-400 mb-1 drop-shadow-[0_0_10px_rgba(250,204,21,0.4)]">
              {dict.claimTitle}
            </h3>
            <p className="text-sm text-foreground/70">{dict.claimSubtitle}</p>
          </div>

          {errorMsg && (
            <div className="mb-4 p-3 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200 text-xs text-center relative z-10">
              {errorMsg}
            </div>
          )}

          <form
            onSubmit={handleSubmit}
            className="space-y-4 text-left relative z-10">
            <div className="space-y-2">
              <Label htmlFor="method">{dict.claimMethodLabel}</Label>
              <select
                id="method"
                value={method}
                onChange={(e) => {
                  setMethod(e.target.value);
                  setDestination(""); // Clear input when switching methods
                }}
                className="flex h-12 w-full rounded-xl border border-white/20 bg-white/5 px-4 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-yellow-400/50 text-foreground transition-all">
                <option value="Multicaixa Express" className="bg-background">
                  {dict.claimExpress || "Multicaixa Express"}
                </option>
                <option value="Unitel Money" className="bg-background">
                  Unitel Money
                </option>
                <option value="Bank Transfer" className="bg-background">
                  {dict.claimBank}
                </option>
                <option
                  value="ATM"
                  disabled={!isAtmEligible}
                  className="bg-background disabled:text-foreground/30">
                  {dict.claimAtm} {!isAtmEligible && dict.claimAtmDisabled}
                </option>
              </select>
            </div>

            {/* 🟢 Conditionally render Account Name field for Bank Transfers */}
            {method === "Bank Transfer" && (
              <div className="space-y-2 animate-in slide-in-from-top-2 fade-in">
                <Label htmlFor="accountName">{dict.claimAccountName}</Label>
                <Input
                  id="accountName"
                  required
                  placeholder="Fernando Bras"
                  value={accountName}
                  onChange={(e) => setAccountName(e.target.value)}
                  disabled={isSubmitting}
                  className="h-12 rounded-xl bg-white/5 border-white/20"
                />
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="destination">
                {method === "Bank Transfer"
                  ? dict.claimIbanLabel
                  : dict.claimPhoneLabel}
              </Label>
              <Input
                id="destination"
                required
                placeholder={
                  method === "Bank Transfer"
                    ? "AO06 0000 0000 0000 0000 0000 0"
                    : "9XX XXX XXX"
                }
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
                disabled={isSubmitting}
                className="h-12 rounded-xl bg-white/5 border-white/20"
              />
            </div>

            <Button
              type="submit"
              variant="warning"
              className="w-full h-12 font-bold tracking-widest mt-4 rounded-xl text-black hover:scale-[1.02] transition-transform"
              disabled={isSubmitting}>
              {isSubmitting ? dict.btnProcessing : dict.btnClaim}
            </Button>
          </form>
        </div>
      )}
    </div>
  );
}

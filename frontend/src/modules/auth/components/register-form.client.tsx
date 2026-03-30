//femabras/frontend/src/modules/auth/components/register-form.client.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { authClientService } from "../services/auth.client.service";
import type { Dictionary } from "@/i18n/get-dictionary";

export function RegisterForm({ dict }: { dict: Dictionary["auth"] }) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const [step, setStep] = useState<"register" | "otp">("register");
  const [pendingUserId, setPendingUserId] = useState<string | null>(null);

  const handleRegister = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    setErrorMsg(null);

    const formData = new FormData(e.currentTarget);
    const email = formData.get("email") as string;
    const phone = formData.get("phone") as string;
    const password = formData.get("password") as string;

    try {
      const res = await authClientService.register(email, phone, password);
      setPendingUserId(String(res.user_id));
      setStep("otp");
    } catch (error) {
      setErrorMsg(
        error instanceof Error ? error.message : "Registration failed.",
      );
    } finally {
      setIsLoading(false);
    }
  };

  // FIXED: Using React.SyntheticEvent
  const handleVerifyOTP = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!pendingUserId) return;

    setIsLoading(true);
    setErrorMsg(null);
    const otp = new FormData(e.currentTarget).get("otp") as string;

    try {
      await authClientService.verifyOTP(pendingUserId, otp);
      router.push("/");
      router.refresh();
    } catch (error) {
      setErrorMsg(error instanceof Error ? error.message : "Invalid OTP.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md mx-auto p-6 sm:p-8 bg-white/5 border border-white/10 rounded-3xl shadow-2xl backdrop-blur-xl animate-in fade-in zoom-in-95 duration-500">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-black tracking-tight text-foreground mb-2">
          {step === "register" ? dict.registerTitle : dict.verifyTitle}
        </h1>
        <p className="text-sm text-foreground/60">
          {step === "register" ? dict.registerSubtitle : dict.verifySubtitle}
        </p>
      </div>

      {errorMsg && (
        <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200 text-sm text-center font-medium">
          {errorMsg}
        </div>
      )}

      {step === "register" ? (
        <form onSubmit={handleRegister} className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="email">{dict.emailLabel}</Label>
            <Input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              placeholder={dict.emailPlaceholder}
              required
              disabled={isLoading}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="phone">{dict.phoneLabel}</Label>
            <Input
              id="phone"
              name="phone"
              type="tel"
              autoComplete="tel"
              placeholder={dict.phonePlaceholder}
              required
              disabled={isLoading}
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="password">{dict.passwordLabel}</Label>
            <Input
              id="password"
              name="password"
              type="password"
              autoComplete="new-password"
              placeholder={dict.passwordMin}
              required
              minLength={8}
              disabled={isLoading}
            />
          </div>
          <Button
            type="submit"
            variant="warning"
            className="w-full mt-4"
            disabled={isLoading}>
            {isLoading ? dict.processing : dict.btnRegister}
          </Button>
        </form>
      ) : (
        <form onSubmit={handleVerifyOTP} className="space-y-5">
          <div className="space-y-2 text-center">
            <Input
              id="otp"
              name="otp"
              type="text"
              inputMode="numeric"
              autoComplete="one-time-code"
              maxLength={6}
              placeholder="000000"
              className="text-center text-2xl tracking-widest font-black"
              required
              disabled={isLoading}
              autoFocus
            />
          </div>
          <Button
            type="submit"
            variant="warning"
            className="w-full mt-4"
            disabled={isLoading}>
            {isLoading ? dict.btnVerifying : dict.btnVerify}
          </Button>
        </form>
      )}

      {step === "register" && (
        <div className="mt-8 text-center text-sm text-foreground/60">
          {dict.haveAccount}{" "}
          <Link
            href="/login"
            className="text-yellow-400 font-bold hover:underline transition-all">
            {dict.loginLink}
          </Link>
        </div>
      )}
    </div>
  );
}

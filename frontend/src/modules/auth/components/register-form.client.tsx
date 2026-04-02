//femabras/frontend/src/modules/auth/components/register-form.client.tsx
"use client";

import { useState, useMemo } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { authClientService } from "../services/auth.client.service";
import type { Dictionary } from "@/i18n/get-dictionary";
import { zxcvbn, zxcvbnOptions } from "@zxcvbn-ts/core";
import * as zxcvbnCommonPackage from "@zxcvbn-ts/language-common";
import * as zxcvbnEnPackage from "@zxcvbn-ts/language-en";

zxcvbnOptions.setOptions({
  translations: zxcvbnEnPackage.translations,
  graphs: zxcvbnCommonPackage.adjacencyGraphs,
  dictionary: {
    ...zxcvbnCommonPackage.dictionary,
    ...zxcvbnEnPackage.dictionary,
  },
});

export function RegisterForm({ dict }: { dict: Dictionary["auth"] }) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const [step, setStep] = useState<"register" | "otp">("register");
  const [pendingUserId, setPendingUserId] = useState<string | null>(null);

  const [password, setPassword] = useState("");

  const handleRegister = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    setErrorMsg(null);

    const formData = new FormData(e.currentTarget);
    const name = formData.get("name") as string;
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
      const res = await authClientService.register(name, email, password);
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

  const passwordStrength = useMemo(() => {
    if (!password) return { score: -1, warning: "", color: "bg-white/5" };
    const result = zxcvbn(password);
    const colors = [
      "bg-red-500/70",
      "bg-red-400",
      "bg-yellow-500/80",
      "bg-yellow-400",
      "bg-foreground shadow-[0_0_10px_rgba(251,255,254,0.6)]",
    ];
    return {
      score: result.score,
      warning: result.feedback.warning || result.feedback.suggestions[0] || "",
      color: colors[result.score],
    };
  }, [password]);

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
            <Label htmlFor="name">{dict.nameLabel}</Label>
            <Input
              id="name"
              name="name"
              type="text"
              autoComplete="username"
              placeholder={dict.namePlaceholder}
              required
              minLength={2}
              maxLength={30}
              disabled={isLoading}
            />
          </div>
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
            <Label htmlFor="password">{dict.passwordLabel}</Label>
            <Input
              id="password"
              name="password"
              type="password"
              autoComplete="new-password"
              placeholder={dict.passwordMin}
              required
              minLength={16}
              disabled={isLoading}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />

            {/* Render the Strength Meter UI */}
            <div className="pt-3">
              <div className="flex gap-2 w-full">
                {[0, 1, 2, 3].map((threshold) => {
                  const isActive =
                    passwordStrength.score === 0
                      ? threshold === 1
                      : passwordStrength.score >= threshold;
                  return (
                    <div
                      key={threshold}
                      className={`h-1.5 w-1/4 rounded-full transition-all duration-500 ease-out ${
                        isActive
                          ? passwordStrength.color
                          : "bg-white/5 shadow-inner"
                      }`}
                    />
                  );
                })}
              </div>
              {passwordStrength.warning && (
                <p className="text-[11px] text-foreground/70 mt-3 font-medium animate-in fade-in slide-in-from-top-1">
                  💡 {passwordStrength.warning}
                </p>
              )}
            </div>
          </div>
          <Button
            type="submit"
            variant="warning"
            className="w-full mt-4"
            disabled={
              isLoading || (password.length > 0 && passwordStrength.score < 2)
            }>
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
              placeholder=""
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

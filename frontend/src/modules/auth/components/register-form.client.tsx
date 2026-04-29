// femabras/frontend/src/modules/auth/components/register-form.client.tsx
"use client";

import { useState, useMemo } from "react";
import { useRouter, useParams } from "next/navigation";
import Link from "next/link";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { authClientService } from "../services/auth.client.service";
import { env } from "@/shared/config/env";
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

// Inline Google SVG so it renders identically on both forms
function GoogleIcon() {
  return (
    <svg className="w-5 h-5 shrink-0" viewBox="0 0 24 24">
      <path
        fill="#4285F4"
        d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
      />
      <path
        fill="#34A853"
        d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
      />
      <path
        fill="#FBBC05"
        d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
      />
      <path
        fill="#EA4335"
        d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
      />
    </svg>
  );
}

export function RegisterForm({ dict }: { dict: Dictionary["auth"] }) {
  const router = useRouter();
  const params = useParams();
  const locale = (params?.locale as string) || "en";

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
      router.push(`/${locale}`);
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
        <>
          {/* Google OAuth sign-up */}
          <a
            href={`${env.apiUrl}/auth/google/login`}
            className="flex items-center justify-center gap-3 w-full py-3 px-4 mb-6 bg-white/5 hover:bg-white/10 border border-white/10 hover:border-white/20 rounded-2xl transition-all text-sm font-bold text-foreground">
            <GoogleIcon />
            Continue with Google
          </a>

          <div className="flex items-center gap-3 mb-6">
            <div className="flex-1 h-px bg-white/10" />
            <span className="text-xs text-foreground/40 font-medium uppercase tracking-widest">
              or
            </span>
            <div className="flex-1 h-px bg-white/10" />
          </div>

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

          <div className="mt-8 text-center text-sm text-foreground/60">
            {dict.haveAccount}{" "}
            <Link
              href={`/${locale}/login`}
              className="text-yellow-400 font-bold hover:underline transition-all">
              {dict.loginLink}
            </Link>
          </div>
        </>
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
    </div>
  );
}

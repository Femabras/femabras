//femabras/frontend/src/modules/auth/components/login-form.client.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { authClientService } from "../services/auth.client.service";
import type { Dictionary } from "@/i18n/get-dictionary";

export function LoginForm({ dict }: { dict: Dictionary["auth"] }) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleSubmit = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);
    setErrorMsg(null);

    const formData = new FormData(e.currentTarget);
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
      await authClientService.login(email, password);

      const savedGuess = sessionStorage.getItem("femabras_saved_guess");
      if (savedGuess) {
        console.log("Found saved guess:", JSON.parse(savedGuess));
      }

      router.push("/");
      router.refresh();
    } catch (error) {
      setErrorMsg(
        error instanceof Error ? error.message : "Invalid email or password.",
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md mx-auto p-6 sm:p-8 bg-white/5 border border-white/10 rounded-3xl shadow-2xl backdrop-blur-xl animate-in fade-in zoom-in-95 duration-500">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-black tracking-tight text-foreground mb-2">
          {dict.loginTitle}
        </h1>
        <p className="text-sm text-foreground/60">{dict.loginSubtitle}</p>
      </div>

      {errorMsg && (
        <div className="mb-6 p-4 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200 text-sm text-center font-medium">
          {errorMsg}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5">
        <div className="space-y-2">
          <Label htmlFor="email">{dict.emailLabel}</Label>
          <Input
            id="email"
            name="email"
            type="text"
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
            placeholder={dict.passwordPlaceholder}
            required
            disabled={isLoading}
          />
        </div>

        <Button
          type="submit"
          variant="warning"
          className="w-full mt-4"
          disabled={isLoading}>
          {isLoading ? dict.processing : dict.loginBtn}
        </Button>
      </form>

      <div className="mt-8 text-center text-sm text-foreground/60">
        {dict.noAccount}{" "}
        <Link
          href="/register"
          className="text-yellow-400 font-bold hover:underline transition-all">
          {dict.registerLink}
        </Link>
      </div>
    </div>
  );
}

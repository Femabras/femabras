//femabras/frontend/src/app/[locale]/error.tsx
"use client";

import { useEffect } from "react";
import { useParams } from "next/navigation";

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

const errorTranslations = {
  en: { title: "Something went wrong!", retry: "Try again" },
  pt: { title: "Algo deu errado!", retry: "Tentar novamente" },
  fr: { title: "Quelque chose s'est mal passé!", retry: "Réessayer" },
};

export default function ErrorPage({ error, reset }: ErrorProps) {
  const params = useParams();
  const locale = (params?.locale as string) || "en";

  const dict =
    errorTranslations[locale as keyof typeof errorTranslations] ||
    errorTranslations.en;

  useEffect(() => {
    console.error("Application Error Caught by Boundary:", error);
  }, [error]);

  return (
    <div className="flex flex-1 flex-col items-center justify-center w-full px-4 py-12 relative">
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-red-500/10 rounded-full blur-[100px] pointer-events-none" />

      <div className="flex flex-col items-center gap-6 text-center z-10 w-full max-w-md">
        <h2 className="text-3xl sm:text-5xl font-black text-foreground tracking-tight leading-tight">
          {dict.title}
        </h2>

        <p className="text-foreground/50 text-sm">
          {error.message || "An unexpected error occurred."}
        </p>

        <button
          onClick={() => reset()}
          className="w-full sm:w-auto px-10 py-4 font-black rounded-2xl uppercase tracking-widest bg-brand-gold text-background shadow-[0_0_15px_var(--color-brand-gold-glow)] hover:scale-105 transition-all active:scale-95 mt-4">
          {dict.retry}
        </button>
      </div>
    </div>
  );
}

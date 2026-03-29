//femabras/frontend/src/app/[locale]/error.tsx
"use client";

import { useEffect } from "react";
import { Button } from "@/shared/components/ui/button";

const errorTranslations = {
  en: {
    message: "Something went wrong.",
    tryAgain: "Try again",
    defaultError: "An unexpected error occurred.",
  },
  pt: {
    message: "Algo deu errado.",
    tryAgain: "Tentar novamente",
    defaultError: "Ocorreu um erro inesperado.",
  },
  fr: {
    message: "Quelque chose s'est mal passé.",
    tryAgain: "Réessayer",
    defaultError: "Une erreur inattendue s'est produite.",
  },
};

interface ErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
  params: { locale: string };
}

export default function ErrorPage({ error, reset, params }: ErrorProps) {
  const dict =
    errorTranslations[params.locale as keyof typeof errorTranslations] ||
    errorTranslations.en;

  useEffect(() => {
    console.error("App error:", error);
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8 text-center">
      <h2 className="mb-4 text-3xl font-bold text-red-400">{dict.message}</h2>
      <p className="mb-8 text-foreground/70 max-w-md">
        {error.message || dict.defaultError}
      </p>

      <Button variant="warning" onClick={reset} className="px-8">
        {dict.tryAgain}
      </Button>

      <p className="mt-6 text-xs opacity-50">Error ID: {error.digest}</p>
    </div>
  );
}

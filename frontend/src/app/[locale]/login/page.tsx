//femabras/frontend/src/app/[locale]/login/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";
import { LoginForm } from "@/modules/auth/components/login-form.client";
import Link from "next/link";

export default async function LoginPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const dict = await getDictionary(locale);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4 relative">
      {/* Subtle background glow effect */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-yellow-500/10 rounded-full blur-[100px] pointer-events-none" />

      <div className="w-full max-w-md mb-8 z-10">
        <Link
          href="/"
          className="text-sm font-bold text-foreground/50 hover:text-yellow-400 transition-colors">
          {dict.auth.backToGame}
        </Link>
      </div>

      <div className="z-10 w-full">
        <LoginForm dict={dict.auth} />
      </div>
    </div>
  );
}

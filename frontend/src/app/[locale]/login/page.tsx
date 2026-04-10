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
    <div className="flex flex-1 flex-col items-center justify-center w-full px-4 py-12 relative">
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-brand-gold/10 rounded-full blur-[100px] pointer-events-none" />

      <div className="w-full max-w-md flex flex-col gap-6 z-10">
        <Link
          href="/"
          className="text-sm font-bold text-foreground/50 hover:text-brand-gold transition-colors">
          {dict.auth.backToGame}
        </Link>
        <LoginForm dict={dict.auth} />
      </div>
    </div>
  );
}

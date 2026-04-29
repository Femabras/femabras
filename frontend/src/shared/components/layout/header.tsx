// femabras/frontend/src/shared/components/layout/header.tsx
import Link from "next/link";
import Image from "next/image";
import { cookies } from "next/headers";
import { LanguageSwitcher } from "./language-switcher.client";
import { MobileMenu } from "./mobile-menu.client";
import { LogoutButton } from "./logout-button.client";
import { AttemptsBadge } from "./attempts-badge.client";
import { challengeServerService } from "@/modules/challenge/services/challenge.server.service";
import type { Dictionary } from "@/i18n/get-dictionary";

interface HeaderProps {
  dict: Dictionary["layout"];
  locale: string;
}

export async function Header({ dict, locale }: HeaderProps) {
  const cookieStore = await cookies();
  const isAuthenticated = cookieStore.has("access_token");

  let initialAttempts = 5;
  if (isAuthenticated) {
    initialAttempts = await challengeServerService.getAttempts();
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b border-white/10 bg-background/80 backdrop-blur-xl">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16 sm:h-20 relative">
          <div className="flex items-center z-10">
            <Link href={`/${locale}`} className="flex items-center gap-2 group">
              <Image
                src="/logos/femabras-logo-light.svg"
                alt="Femabras"
                width={140}
                height={32}
                priority
                className="w-auto h-6 sm:h-8 group-hover:opacity-80 transition-opacity drop-shadow-[0_0_10px_var(--color-brand-gold-glow)]"
              />
            </Link>
          </div>

          {/* Mobile: centered attempts badge */}
          {isAuthenticated && (
            <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 flex lg:hidden items-center z-0 animate-in fade-in zoom-in duration-300">
              <AttemptsBadge
                initialAttempts={initialAttempts}
                label={dict.attempts}
                containerClassName="gap-1.5"
                textClassName="hidden xs:inline-block"
              />
            </div>
          )}

          <nav className="hidden lg:flex absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 items-center gap-8 z-0">
            <Link
              href={`/${locale}`}
              className="text-xs font-bold uppercase tracking-widest text-foreground/70 hover:text-brand-gold transition-colors">
              {dict.home}
            </Link>
            <Link
              href={`/${locale}/products`}
              className="text-xs font-bold uppercase tracking-widest text-foreground/70 hover:text-brand-gold transition-colors">
              {dict.products}
            </Link>
            <Link
              href={`/${locale}/services`}
              className="text-xs font-bold uppercase tracking-widest text-foreground/70 hover:text-brand-gold transition-colors">
              {dict.services}
            </Link>
            <Link
              href={`/${locale}/about`}
              className="text-xs font-bold uppercase tracking-widest text-foreground/70 hover:text-brand-gold transition-colors">
              {dict.about}
            </Link>
            <Link
              href={`/${locale}/contact`}
              className="text-xs font-bold uppercase tracking-widest text-foreground/70 hover:text-brand-gold transition-colors">
              {dict.contact}
            </Link>
          </nav>

          <div className="flex items-center gap-3 sm:gap-6 z-10">
            {/* Desktop: right-aligned attempts badge */}
            {isAuthenticated && (
              <AttemptsBadge
                initialAttempts={initialAttempts}
                label={dict.attempts}
                containerClassName="hidden lg:flex gap-2"
              />
            )}

            <LanguageSwitcher currentLocale={locale} />

            {isAuthenticated ? (
              <div className="hidden lg:flex items-center gap-4">
                <LogoutButton label={dict.logout} />
              </div>
            ) : (
              <div className="hidden lg:flex items-center gap-4">
                <Link
                  href={`/${locale}/login`}
                  className="text-sm font-bold text-foreground/70 hover:text-foreground transition-colors">
                  {dict.login}
                </Link>
                <Link
                  href={`/${locale}/register`}
                  className="text-sm font-bold bg-brand-gold text-background px-5 py-2 rounded-xl hover:scale-105 transition-transform shadow-[0_0_15px_var(--color-brand-gold-glow)]">
                  {dict.register}
                </Link>
              </div>
            )}

            <MobileMenu
              dict={dict}
              locale={locale}
              isAuthenticated={isAuthenticated}
            />
          </div>
        </div>
      </div>
    </header>
  );
}

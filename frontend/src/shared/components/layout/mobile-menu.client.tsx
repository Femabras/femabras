//femabras/frontend/src/shared/components/layout/mobile-menu.client.tsx
"use client";

import { useState, useEffect } from "react";
import { createPortal } from "react-dom";
import Link from "next/link";
import { authClientService } from "@/modules/auth/services/auth.client.service";
import type { Dictionary } from "@/i18n/get-dictionary";

interface MobileMenuProps {
  dict: Dictionary["layout"];
  locale: string;
  isAuthenticated: boolean;
}

export function MobileMenu({ dict, locale, isAuthenticated }: MobileMenuProps) {
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "unset";
    }
    return () => {
      document.body.style.overflow = "unset";
    };
  }, [isOpen]);

  const closeMenu = () => setIsOpen(false);

  return (
    <div className="lg:hidden flex items-center">
      <button
        onClick={() => setIsOpen(true)}
        className="p-2 text-foreground/70 hover:text-brand-gold transition-colors"
        aria-label="Open menu">
        <svg
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round">
          <line x1="3" y1="12" x2="21" y2="12"></line>
          <line x1="3" y1="6" x2="21" y2="6"></line>
          <line x1="3" y1="18" x2="21" y2="18"></line>
        </svg>
      </button>

      {isOpen &&
        createPortal(
          // 🟢 FIX: Restored z-[100] with brackets. Tailwind v4 handles custom z-index exactly like this.
          <div className="fixed inset-0 z-100 bg-background/80 backdrop-blur-3xl flex flex-col p-6 animate-in slide-in-from-top-8 fade-in duration-300">
            <div className="flex justify-end items-center shrink-0">
              <button
                onClick={closeMenu}
                className="p-2 text-foreground/70 hover:text-brand-gold transition-colors"
                aria-label="Close menu">
                <svg
                  width="28"
                  height="28"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>

            <div className="flex flex-col flex-1 overflow-y-auto justify-start items-center w-full pt-12 sm:pt-16 pb-12">
              <nav className="flex flex-col gap-6 sm:gap-8 text-center items-center w-full">
                <Link
                  href={`/${locale}`}
                  onClick={closeMenu}
                  className="text-xl sm:text-2xl font-black uppercase tracking-widest text-foreground hover:text-brand-gold transition-colors">
                  {dict.home}
                </Link>
                <Link
                  href={`/${locale}/products`}
                  onClick={closeMenu}
                  className="text-xl sm:text-2xl font-black uppercase tracking-widest text-foreground hover:text-brand-gold transition-colors">
                  {dict.products}
                </Link>
                <Link
                  href={`/${locale}/services`}
                  onClick={closeMenu}
                  className="text-xl sm:text-2xl font-black uppercase tracking-widest text-foreground hover:text-brand-gold transition-colors">
                  {dict.services}
                </Link>
                <Link
                  href={`/${locale}/about`}
                  onClick={closeMenu}
                  className="text-xl sm:text-2xl font-black uppercase tracking-widest text-foreground hover:text-brand-gold transition-colors">
                  {dict.about}
                </Link>
                <Link
                  href={`/${locale}/contact`}
                  onClick={closeMenu}
                  className="text-xl sm:text-2xl font-black uppercase tracking-widest text-foreground hover:text-brand-gold transition-colors">
                  {dict.contact}
                </Link>
              </nav>

              <div className="mt-10 sm:mt-12 flex flex-col gap-4 pt-8 border-t border-white/10 w-full max-w-xs mx-auto text-center">
                {isAuthenticated ? (
                  <button
                    onClick={async () => {
                      await authClientService.logout();
                      window.location.href = "/";
                    }}
                    className="text-base font-bold text-red-400 hover:text-red-300 transition-colors py-2">
                    {dict.logout}
                  </button>
                ) : (
                  <>
                    <Link
                      href={`/${locale}/login`}
                      onClick={closeMenu}
                      className="text-base font-bold text-foreground/70 hover:text-foreground transition-colors py-2">
                      {dict.login}
                    </Link>
                    <Link
                      href={`/${locale}/register`}
                      onClick={closeMenu}
                      className="text-base font-bold bg-brand-gold text-background px-6 py-4 rounded-2xl shadow-[0_0_15px_var(--color-brand-gold-glow)]">
                      {dict.register}
                    </Link>
                  </>
                )}
              </div>
            </div>
          </div>,
          document.body,
        )}
    </div>
  );
}

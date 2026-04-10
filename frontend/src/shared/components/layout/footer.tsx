// femabras/frontend/src/shared/components/layout/footer.tsx
import Link from "next/link";
import Image from "next/image";
import type { Dictionary } from "@/i18n/get-dictionary";

interface FooterProps {
  dict: Dictionary["layout"];
  locale: string;
}

export function Footer({ dict, locale }: FooterProps) {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="w-full border-t border-white/10 bg-background/50 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row items-center justify-between h-auto md:h-20 py-6 md:py-0 relative gap-4 md:gap-0">
          {/* Left: Logo */}
          <div className="flex items-center z-10">
            <Link href={`/${locale}`}>
              <Image
                src="/logos/femabras-logo-light.svg"
                alt="Femabras"
                width={120}
                height={28}
                className="w-auto h-5 sm:h-6 opacity-60 grayscale hover:grayscale-0 hover:opacity-100 transition-all duration-300"
              />
            </Link>
          </div>

          {/* Center: Copyright (Absolutely centered on desktop) */}
          <p className="md:absolute md:left-1/2 md:top-1/2 md:-translate-x-1/2 md:-translate-y-1/2 text-xs text-foreground/50 text-center z-0 whitespace-nowrap">
            © {currentYear} Femabras. {dict.rights}
          </p>

          {/* Right: Links */}
          <div className="flex items-center gap-4 sm:gap-6 z-10">
            <Link
              href={`/${locale}/terms`}
              className="text-xs font-medium text-foreground/50 hover:text-brand-gold transition-colors">
              {dict.terms}
            </Link>
            <Link
              href={`/${locale}/privacy`}
              className="text-xs font-medium text-foreground/50 hover:text-brand-gold transition-colors">
              {dict.privacy}
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}

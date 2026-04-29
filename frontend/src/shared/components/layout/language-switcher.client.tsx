// femabras/frontend/src/shared/components/layout/language-switcher.client.tsx
"use client";

import { usePathname, useRouter } from "next/navigation";

export function LanguageSwitcher({ currentLocale }: { currentLocale: string }) {
  const router = useRouter();
  const pathname = usePathname();

  const handleLocaleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newLocale = e.target.value;

    // Split on "/" and replace only the locale segment (index 1).
    // e.g. "/en/about" → ["", "en", "about"] → ["", "fr", "about"] → "/fr/about"
    // The old string-replace approach could corrupt URLs where the locale
    // string appeared elsewhere (e.g. query params or slugs like "/en-products").
    const segments = pathname.split("/");
    segments[1] = newLocale;
    const newPath = segments.join("/");

    router.push(newPath);
    router.refresh();
  };

  return (
    <div className="relative">
      <select
        value={currentLocale}
        onChange={handleLocaleChange}
        className="appearance-none bg-white/5 border border-white/10 text-foreground text-xs sm:text-sm font-bold uppercase tracking-widest rounded-xl px-3 py-2 pr-8 focus:outline-none focus:ring-2 focus:ring-yellow-500/50 cursor-pointer transition-all hover:bg-white/10">
        <option value="en" className="bg-background text-foreground">
          EN
        </option>
        <option value="pt" className="bg-background text-foreground">
          PT
        </option>
        <option value="fr" className="bg-background text-foreground">
          FR
        </option>
      </select>
      <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-foreground/50">
        <svg
          className="fill-current h-4 w-4"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20">
          <path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z" />
        </svg>
      </div>
    </div>
  );
}

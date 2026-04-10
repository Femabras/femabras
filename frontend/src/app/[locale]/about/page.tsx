//femabras/frontend/src/app/[locale]/about/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";

export default async function AboutPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const dict = await getDictionary(locale);

  return (
    <div className="flex flex-1 flex-col items-center justify-center w-full max-w-4xl mx-auto px-4 py-12 gap-10 sm:gap-16 animate-in fade-in slide-in-from-bottom-8 duration-700">
      <div className="text-center">
        <h1 className="text-4xl sm:text-6xl font-black text-foreground tracking-tighter mb-4 drop-shadow-[0_0_15px_var(--color-brand-gold-glow)]">
          {dict.pages.aboutTitle}
        </h1>
        <p className="text-brand-gold font-bold uppercase tracking-widest text-sm sm:text-base">
          {dict.pages.aboutSubtitle}
        </p>
      </div>

      <div className="w-full bg-white/5 border border-white/10 rounded-3xl p-8 sm:p-12 backdrop-blur-xl shadow-2xl relative overflow-hidden">
        <div className="absolute -top-32 -right-32 w-64 h-64 bg-brand-gold/10 rounded-full blur-[80px]" />

        <div className="relative z-10 space-y-6 text-foreground/80 leading-relaxed text-base sm:text-lg">
          <p>{dict.pages.aboutBio1}</p>
          <p>{dict.pages.aboutBio2}</p>
        </div>
      </div>
    </div>
  );
}

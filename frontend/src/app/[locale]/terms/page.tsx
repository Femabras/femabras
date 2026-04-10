//femabras/frontend/src/app/[locale]/terms/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";

export default async function TermsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const dict = await getDictionary(locale);

  return (
    <div className="flex flex-1 flex-col justify-start w-full max-w-3xl mx-auto px-4 py-12 gap-10 animate-in fade-in duration-500">
      <div className="border-b border-white/10 pb-8">
        <h1 className="text-3xl sm:text-5xl font-black text-brand-gold tracking-tighter mb-4">
          {dict.layout.terms}
        </h1>
        <p className="text-sm font-bold uppercase tracking-widest text-foreground/40">
          {dict.pages.legalUpdated}
        </p>
      </div>

      <div className="space-y-8 text-foreground/70 leading-relaxed">
        <section>
          <h2 className="text-2xl font-bold text-foreground mb-4">
            1. Introduction
          </h2>
          <p>
            By accessing Femabras, you agree to be bound by these terms. This
            platform operates the Daily Challenge and processes digital asset
            sales...
          </p>
        </section>

        <section>
          <h2 className="text-2xl font-bold text-foreground mb-4">
            2. Digital Products
          </h2>
          <p>
            All beat licenses and e-books are delivered digitally. Due to the
            nature of digital goods, all sales are final once the download link
            is generated...
          </p>
        </section>
      </div>
    </div>
  );
}

//femabras/frontend/src/app/[locale]/privacy/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";

export default async function PrivacyPage({
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
          {dict.layout.privacy}
        </h1>
        <p className="text-sm font-bold uppercase tracking-widest text-foreground/40">
          {dict.pages.legalUpdated}
        </p>
      </div>

      <div className="space-y-8 text-foreground/70 leading-relaxed">
        <section>
          <h2 className="text-2xl font-bold text-foreground mb-4">
            1. Data Collection
          </h2>
          <p>
            We collect necessary information such as your email address and
            payment details strictly for account security and processing prize
            payouts...
          </p>
        </section>

        <section>
          <h2 className="text-2xl font-bold text-foreground mb-4">
            2. Data Security
          </h2>
          <p>
            Your data is encrypted using modern standards. We do not sell or
            share your personal information with third-party advertisers...
          </p>
        </section>
      </div>
    </div>
  );
}

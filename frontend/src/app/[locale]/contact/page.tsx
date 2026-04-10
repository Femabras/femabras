//femabras/frontend/src/app/[locale]/contact/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { Button } from "@/shared/components/ui/button";

export default async function ContactPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const dict = await getDictionary(locale);

  return (
    <div className="flex flex-1 flex-col items-center justify-center w-full max-w-2xl mx-auto px-4 py-12 gap-10 sm:gap-12 animate-in fade-in zoom-in-95 duration-700">
      <div className="text-center">
        <h1 className="text-4xl sm:text-5xl font-black text-foreground tracking-tighter mb-4">
          {dict.pages.contactTitle}
        </h1>
        <p className="text-foreground/60 text-sm sm:text-base">
          {dict.pages.contactSubtitle}
        </p>
      </div>

      <div className="w-full bg-white/5 border border-white/10 rounded-3xl p-6 sm:p-10 backdrop-blur-md">
        <form className="flex flex-col gap-6">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 w-full">
            <div className="flex flex-col gap-2">
              <Label htmlFor="name">{dict.pages.contactName}</Label>
              <Input id="name" required />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">{dict.pages.contactEmail}</Label>
              <Input id="email" type="email" required />
            </div>
          </div>

          <div className="flex flex-col gap-2">
            <Label htmlFor="message">{dict.pages.contactMessage}</Label>
            <textarea
              id="message"
              rows={5}
              required
              className="flex w-full rounded-2xl border-2 border-white/10 bg-white/5 px-4 py-3 text-base text-foreground placeholder:text-foreground/30 focus-visible:outline-none focus-visible:border-brand-gold focus-visible:ring-4 focus-visible:ring-brand-gold/20 transition-all duration-300 resize-none"></textarea>
          </div>

          <Button type="button" variant="warning" className="w-full mt-2">
            {dict.pages.btnSend}
          </Button>
        </form>
      </div>
    </div>
  );
}

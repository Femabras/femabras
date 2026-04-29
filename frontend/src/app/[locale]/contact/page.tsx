// femabras/frontend/src/app/[locale]/contact/page.tsx
"use client";

import { useState } from "react";
import { useParams } from "next/navigation";
import { Input } from "@/shared/components/ui/input";
import { Label } from "@/shared/components/ui/label";
import { Button } from "@/shared/components/ui/button";
import { env } from "@/shared/config/env";

// Inline type — the contact page is a client component and can't import the
// server-only dictionary loader. Strings are duplicated from pages.* keys
// to keep this self-contained.
type ContactDict = {
  contactTitle: string;
  contactSubtitle: string;
  contactName: string;
  contactEmail: string;
  contactMessage: string;
  btnSend: string;
  contactSending: string;
  contactSuccessTitle: string;
  contactSuccessMsg: string;
};

const contactDicts: Record<string, ContactDict> = {
  en: {
    contactTitle: "Get in Touch",
    contactSubtitle: "Let's build something great.",
    contactName: "Your Name",
    contactEmail: "Your Email",
    contactMessage: "Message",
    btnSend: "Send Message",
    contactSending: "Sending...",
    contactSuccessTitle: "Message sent!",
    contactSuccessMsg: "We'll get back to you soon.",
  },
  pt: {
    contactTitle: "Fale Connosco",
    contactSubtitle: "Vamos construir algo incrível.",
    contactName: "O seu Nome",
    contactEmail: "O seu Email",
    contactMessage: "Mensagem",
    btnSend: "Enviar Mensagem",
    contactSending: "A enviar...",
    contactSuccessTitle: "Mensagem enviada!",
    contactSuccessMsg: "Entraremos em contacto em breve.",
  },
  fr: {
    contactTitle: "Contactez-nous",
    contactSubtitle: "Construisons quelque chose de grand.",
    contactName: "Votre Nom",
    contactEmail: "Votre e-mail",
    contactMessage: "Message",
    btnSend: "Envoyer le Message",
    contactSending: "Envoi en cours...",
    contactSuccessTitle: "Message envoyé !",
    contactSuccessMsg: "Nous vous répondrons bientôt.",
  },
};

export default function ContactPage() {
  const params = useParams();
  const locale = (params?.locale as string) || "en";
  const dict = contactDicts[locale] ?? contactDicts.en;

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrorMsg(null);

    const formData = new FormData(e.currentTarget);
    const name = formData.get("name") as string;
    const email = formData.get("email") as string;
    const message = formData.get("message") as string;

    try {
      const res = await fetch(`${env.apiUrl}/contact`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, email, message }),
      });

      if (!res.ok) {
        const err = await res.json().catch(() => null);
        throw new Error(err?.error || "Failed to send message.");
      }

      setSuccess(true);
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : "An error occurred.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-1 flex-col items-center justify-center w-full max-w-2xl mx-auto px-4 py-12 gap-10 sm:gap-12 animate-in fade-in zoom-in-95 duration-700">
      <div className="text-center">
        <h1 className="text-4xl sm:text-5xl font-black text-foreground tracking-tighter mb-4">
          {dict.contactTitle}
        </h1>
        <p className="text-foreground/60 text-sm sm:text-base">
          {dict.contactSubtitle}
        </p>
      </div>

      <div className="w-full bg-white/5 border border-white/10 rounded-3xl p-6 sm:p-10 backdrop-blur-md">
        {success ? (
          <div className="flex flex-col items-center gap-4 py-8 text-center animate-in zoom-in duration-500">
            <div className="text-5xl">✉️</div>
            <p className="text-lg font-bold text-brand-gold">
              {dict.contactSuccessTitle}
            </p>
            <p className="text-sm text-foreground/60">
              {dict.contactSuccessMsg}
            </p>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="flex flex-col gap-6">
            {errorMsg && (
              <div className="p-3 bg-red-500/20 border border-red-500/50 rounded-xl text-red-200 text-sm text-center">
                {errorMsg}
              </div>
            )}

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 w-full">
              <div className="flex flex-col gap-2">
                <Label htmlFor="name">{dict.contactName}</Label>
                <Input id="name" name="name" required disabled={isSubmitting} />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="email">{dict.contactEmail}</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  required
                  disabled={isSubmitting}
                />
              </div>
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="message">{dict.contactMessage}</Label>
              <textarea
                id="message"
                name="message"
                rows={5}
                required
                disabled={isSubmitting}
                className="flex w-full rounded-2xl border-2 border-white/10 bg-white/5 px-4 py-3 text-base text-foreground placeholder:text-foreground/30 focus-visible:outline-none focus-visible:border-brand-gold focus-visible:ring-4 focus-visible:ring-brand-gold/20 transition-all duration-300 resize-none disabled:opacity-50"
              />
            </div>

            <Button
              type="submit"
              variant="warning"
              className="w-full mt-2"
              disabled={isSubmitting}>
              {isSubmitting ? dict.contactSending : dict.btnSend}
            </Button>
          </form>
        )}
      </div>
    </div>
  );
}

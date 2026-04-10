//femabras/frontend/src/app/[locale]/layout.tsx
import type { Metadata } from "next";
import { SpeedInsights } from "@vercel/speed-insights/next";
import { siteConfig } from "@/shared/config/site";
import { chillax } from "@/shared/assets/fonts";
import { Header } from "@/shared/components/layout/header";
import { Footer } from "@/shared/components/layout/footer";
import { getDictionary } from "@/i18n/get-dictionary";
import "../globals.css";

export const metadata: Metadata = {
  title: siteConfig.title,
  description: siteConfig.description,
};

type LayoutProps = {
  children: React.ReactNode;
  params: Promise<{ locale: string }>;
};

export default async function RootLayout({ children, params }: LayoutProps) {
  const { locale } = await params;
  const isProduction = process.env.NODE_ENV === "production";

  // Load the dictionary on the server to pass to the layout components
  const dict = await getDictionary(locale);

  return (
    <html lang={locale}>
      <body
        className={`${chillax.variable} bg-background text-foreground antialiased min-h-screen flex flex-col`}>
        <Header dict={dict.layout} locale={locale} />
        <main className="flex flex-col grow">{children}</main>
        <Footer dict={dict.layout} locale={locale} />
        {isProduction && <SpeedInsights />}
      </body>
    </html>
  );
}

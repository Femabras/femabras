//femabras/frontend/src/app/[locale]/services/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";

export default async function ServicesPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const dict = await getDictionary(locale);

  const services = [
    {
      title: dict.pages.svcSoftware,
      desc: dict.pages.svcSoftwareDesc,
      icon: "💻",
    },
    { title: dict.pages.svcMusic, desc: dict.pages.svcMusicDesc, icon: "🎹" },
  ];

  return (
    <div className="flex flex-1 flex-col items-center justify-center w-full max-w-6xl mx-auto px-4 py-12 gap-10 sm:gap-16 animate-in fade-in duration-700">
      <div className="text-center">
        <h1 className="text-4xl sm:text-6xl font-black text-foreground tracking-tighter mb-4">
          {dict.pages.servicesTitle}
        </h1>
        <p className="text-brand-gold font-bold uppercase tracking-widest text-sm">
          {dict.pages.servicesSubtitle}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8 w-full">
        {services.map((svc, idx) => (
          <div
            key={idx}
            className="group bg-white/5 border border-white/10 rounded-3xl p-8 hover:bg-white/10 transition-colors backdrop-blur-md flex flex-col gap-4">
            <div className="text-4xl sm:text-5xl group-hover:scale-110 transition-transform origin-left">
              {svc.icon}
            </div>
            <h3 className="text-xl sm:text-2xl font-black text-brand-gold">
              {svc.title}
            </h3>
            <p className="text-sm sm:text-base text-foreground/70 leading-relaxed">
              {svc.desc}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}

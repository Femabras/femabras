//femabras/frontend/src/app/[locale]/products/page.tsx
import { getDictionary } from "@/i18n/get-dictionary";

export default async function ProductsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const dict = await getDictionary(locale);

  const products = [
    {
      title: "Afropop Anthem Vol. 1",
      type: dict.pages.prodBeat,
      desc: dict.pages.prodBeatDesc,
      price: "15,000 AOA",
    },
    {
      title: "Kizomba Nightfall",
      type: dict.pages.prodBeat,
      desc: dict.pages.prodBeatDesc,
      price: "20,000 AOA",
    },
    {
      title: "Go Zero-Dependency",
      type: dict.pages.prodBook,
      desc: dict.pages.prodBookDesc,
      price: "5,000 AOA",
    },
  ];

  return (
    <div className="flex flex-1 flex-col items-center justify-center w-full max-w-6xl mx-auto px-4 py-12 gap-10 sm:gap-16 animate-in fade-in duration-700">
      <div className="text-center">
        <h1 className="text-4xl sm:text-6xl font-black text-foreground tracking-tighter mb-4">
          {dict.pages.productsTitle}
        </h1>
        <p className="text-brand-gold font-bold uppercase tracking-widest text-sm">
          {dict.pages.productsSubtitle}
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
        {products.map((prod, idx) => (
          <div
            key={idx}
            className="flex flex-col bg-white/5 border border-white/10 rounded-3xl p-6 hover:border-brand-gold/30 transition-all gap-4">
            <div className="flex flex-col gap-2 flex-grow">
              <span className="text-[10px] font-bold uppercase tracking-widest text-brand-gold px-3 py-1 bg-brand-gold/10 rounded-full w-fit">
                {prod.type}
              </span>
              <h3 className="text-xl font-black text-foreground">
                {prod.title}
              </h3>
              <p className="text-sm text-foreground/60">{prod.desc}</p>
            </div>
            <div className="flex items-center justify-between pt-6 border-t border-white/10 mt-auto">
              <span className="font-bold text-foreground text-sm sm:text-base">
                {prod.price}
              </span>
              <button className="text-[10px] sm:text-xs font-black uppercase tracking-widest text-background bg-brand-gold px-4 py-2 rounded-lg hover:scale-105 transition-transform">
                {dict.pages.btnBuy}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

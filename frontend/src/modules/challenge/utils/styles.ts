//frontend/src/modules/challenge/utils/styles.ts
export const UI = {
  authOverlay:
    "absolute inset-0 z-50 flex flex-col items-center justify-center bg-background/95 backdrop-blur-md rounded-2xl p-6",

  guestBanner:
    "bg-white/5 border border-white/10 px-5 py-2 rounded-xl shadow-lg backdrop-blur-sm",

  toast:
    "bg-red-500/90 text-white px-6 py-3 rounded-full shadow-lg font-bold text-sm flex items-center gap-2",

  title: "font-bold text-2xl sm:text-4xl tracking-tight mb-2",
} as const;

//femabras/frontend/src/shared/config/gameStyles.ts
import { cn } from "@/shared/lib/utils";

export const GAME = {
  slot: (isShaking: boolean, isOver: boolean, hasValue: boolean) =>
    cn(
      "relative flex w-full max-w-20 aspect-4/5 items-center justify-center rounded-lg border-2 transition-all duration-300",
      isShaking &&
        "border-red-500 bg-red-500/15 border-solid shadow-lg shadow-red-500/20",
      !isShaking && hasValue && "border-foreground bg-foreground border-solid",
      !isShaking &&
        !hasValue &&
        isOver &&
        "border-foreground/50 bg-foreground/10 border-dashed",
      !isShaking &&
        !hasValue &&
        !isOver &&
        "border-foreground/20 bg-transparent border-dashed",
    ),
  number: (isDragging: boolean) =>
    cn(
      "flex h-16 w-16 select-none cursor-grab items-center justify-center rounded-full bg-foreground text-background text-2xl font-black active:cursor-grabbing active:scale-110",
      isDragging
        ? "opacity-0 scale-50"
        : "opacity-100 scale-100 shadow-xl transition-all duration-300 ease-out",
    ),
  submitBtn:
    "w-full sm:w-auto px-10 sm:px-14 py-4 sm:py-5 bg-foreground text-background font-black rounded-2xl shadow-xl transition-all uppercase tracking-widest active:scale-95 disabled:opacity-50",
  adjustBtn:
    "w-full sm:w-auto px-10 sm:px-14 py-4 sm:py-5 border-2 border-foreground/20 text-foreground/80 font-bold rounded-2xl transition-all uppercase tracking-widest hover:bg-white/5",
  animations: { popIn: "animate-pop-in", shake: "shake-wrapper" },
};

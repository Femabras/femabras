//femabras/frontend/src/shared/components/ui/input.tsx
import { forwardRef, InputHTMLAttributes } from "react";
import { cn } from "@/shared/lib/utils";

export const Input = forwardRef<
  HTMLInputElement,
  InputHTMLAttributes<HTMLInputElement>
>(({ className, type, ...props }, ref) => {
  return (
    <input
      type={type}
      className={cn(
        "flex h-14 w-full rounded-2xl border-2 border-white/10 bg-white/5 px-4 py-2 text-base text-foreground ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-foreground/30 focus-visible:outline-none focus-visible:border-yellow-500 focus-visible:ring-4 focus-visible:ring-yellow-500/20 disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-300",
        className,
      )}
      ref={ref}
      {...props}
    />
  );
});
Input.displayName = "Input";

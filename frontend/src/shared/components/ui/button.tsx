//frontend/src/shared/components/ui/button.tsx
import { ButtonHTMLAttributes, forwardRef } from "react";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "warning";
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className = "", variant = "primary", children, ...props }, ref) => {
    // Base classes applied to ALL buttons
    const baseClass =
      "w-full sm:w-auto px-10 sm:px-14 py-4 sm:py-5 font-black rounded-2xl transition-all uppercase tracking-widest active:scale-95 disabled:opacity-50";

    // Variant-specific classes
    const variants = {
      primary: "bg-foreground text-background shadow-xl",
      secondary:
        "border-2 border-foreground/20 text-foreground/80 hover:bg-white/5",
      warning:
        "bg-yellow-500 text-black shadow-xl shadow-yellow-500/40 hover:scale-105",
    };

    return (
      <button
        ref={ref}
        className={`${baseClass} ${variants[variant]} ${className}`}
        {...props}>
        {children}
      </button>
    );
  },
);
Button.displayName = "Button";

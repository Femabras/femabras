//frontend/src/shared/components/ui/label.tsx
import { LabelHTMLAttributes, forwardRef } from "react";
import { cn } from "@/shared/lib/utils";

export const Label = forwardRef<
  HTMLLabelElement,
  LabelHTMLAttributes<HTMLLabelElement>
>(({ className, ...props }, ref) => {
  return (
    <label
      ref={ref}
      className={cn(
        "text-sm font-bold leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 text-foreground/90 ml-1",
        className,
      )}
      {...props}
    />
  );
});
Label.displayName = "Label";

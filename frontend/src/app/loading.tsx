//frontend/src/app/loading.tsx
export default function Loading() {
  return (
    <div className="flex min-h-screen items-center justify-center text-foreground">
      {/* A universally understood, pulsing loading indicator */}
      <div className="flex flex-col items-center gap-4 animate-pulse">
        <div className="h-12 w-12 rounded-full border-4 border-yellow-500 border-t-transparent animate-spin" />
        <p className="text-sm font-bold uppercase tracking-widest text-foreground/70">
          Loading / Carregando
        </p>
      </div>
    </div>
  );
}

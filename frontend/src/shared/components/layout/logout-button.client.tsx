//femabras/frontend/src/shared/components/layout/logout-button.client.tsx
"use client";

import { authClientService } from "@/modules/auth/services/auth.client.service";

export function LogoutButton({ label }: { label: string }) {
  const handleLogout = async () => {
    try {
      await authClientService.logout();
      sessionStorage.removeItem("femabras_saved_guess");
      localStorage.removeItem("femabras_attempts");
      window.location.href = "/";
    } catch (error) {
      console.error("Logout failed:", error);
      window.location.href = "/";
    }
  };

  return (
    <button
      onClick={handleLogout}
      className="text-sm font-bold text-foreground/70 hover:text-red-400 transition-colors">
      {label}
    </button>
  );
}

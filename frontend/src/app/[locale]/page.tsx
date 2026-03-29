//femabras/frontend/src/app/[locale]/page.tsx
import { cookies } from "next/headers";
import { challengeServerService } from "@/modules/challenge/services/challenge.server.service";
import { GameBoard } from "@/modules/challenge/components/game-board.client";
import { WinnerPodium } from "@/modules/challenge/components/winner-podium.client";
import { getDictionary } from "@/i18n/get-dictionary";

export const dynamic = "force-dynamic";

export default async function Home({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const challenge = await challengeServerService.getDailyChallenge();
  const dict = await getDictionary(locale);

  const cookieStore = await cookies();
  const isAuthenticated = cookieStore.has("auth_token");

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4">
      {challenge.status === "solved" ? (
        <WinnerPodium
          challenge={challenge}
          dict={dict.challenge}
          isAuthenticated={isAuthenticated}
        />
      ) : (
        <GameBoard
          challenge={challenge}
          isAuthenticated={isAuthenticated}
          dict={dict.challenge}
        />
      )}
    </div>
  );
}

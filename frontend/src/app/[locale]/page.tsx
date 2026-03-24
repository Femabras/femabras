//frontend/src/app/[locale]/page.tsx
import { cookies } from "next/headers";
import { challengeServerService } from "@/modules/challenge/services/challenge.server.service";
import { GameBoard } from "@/modules/challenge/components/game-board.client";
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
      {/* Pass the dictionary to the client so UI components stay pure */}
      <GameBoard
        challenge={challenge}
        isAuthenticated={isAuthenticated}
        dict={dict.challenge}
      />
    </div>
  );
}

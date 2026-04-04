//femabras/frontend/src/app/[locale]/page.tsx
import { cookies } from "next/headers";
import { challengeServerService } from "@/modules/challenge/services/challenge.server.service";
import { GameBoard } from "@/modules/challenge/components/game-board.client";
import { WinnerPodium } from "@/modules/challenge/components/winner-podium.client";
import { ClaimPrizeForm } from "@/modules/challenge/components/claim-prize-form.client";
import { WinnerStatus } from "@/modules/challenge/components/winner-status.client";
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

  let isWinner = false;
  let payoutStatus: "unclaimed" | "pending" | "paid" | "rejected" = "unclaimed";

  if (challenge.status === "solved" && isAuthenticated) {
    const myStatus = await challengeServerService.getMyStatus();
    isWinner = myStatus.is_winner;
    payoutStatus = myStatus.payout_status;
  }

  let pageContent;

  if (challenge.status === "solved") {
    if (isWinner) {
      if (payoutStatus === "unclaimed") {
        pageContent = (
          <ClaimPrizeForm prizeAmount={challenge.prize} dict={dict.challenge} />
        );
      } else {
        pageContent = (
          <WinnerStatus
            prize={challenge.prize}
            payoutStatus={challenge.payout_status}
            dict={dict.challenge}
          />
        );
      }
    } else {
      pageContent = (
        <WinnerPodium
          challenge={challenge}
          dict={dict.challenge}
          isAuthenticated={isAuthenticated}
        />
      );
    }
  } else {
    pageContent = (
      <GameBoard
        challenge={challenge}
        isAuthenticated={isAuthenticated}
        dict={dict.challenge}
      />
    );
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4">
      {pageContent}
    </div>
  );
}

//femabras/frontend/src/shared/lib/events.ts
export const ATTEMPTS_EVENT = "femabras:update_attempts";

export function dispatchAttemptsUpdate(newAttempts: number) {
  if (typeof window !== "undefined") {
    window.dispatchEvent(
      new CustomEvent<number>(ATTEMPTS_EVENT, { detail: newAttempts }),
    );
  }
}

// femabras/frontend/src/shared/lib/api.client.ts
//
// Centralised fetch wrapper that automatically attaches the CSRF token to
// state-changing requests. The backend's middleware/csrf.go uses the
// double-submit cookie pattern: a non-HttpOnly cookie named "csrf_token"
// is set on every safe request, and we echo it in X-CSRF-Token on POSTs.
//
// All client-side modules should call apiFetch() instead of fetch() directly
// for any request to our backend that mutates state.

import { env } from "@/shared/config/env";

const CSRF_COOKIE = "csrf_token";
const CSRF_HEADER = "X-CSRF-Token";

const MUTATING_METHODS = new Set(["POST", "PUT", "PATCH", "DELETE"]);

function readCsrfCookie(): string | null {
  if (typeof document === "undefined") return null;
  const match = document.cookie.match(
    new RegExp(`(?:^|;\\s*)${CSRF_COOKIE}=([^;]+)`),
  );
  return match ? decodeURIComponent(match[1]) : null;
}

/**
 * Ensure the CSRF cookie has been issued by the backend before a mutating
 * request fires. If we don't have the cookie yet (first visit, fresh tab),
 * call GET /csrf which issues it. Idempotent and cheap.
 */
async function ensureCsrfToken(): Promise<string | null> {
  const token = readCsrfCookie();
  if (token) return token;

  try {
    await fetch(`${env.apiUrl}/csrf`, {
      method: "GET",
      credentials: "include",
    });
  } catch {
    // Network error — return null and let the actual request fail with a
    // clear error rather than silently dropping it.
    return null;
  }

  return readCsrfCookie();
}

/**
 * Wrapper around fetch that adds CSRF headers to state-changing requests.
 * Use this for every authenticated mutation.
 */
export async function apiFetch(
  input: string,
  init: RequestInit = {},
): Promise<Response> {
  const method = (init.method || "GET").toUpperCase();
  const headers = new Headers(init.headers);

  if (MUTATING_METHODS.has(method)) {
    const token = await ensureCsrfToken();
    if (token) {
      headers.set(CSRF_HEADER, token);
    }
  }

  return fetch(input, {
    ...init,
    headers,
    credentials: init.credentials ?? "include",
  });
}

// femabras/frontend/src/modules/auth/actions/auth.actions.ts
"use server";

import { cookies } from "next/headers";

/**
 * Server-side logout action.
 *
 * Clears both auth cookies that the backend sets:
 *   - access_token  (was incorrectly named "auth_token" before)
 *   - refresh_token
 *
 * This is called as a fallback or from Server Components. Client Components
 * use authClientService.logout() which also hits the backend /logout endpoint
 * to revoke the refresh token in the database.
 */
export async function logoutAction() {
  const cookieStore = await cookies();
  cookieStore.delete("access_token");
  cookieStore.delete("refresh_token");
}

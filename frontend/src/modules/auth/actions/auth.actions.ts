//frontend/src/modules/auth/actions/auth.actions.ts
"use server";

import { cookies } from "next/headers";

export async function logoutAction() {
  const cookieStore = await cookies();
  cookieStore.delete("auth_token");
}

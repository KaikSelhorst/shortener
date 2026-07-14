import { redirect } from "@sveltejs/kit";
import { ACCESS_TOKEN_COOKIE } from "$lib/server/auth";
import type { LayoutServerLoad } from "./$types";

// Already-authenticated users have no business on login/register/mfa —
// send them straight to the app instead of letting them re-auth.
export const load: LayoutServerLoad = ({ cookies }) => {
  if (cookies.get(ACCESS_TOKEN_COOKIE)) redirect(303, "/");
};

import { fail, redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { MFA_SESSION_COOKIE, setAuthCookies } from "$lib/server/auth";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ cookies }) => {
  if (!cookies.get(MFA_SESSION_COOKIE)) redirect(303, "/login");
};

export const actions: Actions = {
  default: async ({ request, cookies, fetch }) => {
    const session = cookies.get(MFA_SESSION_COOKIE);
    if (!session) return fail(401, { error: "Session expired, please sign in again." });

    const data = await request.formData();
    const code = data.get("code");

    if (typeof code !== "string" || !code) {
      return fail(400, { error: "Code is required." });
    }

    const res = await fetch(`${env.API_URL}/auth/mfa/totp`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ session, code }),
    });
    const body = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: body.error ?? "Invalid code." });
    }

    cookies.delete(MFA_SESSION_COOKIE, { path: "/" });
    setAuthCookies(cookies, body);
    redirect(303, "/p");
  },
};

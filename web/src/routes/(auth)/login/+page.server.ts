import { fail, redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { setAuthCookies, setMfaSessionCookie } from "$lib/server/auth";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ request, cookies, fetch }) => {
    const data = await request.formData();
    const email = data.get("email");
    const password = data.get("password");

    if (typeof email !== "string" || typeof password !== "string" || !email || !password) {
      return fail(400, { error: "Email and password are required." });
    }

    const res = await fetch(`${env.API_URL}/auth/login`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
    const body = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: body.error ?? "Something went wrong." });
    }

    if (body.next === "totp") {
      setMfaSessionCookie(cookies, body.session);
      redirect(303, "/mfa");
    }

    setAuthCookies(cookies, body);
    redirect(303, "/");
  },
};

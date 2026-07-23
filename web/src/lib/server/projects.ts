import type { Cookies } from "@sveltejs/kit";
import { dev } from "$app/environment";

export const LAST_PROJECT_COOKIE = "last_project";

const LAST_PROJECT_MAX_AGE = 60 * 60 * 24 * 365; // 1 year — pure UX convenience, not auth

export function setLastProjectCookie(cookies: Cookies, slug: string) {
  cookies.set(LAST_PROJECT_COOKIE, slug, {
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure: !dev,
    maxAge: LAST_PROJECT_MAX_AGE,
  });
}

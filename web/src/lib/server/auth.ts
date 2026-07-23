import type { Cookies } from "@sveltejs/kit";
import { dev } from "$app/environment";
import { env } from "$env/dynamic/private";

export const MFA_SESSION_COOKIE = "mfa_session";
export const ACCESS_TOKEN_COOKIE = "access_token";
const REFRESH_TOKEN_COOKIE = "refresh_token";

// Mirrors the API's service.RefreshTokenTTL (api/internal/service/auth.go) —
// keep in sync if that value changes.
const REFRESH_TOKEN_MAX_AGE = 60 * 60 * 24 * 7;
// Mirrors the API's 5-minute TOTP session window (api/internal/handler/totp.go).
const MFA_SESSION_MAX_AGE = 60 * 5;

interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export function setAuthCookies(cookies: Cookies, tokens: AuthTokens) {
  const secure = !dev;

  // Both cookies live as long as the refresh token, not just the access token's
  // own 15-minute expiry — the cookie has to outlive the JWT so apiFetch() gets
  // a chance to see a 401 and transparently refresh, instead of the browser
  // deleting the cookie out from under it and forcing a full re-login.
  cookies.set(ACCESS_TOKEN_COOKIE, tokens.access_token, {
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure,
    maxAge: REFRESH_TOKEN_MAX_AGE,
  });

  cookies.set(REFRESH_TOKEN_COOKIE, tokens.refresh_token, {
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure,
    maxAge: REFRESH_TOKEN_MAX_AGE,
  });
}

export function setMfaSessionCookie(cookies: Cookies, session: string) {
  cookies.set(MFA_SESSION_COOKIE, session, {
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure: !dev,
    maxAge: MFA_SESSION_MAX_AGE,
  });
}

export function getRefreshToken(cookies: Cookies) {
  return cookies.get(REFRESH_TOKEN_COOKIE);
}

export function clearAuthCookies(cookies: Cookies) {
  cookies.delete(ACCESS_TOKEN_COOKIE, { path: "/" });
  cookies.delete(REFRESH_TOKEN_COOKIE, { path: "/" });
}

// Keyed by the refresh token value itself (not per-request) so concurrent
// requests trying to refresh the SAME token — two tabs, or SvelteKit's own
// hover-preload firing a background request right before a click — share one
// in-flight call instead of racing to spend a single-use token. Only the
// winner's `cookies` gets mutated; a loser's response just omits Set-Cookie
// for these two cookies, leaving the browser's jar untouched by it, so the
// winner's write (independent of response ordering) is what sticks.
// In-memory and single-process — fine for this app's current adapter-node,
// single-instance deployment. Won't help across multiple server instances.
const inflightRefreshes = new Map<string, Promise<string | null>>();

// Exchanges the refresh token cookie for a new token pair via POST /auth/refresh,
// rotating both cookies on success. Returns null (and clears both cookies) if the
// refresh token is missing, expired, or already revoked — the caller should treat
// that as "session is over."
export async function refreshAccessToken(
  cookies: Cookies,
  fetchFn: typeof fetch,
): Promise<string | null> {
  const refreshToken = cookies.get(REFRESH_TOKEN_COOKIE);
  if (!refreshToken) return null;

  const inflight = inflightRefreshes.get(refreshToken);
  if (inflight) return inflight;

  const promise = (async () => {
    try {
      const res = await fetchFn(`${env.API_URL}/auth/refresh`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (!res.ok) {
        clearAuthCookies(cookies);
        return null;
      }

      const tokens: AuthTokens = await res.json();
      setAuthCookies(cookies, tokens);
      return tokens.access_token;
    } finally {
      inflightRefreshes.delete(refreshToken);
    }
  })();

  inflightRefreshes.set(refreshToken, promise);
  return promise;
}

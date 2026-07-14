import type { Cookies } from "@sveltejs/kit";
import { dev } from "$app/environment";

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

  cookies.set(ACCESS_TOKEN_COOKIE, tokens.access_token, {
    path: "/",
    httpOnly: true,
    sameSite: "lax",
    secure,
    maxAge: tokens.expires_in,
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

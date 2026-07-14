import { redirect, type Cookies } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { ACCESS_TOKEN_COOKIE, refreshAccessToken } from "./auth";

interface ApiFetchInit {
  method?: string;
  headers?: Record<string, string>;
  body?: string;
}

interface ApiFetchEvent {
  cookies: Cookies;
  fetch: typeof fetch;
}

// Authenticated fetch against the Go API. On a 401 (expired access token), it
// transparently exchanges the refresh token for a new pair and retries once
// before giving up — callers never have to think about token expiry, only
// about redirect(303, "/login") never coming back if the refresh also fails.
export async function apiFetch(
  { cookies, fetch }: ApiFetchEvent,
  path: string,
  init: ApiFetchInit = {},
): Promise<Response> {
  const request = (accessToken: string | undefined) =>
    fetch(`${env.API_URL}${path}`, {
      ...init,
      headers: { ...init.headers, authorization: `Bearer ${accessToken}` },
    });

  const res = await request(cookies.get(ACCESS_TOKEN_COOKIE));
  if (res.status !== 401) return res;

  const newAccessToken = await refreshAccessToken(cookies, fetch);
  if (!newAccessToken) redirect(303, "/login");

  return request(newAccessToken);
}

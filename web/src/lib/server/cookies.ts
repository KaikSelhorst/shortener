import { dev } from '$app/environment'
import type { Cookies } from '@sveltejs/kit'

const BASE_OPTS = {
	path: '/',
	httpOnly: true,
	sameSite: 'lax' as const,
	// $app/environment.dev is the idiomatic SvelteKit way to detect the
	// current environment — unlike process.env.NODE_ENV it is guaranteed
	// to be tree-shaken correctly and works in all SvelteKit adapters.
	secure: !dev,
}

// Tokens required to set auth cookies. Matches both TokenResponse and a
// completed AuthState (next === 'complete'), where both fields are defined.
interface AuthTokens {
	access_token: string
	refresh_token: string
}

// assertComplete narrows an AuthState to AuthTokens by asserting that
// access_token and refresh_token are present. Call this only after confirming
// next === 'complete' or when the endpoint always returns a completed state
// (e.g. /auth/register).
export function assertComplete(state: {
	access_token?: string
	refresh_token?: string
}): AuthTokens {
	if (!state.access_token || !state.refresh_token) {
		throw new Error('Expected a completed auth state but tokens are missing')
	}
	return { access_token: state.access_token, refresh_token: state.refresh_token }
}

export function setAuthCookies(cookies: Cookies, tokens: AuthTokens): void {
	cookies.set('access_token', tokens.access_token, { ...BASE_OPTS, maxAge: 60 * 15 })
	cookies.set('refresh_token', tokens.refresh_token, {
		...BASE_OPTS,
		maxAge: 60 * 60 * 24 * 7,
	})
}

export function deleteAuthCookies(cookies: Cookies): void {
	cookies.delete('access_token', { path: '/' })
	cookies.delete('refresh_token', { path: '/' })
}

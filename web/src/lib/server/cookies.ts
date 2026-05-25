import { dev } from '$app/environment'
import type { Cookies } from '@sveltejs/kit'
import type { TokenResponse } from '$lib/types'

const BASE_OPTS = {
	path: '/',
	httpOnly: true,
	sameSite: 'lax' as const,
	// $app/environment.dev is the idiomatic SvelteKit way to detect the
	// current environment — unlike process.env.NODE_ENV it is guaranteed
	// to be tree-shaken correctly and works in all SvelteKit adapters.
	secure: !dev,
}

export function setAuthCookies(
	cookies: Cookies,
	tokens: Pick<TokenResponse, 'access_token' | 'refresh_token'>,
): void {
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

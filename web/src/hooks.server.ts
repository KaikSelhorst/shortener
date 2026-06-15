import type { Handle } from '@sveltejs/kit'
import { API_URL } from '$lib/server/config'
import { setAuthCookies, deleteAuthCookies } from '$lib/server/cookies'

function decodeJwtPayload(token: string): { user_id: number; exp: number } | null {
	try {
		const part = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
		return JSON.parse(atob(part))
	} catch {
		return null
	}
}

export const handle: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('access_token')

	if (token) {
		const payload = decodeJwtPayload(token)
		if (payload && payload.exp * 1000 > Date.now()) {
			event.locals.user = { id: payload.user_id }
			return resolve(event)
		}
	}

	const refreshToken = event.cookies.get('refresh_token')
	if (!refreshToken) return resolve(event)

	try {
		const res = await fetch(`${API_URL}/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: refreshToken }),
		})

		if (res.ok) {
			const data = await res.json()
			setAuthCookies(event.cookies, data)
			const newPayload = decodeJwtPayload(data.access_token)
			if (newPayload) event.locals.user = { id: newPayload.user_id }
		} else {
			deleteAuthCookies(event.cookies)
		}
	} catch {
		// network error — proceed unauthenticated
	}

	return resolve(event)
}

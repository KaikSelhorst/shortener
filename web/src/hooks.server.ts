import type { Handle } from '@sveltejs/kit'

function decodeJwtPayload(token: string): { user_id: number; exp: number } | null {
	try {
		const part = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
		return JSON.parse(atob(part))
	} catch {
		return null
	}
}

const COOKIE_OPTS = {
	path: '/',
	httpOnly: true,
	sameSite: 'lax' as const,
	secure: process.env.NODE_ENV === 'production',
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
		const base = process.env.PUBLIC_API_URL ?? 'http://localhost:8080'
		const res = await fetch(`${base}/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: refreshToken }),
		})

		if (res.ok) {
			const data = await res.json()
			event.cookies.set('access_token', data.access_token, {
				...COOKIE_OPTS,
				maxAge: 60 * 15,
			})
			event.cookies.set('refresh_token', data.refresh_token, {
				...COOKIE_OPTS,
				maxAge: 60 * 60 * 24 * 7,
			})
			const newPayload = decodeJwtPayload(data.access_token)
			if (newPayload) event.locals.user = { id: newPayload.user_id }
		} else {
			event.cookies.delete('access_token', { path: '/' })
			event.cookies.delete('refresh_token', { path: '/' })
		}
	} catch {
		// network error — proceed unauthenticated
	}

	return resolve(event)
}

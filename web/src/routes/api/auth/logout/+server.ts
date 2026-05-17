import { redirect } from '@sveltejs/kit'
import type { RequestHandler } from './$types'
import { createApi } from '$lib/api'

export const POST: RequestHandler = async ({ cookies, fetch }) => {
	const refreshToken = cookies.get('refresh_token')

	if (refreshToken) {
		const token = cookies.get('access_token')
		try {
			await createApi(fetch, token).auth.logout(refreshToken)
		} catch {
			// best-effort: clear cookies regardless
		}
	}

	cookies.delete('access_token', { path: '/' })
	cookies.delete('refresh_token', { path: '/' })

	redirect(303, '/login')
}

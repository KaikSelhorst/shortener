import { redirect, fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'

export const load: PageServerLoad = ({ locals }) => {
	if (locals.user) redirect(302, '/dashboard')
}

export const actions: Actions = {
	default: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const email = data.get('email') as string
		const password = data.get('password') as string

		try {
			const tokens = await createApi(fetch).auth.login(email, password)

			cookies.set('access_token', tokens.access_token, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				secure: process.env.NODE_ENV === 'production',
				maxAge: 60 * 15,
			})
			cookies.set('refresh_token', tokens.refresh_token, {
				path: '/',
				httpOnly: true,
				sameSite: 'lax',
				secure: process.env.NODE_ENV === 'production',
				maxAge: 60 * 60 * 24 * 7,
			})
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Login failed' })
		}

		redirect(302, '/dashboard')
	},
}

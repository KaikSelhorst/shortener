import { redirect, fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'
import { setAuthCookies, assertComplete } from '$lib/server/cookies'

export const load: PageServerLoad = ({ locals }) => {
	if (locals.user) redirect(302, '/dashboard')
}

export const actions: Actions = {
	default: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const email = data.get('email') as string
		const password = data.get('password') as string

		try {
			// Register always returns next === 'complete' — new users never have TOTP enabled.
			const state = await createApi(fetch).auth.register(email, password)
			setAuthCookies(cookies, assertComplete(state))
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Registration failed' })
		}

		redirect(302, '/dashboard')
	},
}

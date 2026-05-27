import { redirect, fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi, ApiError } from '$lib/api'
import { setAuthCookies, assertComplete } from '$lib/server/cookies'

export const load: PageServerLoad = ({ locals }) => {
	if (locals.user) redirect(302, '/dashboard')
}

export const actions: Actions = {
	// Step 1: validate email + password.
	// If TOTP is enabled, returns { session } so the page shows the code input.
	// Otherwise sets cookies and redirects to the dashboard.
	credentials: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const email = data.get('email') as string
		const password = data.get('password') as string

		try {
			const state = await createApi(fetch).auth.login(email, password)

			if (state.next === 'totp') {
				return { session: state.session }
			}

			setAuthCookies(cookies, assertComplete(state))
		} catch (err) {
			return fail(400, { error: err instanceof ApiError ? err.message : 'Login failed' })
		}

		redirect(302, '/dashboard')
	},

	// Step 2: validate the TOTP code using the session token from step 1.
	totp: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const session = data.get('session') as string
		const code = data.get('code') as string

		if (!session || !code) {
			return fail(400, { error: 'Session and code are required' })
		}

		try {
			const state = await createApi(fetch).auth.totp.validateMFA(session, code)
			setAuthCookies(cookies, assertComplete(state))
		} catch (err) {
			return fail(400, {
				session,
				error: err instanceof ApiError ? err.message : 'Invalid code',
			})
		}

		redirect(302, '/dashboard')
	},
}

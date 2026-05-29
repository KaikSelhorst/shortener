import { fail, error } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi, ApiError } from '$lib/api'

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const token = cookies.get('access_token')
	const api = createApi(fetch, token)
	try {
		const [apiKeys, projects, me] = await Promise.all([
			api.apiKeys.list(),
			api.projects.list(),
			api.auth.me(),
		])
		return { apiKeys, projects, me }
	} catch (err) {
		if (err instanceof ApiError) error(err.status, err.message)
		error(500, 'Failed to load settings')
	}
}

export const actions: Actions = {
	create: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const name = data.get('name') as string
		const scopes = data.getAll('scopes') as string[]
		const project_id_raw = data.get('project_id') as string
		const project_id = project_id_raw ? Number(project_id_raw) : undefined

		if (!name?.trim()) return fail(400, { createError: 'Name is required' })
		if (!scopes.length) return fail(400, { createError: 'Select at least one scope' })

		try {
			const token = cookies.get('access_token')
			const key = await createApi(fetch, token).apiKeys.create({
				name: name.trim(),
				scopes,
				project_id,
			})
			return { created: key }
		} catch (err) {
			return fail(400, { createError: err instanceof Error ? err.message : 'Failed to create key' })
		}
	},

	delete: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const id = Number(data.get('id'))

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).apiKeys.delete(id)
		} catch (err) {
			return fail(400, { deleteError: err instanceof Error ? err.message : 'Failed to delete key' })
		}
	},

	// Initiates TOTP setup: generates a secret and returns the otpauth:// URI.
	// Returns 409 (propagated from the backend) if TOTP is already enabled.
	totpSetup: async ({ cookies, fetch }) => {
		try {
			const token = cookies.get('access_token')
			const setup = await createApi(fetch, token).auth.totp.setup()
			return { totpSetup: setup }
		} catch (err) {
			const status = err instanceof ApiError && err.status === 409 ? 409 : 400
			return fail(status, { totpError: err instanceof Error ? err.message : 'Failed to start TOTP setup' })
		}
	},

	// Confirms TOTP setup with a valid authenticator code, activating TOTP.
	totpConfirm: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const code = data.get('code') as string

		if (!code?.trim()) return fail(400, { totpError: 'Code is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).auth.totp.confirm(code.trim())
			return { totpEnabled: true }
		} catch (err) {
			return fail(400, { totpError: err instanceof Error ? err.message : 'Invalid code' })
		}
	},

	// Disables TOTP by validating the current authenticator code.
	totpDisable: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const code = data.get('code') as string

		if (!code?.trim()) return fail(400, { totpError: 'Code is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).auth.totp.disable(code.trim())
			return { totpDisabled: true }
		} catch (err) {
			return fail(400, { totpError: err instanceof Error ? err.message : 'Invalid code' })
		}
	},
}

import { fail, error } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi, ApiError } from '$lib/api'

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const token = cookies.get('access_token')
	try {
		const projects = await createApi(fetch, token).projects.list()
		return { projects }
	} catch (err) {
		if (err instanceof ApiError) error(err.status, err.message)
		error(500, 'Failed to load projects')
	}
}

export const actions: Actions = {
	create: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const name = data.get('name') as string

		if (!name?.trim()) return fail(400, { error: 'Name is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).projects.create(name.trim())
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Failed to create project' })
		}
	},

	delete: async ({ request, cookies, fetch }) => {
		const data = await request.formData()
		const slug = data.get('slug') as string

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).projects.delete(slug)
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Failed to delete project' })
		}
	},
}

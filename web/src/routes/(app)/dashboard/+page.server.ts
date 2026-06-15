import { fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'
import { loadOrError } from '$lib/server/load'

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const token = cookies.get('access_token')
	const projects = await loadOrError(
		() => createApi(fetch, token).projects.list(),
		'Failed to load projects',
	)
	return { projects }
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

import { fail, error } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
	const token = cookies.get('access_token')
	const api = createApi(fetch, token)

	try {
		const links = await api.links.list(params.slug)
		return { slug: params.slug, links }
	} catch {
		error(404, 'Project not found')
	}
}

export const actions: Actions = {
	create: async ({ request, params, cookies, fetch }) => {
		const data = await request.formData()
		const url = data.get('url') as string
		const title = (data.get('title') as string) || undefined

		if (!url?.trim()) return fail(400, { error: 'URL is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).links.create(params.slug, { url: url.trim(), title })
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Failed to create link' })
		}
	},

	delete: async ({ request, params, cookies, fetch }) => {
		const data = await request.formData()
		const code = data.get('code') as string

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).links.delete(params.slug, code)
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Failed to delete link' })
		}
	},
}

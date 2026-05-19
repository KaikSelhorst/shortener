import { fail, error } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'

export const load: PageServerLoad = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const cursor = url.searchParams.get('cursor') ?? undefined

	try {
		const links = await createApi(fetch, token).links.list(params.slug, cursor)
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
		const description = (data.get('description') as string) || undefined
		const og_image = (data.get('og_image') as string) || undefined
		const expires_at_raw = data.get('expires_at') as string
		const expires_at = expires_at_raw ? new Date(expires_at_raw).toISOString() : undefined

		if (!url?.trim()) return fail(400, { error: 'URL is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).links.create(params.slug, {
				url: url.trim(),
				title,
				description,
				og_image,
				expires_at,
			})
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

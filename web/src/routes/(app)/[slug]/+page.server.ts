import { fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'
import { loadOrError } from '$lib/server/load'

export const load: PageServerLoad = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const cursor = url.searchParams.get('cursor') ?? undefined
	const links = await loadOrError(
		() => createApi(fetch, token).links.list(params.slug, cursor),
		'Failed to load links',
	)
	return { slug: params.slug, links }
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
		const custom_code = (data.get('custom_code') as string) || undefined

		if (!url?.trim()) return fail(400, { error: 'URL is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).links.create(params.slug, {
				url: url.trim(),
				title,
				description,
				og_image,
				expires_at,
				custom_code,
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

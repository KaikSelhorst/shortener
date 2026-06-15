import { fail } from '@sveltejs/kit'
import type { Actions, PageServerLoad } from './$types'
import { createApi } from '$lib/api'
import { loadOrError } from '$lib/server/load'

export const load: PageServerLoad = async ({ params, cookies, fetch }) => {
	const token = cookies.get('access_token')
	const webhooks = await loadOrError(
		() => createApi(fetch, token).webhooks.list(params.slug),
		'Failed to load webhooks',
	)
	return { slug: params.slug, webhooks }
}

export const actions: Actions = {
	create: async ({ request, params, cookies, fetch }) => {
		const data = await request.formData()
		const url = data.get('url') as string
		const events = data.getAll('events') as string[]

		if (!url?.trim()) return fail(400, { createError: 'URL is required' })
		if (!events.length) return fail(400, { createError: 'Select at least one event' })

		try {
			const token = cookies.get('access_token')
			const webhook = await createApi(fetch, token).webhooks.create(params.slug, {
				url: url.trim(),
				events,
			})
			return { created: webhook }
		} catch (err) {
			return fail(400, { createError: err instanceof Error ? err.message : 'Failed to create webhook' })
		}
	},

	delete: async ({ request, params, cookies, fetch }) => {
		const data = await request.formData()
		const id = data.get('id') as string

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).webhooks.delete(params.slug, id)
		} catch (err) {
			return fail(400, { deleteError: err instanceof Error ? err.message : 'Failed to delete webhook' })
		}
	},
}

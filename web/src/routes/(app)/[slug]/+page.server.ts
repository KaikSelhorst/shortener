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

const UTM_KEYS = ['utm_source', 'utm_medium', 'utm_campaign', 'utm_term', 'utm_content'] as const

function appendUtm(base: string, data: FormData): string {
	const params = new URLSearchParams()
	for (const key of UTM_KEYS) {
		const val = (data.get(key) as string | null)?.trim()
		if (val) params.set(key, val)
	}
	const qs = params.toString()
	if (!qs) return base
	return base + (base.includes('?') ? '&' : '?') + qs
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
		const max_clicks_raw = data.get('max_clicks') as string
		const max_clicks = max_clicks_raw ? parseInt(max_clicks_raw, 10) : undefined
		const custom_code = (data.get('custom_code') as string) || undefined

		if (!url?.trim()) return fail(400, { error: 'URL is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).links.create(params.slug, {
				url: appendUtm(url.trim(), data),
				title,
				description,
				og_image,
				expires_at,
				max_clicks,
				custom_code,
			})
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Failed to create link' })
		}
	},

	update: async ({ request, params, cookies, fetch }) => {
		const data = await request.formData()
		const code = data.get('code') as string
		const url = data.get('url') as string
		const title = (data.get('title') as string) || undefined
		const expires_at_raw = data.get('expires_at') as string
		const expires_at = expires_at_raw ? new Date(expires_at_raw).toISOString() : undefined
		const max_clicks_raw = data.get('max_clicks') as string
		const max_clicks = max_clicks_raw ? parseInt(max_clicks_raw, 10) : undefined

		if (!url?.trim()) return fail(400, { error: 'URL is required' })

		try {
			const token = cookies.get('access_token')
			await createApi(fetch, token).links.update(params.slug, code, {
				url: appendUtm(url.trim(), data),
				title,
				expires_at,
				max_clicks,
			})
		} catch (err) {
			return fail(400, { error: err instanceof Error ? err.message : 'Failed to update link' })
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

import { error } from '@sveltejs/kit'
import type { PageServerLoad } from './$types'
import { createApi } from '$lib/api'
import { loadOrError } from '$lib/server/load'

export const load: PageServerLoad = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const api = createApi(fetch, token)
	const webhookId = Number(params.id)
	if (!Number.isInteger(webhookId) || webhookId < 1) error(404, 'Webhook not found')
	const page = Math.max(1, Number(url.searchParams.get('page')) || 1)

	const [webhooks, deliveriesRes] = await Promise.all([
		loadOrError(() => api.webhooks.list(params.slug), 'Failed to load webhooks'),
		loadOrError(
			() => api.webhooks.deliveries(params.slug, webhookId, page),
			'Failed to load deliveries',
		),
	])

	const webhook = webhooks.find((w) => w.id === webhookId)
	if (!webhook) error(404, 'Webhook not found')

	return {
		slug: params.slug,
		webhook,
		deliveries: deliveriesRes.data,
		hasMore: deliveriesRes.has_more,
		page,
	}
}

import { createApi } from '$lib/api'
import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'

export const GET: RequestHandler = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const limit = Number(url.searchParams.get('limit') ?? '20')
	const deliveries = await createApi(fetch, token).webhooks.deliveries(
		params.slug,
		Number(params.id),
		limit,
	)
	return json(deliveries)
}

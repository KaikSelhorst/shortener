import type { PageServerLoad } from './$types'
import { createApi } from '$lib/api'
import { parsePeriod } from '$lib/analytics'
import { loadOrError } from '$lib/server/load'

export const load: PageServerLoad = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const period = parsePeriod(url.searchParams.get('period'))
	const analytics = await loadOrError(
		() => createApi(fetch, token).analytics.project(params.slug, period),
		'Failed to load analytics',
	)
	return { slug: params.slug, analytics, period }
}

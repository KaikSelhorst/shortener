import { error } from '@sveltejs/kit'
import type { PageServerLoad } from './$types'
import { createApi, ApiError } from '$lib/api'
import { parsePeriod } from '$lib/analytics'

export const load: PageServerLoad = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const period = parsePeriod(url.searchParams.get('period'))

	try {
		const analytics = await createApi(fetch, token).analytics.link(
			params.slug,
			params.code,
			period,
		)
		return { slug: params.slug, code: params.code, analytics, period }
	} catch (err) {
		if (err instanceof ApiError) error(err.status, err.message)
		error(500, 'Failed to load analytics')
	}
}

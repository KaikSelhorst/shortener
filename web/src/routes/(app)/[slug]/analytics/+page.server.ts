import { error } from '@sveltejs/kit'
import type { PageServerLoad } from './$types'
import { createApi, ApiError } from '$lib/api'

const VALID_PERIODS = ['7d', '30d', '90d'] as const
type Period = (typeof VALID_PERIODS)[number]

export const load: PageServerLoad = async ({ params, cookies, fetch, url }) => {
	const token = cookies.get('access_token')
	const raw = url.searchParams.get('period') ?? '30d'
	const period: Period = (VALID_PERIODS as readonly string[]).includes(raw)
		? (raw as Period)
		: '30d'

	try {
		const analytics = await createApi(fetch, token).analytics.project(params.slug, period)
		return { slug: params.slug, analytics, period }
	} catch (err) {
		if (err instanceof ApiError) error(err.status, err.message)
		error(500, 'Failed to load analytics')
	}
}

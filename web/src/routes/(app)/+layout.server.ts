import { redirect } from '@sveltejs/kit'
import type { LayoutServerLoad } from './$types'
import { createApi } from '$lib/api'
import { loadOrError } from '$lib/server/load'

export const load: LayoutServerLoad = async ({ locals, cookies, fetch, depends }) => {
	if (!locals.user) redirect(302, '/login')

	depends('app:projects')
	const token = cookies.get('access_token')
	const projects = await loadOrError(
		() => createApi(fetch, token).projects.list(),
		'Failed to load projects',
	)

	return { user: locals.user, projects }
}

import { redirect } from '@sveltejs/kit'
import type { PageServerLoad } from './$types'

export const load: PageServerLoad = async ({ parent }) => {
	const { projects } = await parent()
	if (projects.length > 0) redirect(302, `/${projects[0].slug}`)
}

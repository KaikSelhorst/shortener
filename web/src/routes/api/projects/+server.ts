import { json } from '@sveltejs/kit'
import type { RequestHandler } from './$types'
import { createApi, ApiError } from '$lib/api'

export const POST: RequestHandler = async ({ request, cookies, fetch }) => {
	const { name } = (await request.json()) as { name?: string }

	if (!name?.trim()) return json({ error: 'Name is required' }, { status: 400 })

	try {
		const token = cookies.get('access_token')
		const project = await createApi(fetch, token).projects.create(name.trim())
		return json(project)
	} catch (err) {
		return json({ error: err instanceof ApiError ? err.message : 'Failed to create project' }, { status: 400 })
	}
}

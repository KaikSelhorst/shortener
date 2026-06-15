import { env } from '$env/dynamic/private'
import type { RequestHandler } from './$types'

const BASE = env.API_URL ?? 'http://localhost:8080'

export const GET: RequestHandler = async ({ cookies, request }) => {
	const token = cookies.get('access_token')
	if (!token) return new Response(null, { status: 401 })

	const upstream = await fetch(`${BASE}/projects/stream`, {
		headers: { Authorization: `Bearer ${token}` },
		signal: request.signal,
	})
	if (!upstream.ok) return new Response(null, { status: upstream.status })

	return new Response(upstream.body, {
		headers: {
			'Content-Type': 'text/event-stream',
			'Cache-Control': 'no-cache',
			Connection: 'keep-alive',
		},
	})
}

import type { Cookies } from '@sveltejs/kit'
import { API_URL } from './config'

export async function proxySSE(
	path: string,
	cookies: Cookies,
	signal: AbortSignal,
): Promise<Response> {
	const token = cookies.get('access_token')
	if (!token) return new Response(null, { status: 401 })

	const upstream = await fetch(`${API_URL}${path}`, {
		headers: { Authorization: `Bearer ${token}` },
		signal,
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

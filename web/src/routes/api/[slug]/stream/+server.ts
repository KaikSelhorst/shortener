import { proxySSE } from '$lib/server/sse'
import type { RequestHandler } from './$types'

export const GET: RequestHandler = ({ params, cookies, request }) =>
	proxySSE(`/projects/${params.slug}/stream`, cookies, request.signal)

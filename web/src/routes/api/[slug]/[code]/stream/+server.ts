import { proxySSE } from '$lib/server/sse'
import type { RequestHandler } from './$types'

export const GET: RequestHandler = ({ params, cookies, request }) =>
	proxySSE(`/projects/${params.slug}/links/${params.code}/stream`, cookies, request.signal)

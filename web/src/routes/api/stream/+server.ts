import { proxySSE } from '$lib/server/sse'
import type { RequestHandler } from './$types'

export const GET: RequestHandler = ({ cookies, request }) =>
	proxySSE('/projects/stream', cookies, request.signal)

import { API_URL as BASE } from '$lib/server/config'

export type Fetch = typeof globalThis.fetch
export type Req = <T>(method: string, path: string, body?: unknown) => Promise<T>

/** Thrown by the API client when the server returns a non-2xx status. */
export class ApiError extends Error {
	constructor(
		public readonly status: number,
		message: string,
	) {
		super(message)
		this.name = 'ApiError'
	}
}

export function createRequester(fetch: Fetch, token?: string): Req {
	const headers: Record<string, string> = { 'Content-Type': 'application/json' }
	if (token) headers['Authorization'] = `Bearer ${token}`

	return async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
		const res = await fetch(`${BASE}${path}`, {
			method,
			headers,
			body: body !== undefined ? JSON.stringify(body) : undefined,
		})
		if (!res.ok) {
			let message = `${res.status} ${res.statusText}`
			try {
				const data = await res.json()
				if (typeof data.error === 'string') message = data.error
			} catch {
				/* ignore parse error — keep status-based message */
			}
			throw new ApiError(res.status, message)
		}
		if (res.status === 204) return undefined as T
		return res.json() as Promise<T>
	}
}

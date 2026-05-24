import { env } from '$env/dynamic/private'
import type {
	TokenResponse,
	Project,
	Link,
	ListLinksResponse,
	CreateLinkRequest,
	UpdateLinkRequest,
	ApiKey,
	CreateApiKeyRequest,
	CreateApiKeyResponse,
} from './types'

const BASE = env.API_URL ?? 'http://localhost:8080'

type Fetch = typeof globalThis.fetch

export function createApi(fetch: Fetch, token?: string) {
	const headers: Record<string, string> = { 'Content-Type': 'application/json' }
	if (token) headers['Authorization'] = `Bearer ${token}`

	async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
		const res = await fetch(`${BASE}${path}`, {
			method,
			headers,
			body: body !== undefined ? JSON.stringify(body) : undefined,
		})
		if (!res.ok) {
			const text = await res.text()
			throw new Error(text || `${res.status} ${res.statusText}`)
		}
		if (res.status === 204) return undefined as T
		return res.json() as Promise<T>
	}

	return {
		auth: {
			login: (email: string, password: string) =>
				req<TokenResponse>('POST', '/auth/login', { email, password }),
			register: (email: string, password: string) =>
				req<TokenResponse>('POST', '/auth/register', { email, password }),
			refresh: (refresh_token: string) =>
				req<TokenResponse>('POST', '/auth/refresh', { refresh_token }),
			logout: (refresh_token: string) =>
				req<void>('POST', '/auth/logout', { refresh_token }),
		},
		projects: {
			list: () => req<Project[]>('GET', '/projects'),
			create: (name: string) => req<Project>('POST', '/projects', { name }),
			update: (slug: string, name: string) =>
				req<Project>('PUT', `/projects/${slug}`, { name }),
			delete: (slug: string) => req<void>('DELETE', `/projects/${slug}`),
		},
		apiKeys: {
			list: () => req<ApiKey[]>('GET', '/api-keys'),
			create: (data: CreateApiKeyRequest) =>
				req<CreateApiKeyResponse>('POST', '/api-keys', data),
			delete: (id: number) => req<void>('DELETE', `/api-keys/${id}`),
		},
		links: {
			list: (slug: string, cursor?: string) =>
				req<ListLinksResponse>(
					'GET',
					`/projects/${slug}/links${cursor ? `?cursor=${cursor}` : ''}`,
				),
			create: (slug: string, data: CreateLinkRequest) =>
				req<Link>('POST', `/projects/${slug}/links`, data),
			get: (slug: string, code: string) =>
				req<Link>('GET', `/projects/${slug}/links/${code}`),
			update: (slug: string, code: string, data: UpdateLinkRequest) =>
				req<Link>('PUT', `/projects/${slug}/links/${code}`, data),
			delete: (slug: string, code: string) =>
				req<void>('DELETE', `/projects/${slug}/links/${code}`),
		},
	}
}

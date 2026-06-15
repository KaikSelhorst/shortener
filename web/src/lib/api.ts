import { API_URL as BASE } from '$lib/server/config'
import type {
	AuthState,
	TOTPSetupResponse,
	MeResponse,
	Project,
	Link,
	ListLinksResponse,
	LinkRequest,
	ApiKeyResponse,
	CreateApiKeyRequest,
	CreateApiKeyResponse,
	Webhook,
	CreateWebhookRequest,
	CreateWebhookResponse,
	WebhookDeliveriesResponse,
	LinkAnalytics,
	ProjectAnalytics,
} from './types'


type Fetch = typeof globalThis.fetch

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

	return {
		auth: {
			login: (email: string, password: string) =>
				req<AuthState>('POST', '/auth/login', { email, password }),
			register: (email: string, password: string) =>
				req<AuthState>('POST', '/auth/register', { email, password }),
			refresh: (refresh_token: string) =>
				req<AuthState>('POST', '/auth/refresh', { refresh_token }),
			logout: (refresh_token: string) =>
				req<void>('POST', '/auth/logout', { refresh_token }),
			me: () => req<MeResponse>('GET', '/auth/me'),
			totp: {
				validateMFA: (session: string, code: string) =>
					req<AuthState>('POST', '/auth/mfa/totp', { session, code }),
				setup: () => req<TOTPSetupResponse>('POST', '/auth/totp/setup'),
				confirm: (code: string) => req<void>('POST', '/auth/totp/confirm', { code }),
				disable: (code: string) => req<void>('DELETE', '/auth/totp', { code }),
			},
		},
		projects: {
			list: () => req<Project[]>('GET', '/projects'),
			create: (name: string) => req<Project>('POST', '/projects', { name }),
			update: (slug: string, name: string) =>
				req<Project>('PUT', `/projects/${slug}`, { name }),
			delete: (slug: string) => req<void>('DELETE', `/projects/${slug}`),
		},
		apiKeys: {
			list: () => req<ApiKeyResponse[]>('GET', '/api-keys'),
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
			create: (slug: string, data: LinkRequest) =>
				req<Link>('POST', `/projects/${slug}/links`, data),
			get: (slug: string, code: string) =>
				req<Link>('GET', `/projects/${slug}/links/${code}`),
			update: (slug: string, code: string, data: LinkRequest) =>
				req<Link>('PUT', `/projects/${slug}/links/${code}`, data),
			delete: (slug: string, code: string) =>
				req<void>('DELETE', `/projects/${slug}/links/${code}`),
		},
		webhooks: {
			list: (slug: string) => req<Webhook[]>('GET', `/projects/${slug}/webhooks`),
			create: (slug: string, data: CreateWebhookRequest) =>
				req<CreateWebhookResponse>('POST', `/projects/${slug}/webhooks`, data),
			delete: (slug: string, id: number) => req<void>('DELETE', `/projects/${slug}/webhooks/${id}`),
			deliveries: (slug: string, webhookId: number, page = 1, limit = 20) =>
				req<WebhookDeliveriesResponse>('GET', `/projects/${slug}/webhooks/${webhookId}/deliveries?limit=${limit}&page=${page}`),
		},
		analytics: {
			project: (slug: string, period?: string) =>
				req<ProjectAnalytics>(
					'GET',
					`/projects/${slug}/analytics${period ? `?period=${period}` : ''}`,
				),
			link: (slug: string, code: string, period?: string) =>
				req<LinkAnalytics>(
					'GET',
					`/projects/${slug}/links/${code}/analytics${period ? `?period=${period}` : ''}`,
				),
		},
	}
}

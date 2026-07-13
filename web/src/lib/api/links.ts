import type { Req } from './client'
import type { Link, ListLinksResponse, LinkRequest } from '../types'

export function createLinksApi(req: Req) {
	return {
		list: (slug: string, cursor?: string) =>
			req<ListLinksResponse>(
				'GET',
				`/projects/${slug}/links${cursor ? `?cursor=${cursor}` : ''}`,
			),
		create: (slug: string, data: LinkRequest) => req<Link>('POST', `/projects/${slug}/links`, data),
		get: (slug: string, code: string) => req<Link>('GET', `/projects/${slug}/links/${code}`),
		update: (slug: string, code: string, data: LinkRequest) =>
			req<Link>('PUT', `/projects/${slug}/links/${code}`, data),
		delete: (slug: string, code: string) => req<void>('DELETE', `/projects/${slug}/links/${code}`),
	}
}

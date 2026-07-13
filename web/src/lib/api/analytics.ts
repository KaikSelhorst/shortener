import type { Req } from './client'
import type { LinkAnalytics, ProjectAnalytics } from '../types'

export function createAnalyticsApi(req: Req) {
	return {
		project: (slug: string, period?: string) =>
			req<ProjectAnalytics>('GET', `/projects/${slug}/analytics${period ? `?period=${period}` : ''}`),
		link: (slug: string, code: string, period?: string) =>
			req<LinkAnalytics>(
				'GET',
				`/projects/${slug}/links/${code}/analytics${period ? `?period=${period}` : ''}`,
			),
	}
}

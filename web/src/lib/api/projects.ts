import type { Req } from './client'
import type { Project } from '../types'

export function createProjectsApi(req: Req) {
	return {
		list: () => req<Project[]>('GET', '/projects'),
		create: (name: string) => req<Project>('POST', '/projects', { name }),
		update: (slug: string, name: string) => req<Project>('PUT', `/projects/${slug}`, { name }),
		delete: (slug: string) => req<void>('DELETE', `/projects/${slug}`),
	}
}

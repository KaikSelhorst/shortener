import { createRequester, type Fetch } from './client'
import { createAuthApi } from './auth'
import { createProjectsApi } from './projects'
import { createApiKeysApi } from './api-keys'
import { createLinksApi } from './links'
import { createWebhooksApi } from './webhooks'
import { createAnalyticsApi } from './analytics'

export { ApiError } from './client'

export function createApi(fetch: Fetch, token?: string) {
	const req = createRequester(fetch, token)

	return {
		auth: createAuthApi(req),
		projects: createProjectsApi(req),
		apiKeys: createApiKeysApi(req),
		links: createLinksApi(req),
		webhooks: createWebhooksApi(req),
		analytics: createAnalyticsApi(req),
	}
}

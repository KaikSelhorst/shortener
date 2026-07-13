import type { Req } from './client'
import type { Webhook, CreateWebhookRequest, CreateWebhookResponse, WebhookDeliveriesResponse } from '../types'

export function createWebhooksApi(req: Req) {
	return {
		list: (slug: string) => req<Webhook[]>('GET', `/projects/${slug}/webhooks`),
		create: (slug: string, data: CreateWebhookRequest) =>
			req<CreateWebhookResponse>('POST', `/projects/${slug}/webhooks`, data),
		delete: (slug: string, id: string) => req<void>('DELETE', `/projects/${slug}/webhooks/${id}`),
		deliveries: (slug: string, webhookId: string, page = 1, limit = 20) =>
			req<WebhookDeliveriesResponse>(
				'GET',
				`/projects/${slug}/webhooks/${webhookId}/deliveries?limit=${limit}&page=${page}`,
			),
	}
}

import type { Req } from './client'
import type { ApiKeyResponse, CreateApiKeyRequest, CreateApiKeyResponse } from '../types'

export function createApiKeysApi(req: Req) {
	return {
		list: () => req<ApiKeyResponse[]>('GET', '/api-keys'),
		create: (data: CreateApiKeyRequest) => req<CreateApiKeyResponse>('POST', '/api-keys', data),
		delete: (id: number) => req<void>('DELETE', `/api-keys/${id}`),
	}
}

export interface TokenResponse {
	access_token: string
	refresh_token: string
	token_type: string
	expires_in: number
}

export interface Project {
	id: number
	user_id: number
	name: string
	slug: string
	created_at: string
}

export interface Link {
	id: number
	project_id: number
	short_code: string
	original_url: string
	title: string | null
	description: string | null
	og_image: string | null
	expires_at: string | null
	created_at: string
	short_url: string
}

export interface ListLinksResponse {
	data: Link[]
	next_cursor: string | null
	prev_cursor: string | null
	limit: number
}

export interface CreateLinkRequest {
	url: string
	title?: string
	description?: string
	og_image?: string
	expires_at?: string
}

export interface UpdateLinkRequest {
	url: string
	title?: string
	description?: string
	og_image?: string
	expires_at?: string
}

export interface ApiKey {
	id: number
	user_id: number
	project_id: number | null
	name: string
	key_prefix: string
	scopes: string[]
	last_used_at: string | null
	created_at: string
}

export interface CreateApiKeyRequest {
	name: string
	scopes: string[]
	project_id?: number
}

export interface CreateApiKeyResponse extends ApiKey {
	token: string
}

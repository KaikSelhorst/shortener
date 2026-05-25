export interface TokenResponse {
	access_token: string
	refresh_token: string
	token_type: string
	expires_in: number
}

export interface Project {
	id: number
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

// LinkRequest is used for both creating and updating a link — mirrors the backend's LinkRequest DTO.
export interface LinkRequest {
	url: string
	title?: string
	description?: string
	og_image?: string
	expires_at?: string
}

export interface ApiKeyResponse {
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

export interface CreateApiKeyResponse extends ApiKeyResponse {
	token: string
}

export interface TokenResponse {
	access_token: string
	refresh_token: string
	token_type: string
	expires_in: number
}

// AuthState is the unified state machine response returned by login and register.
// next === 'complete' → tokens are present and the user is authenticated.
// next === 'totp'     → a second factor is required; session must be passed to /auth/mfa/totp.
export interface AuthState {
	next: 'complete' | 'totp'
	session?: string
	access_token?: string
	refresh_token?: string
	token_type?: string
	expires_in?: number
}

export interface TOTPSetupResponse {
	uri: string
	secret: string
}

export interface MeResponse {
	id: number
	email: string
	totp_enabled: boolean
	created_at: string
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

export interface ClicksOverTime {
	date: string
	count: number
}

export interface DeviceBreakdown {
	mobile: number
	desktop: number
	tablet: number
	bot: number
	unknown: number
}

export interface BrowserBreakdown {
	chrome: number
	firefox: number
	safari: number
	edge: number
	opera: number
	samsung: number
	ie: number
	other: number
	unknown: number
}

export interface ReferrerBreakdown {
	direct: number
	instagram: number
	facebook: number
	twitter: number
	tiktok: number
	linkedin: number
	whatsapp: number
	youtube: number
	google: number
	discord: number
	other: number
}

export interface TopLink {
	short_code: string
	original_url: string
	title: string | null
	total_clicks: number
}

export interface LinkAnalytics {
	link_id: number
	short_code: string
	total_clicks: number
	unique_clicks: number
	over_time: ClicksOverTime[]
	devices: DeviceBreakdown
	referrers: ReferrerBreakdown
	browsers: BrowserBreakdown
}

export interface ProjectAnalytics {
	total_clicks: number
	unique_clicks: number
	over_time: ClicksOverTime[]
	devices: DeviceBreakdown
	referrers: ReferrerBreakdown
	browsers: BrowserBreakdown
	top_links: TopLink[]
}

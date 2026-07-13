import type { Req } from './client'
import type { AuthState, TOTPSetupResponse, MeResponse } from '../types'

export function createAuthApi(req: Req) {
	return {
		login: (email: string, password: string) =>
			req<AuthState>('POST', '/auth/login', { email, password }),
		register: (email: string, password: string) =>
			req<AuthState>('POST', '/auth/register', { email, password }),
		refresh: (refresh_token: string) =>
			req<AuthState>('POST', '/auth/refresh', { refresh_token }),
		logout: (refresh_token: string) => req<void>('POST', '/auth/logout', { refresh_token }),
		me: () => req<MeResponse>('GET', '/auth/me'),
		totp: {
			validateMFA: (session: string, code: string) =>
				req<AuthState>('POST', '/auth/mfa/totp', { session, code }),
			setup: () => req<TOTPSetupResponse>('POST', '/auth/totp/setup'),
			confirm: (code: string) => req<void>('POST', '/auth/totp/confirm', { code }),
			disable: (code: string) => req<void>('DELETE', '/auth/totp', { code }),
		},
	}
}

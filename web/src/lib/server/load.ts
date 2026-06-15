import { error } from '@sveltejs/kit'
import { ApiError } from '$lib/api'

export async function loadOrError<T>(fn: () => Promise<T>, fallback = 'Failed to load'): Promise<T> {
	try {
		return await fn()
	} catch (err) {
		if (err instanceof ApiError) error(err.status, err.message)
		error(500, fallback)
	}
}

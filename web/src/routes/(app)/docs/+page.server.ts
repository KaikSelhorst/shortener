import { API_URL } from '$lib/server/config'
import { error } from '@sveltejs/kit'

export async function load() {
	try {
		const res = await fetch(`${API_URL}/openapi.json`)
		if (!res.ok) error(502, 'API spec unavailable')
		const spec = await res.json()
		return { spec }
	} catch (e) {
		error(502, 'Could not reach API server')
	}
}

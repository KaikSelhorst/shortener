const dateFormat = new Intl.DateTimeFormat(undefined, { dateStyle: 'short' })
const dateTimeFormat = new Intl.DateTimeFormat(undefined, { dateStyle: 'short', timeStyle: 'short' })

export function formatDate(value: string | null | undefined): string {
	if (!value) return '—'
	return dateFormat.format(new Date(value))
}

export function formatDateTime(value: string | null | undefined): string {
	if (!value) return '—'
	return dateTimeFormat.format(new Date(value))
}

export function toDatetimeLocal(iso: string | null): string {
	if (!iso) return ''
	return iso.slice(0, 16)
}

const UTM_KEYS = ['utm_source', 'utm_medium', 'utm_campaign', 'utm_term', 'utm_content'] as const

export function parseUtm(rawUrl: string) {
	try {
		const u = new URL(rawUrl)
		const utm = {
			source: u.searchParams.get('utm_source') ?? '',
			medium: u.searchParams.get('utm_medium') ?? '',
			campaign: u.searchParams.get('utm_campaign') ?? '',
			term: u.searchParams.get('utm_term') ?? '',
			content: u.searchParams.get('utm_content') ?? '',
		}
		for (const key of UTM_KEYS) u.searchParams.delete(key)
		return { baseUrl: u.toString(), ...utm }
	} catch {
		return { baseUrl: rawUrl, source: '', medium: '', campaign: '', term: '', content: '' }
	}
}

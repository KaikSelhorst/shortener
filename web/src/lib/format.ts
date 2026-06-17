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

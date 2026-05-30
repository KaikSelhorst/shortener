import type { DeviceBreakdown, ReferrerBreakdown, BrowserBreakdown } from './types'

export const VALID_PERIODS = ['7d', '30d', '90d'] as const
export type Period = (typeof VALID_PERIODS)[number]

export const periods: { value: Period; label: string }[] = [
	{ value: '7d', label: '7 days' },
	{ value: '30d', label: '30 days' },
	{ value: '90d', label: '90 days' },
]

export function parsePeriod(raw: string | null): Period {
	return (VALID_PERIODS as readonly string[]).includes(raw ?? '')
		? (raw as Period)
		: '30d'
}

export function deviceItems(d: DeviceBreakdown) {
	return [
		{ label: 'Mobile', value: d.mobile },
		{ label: 'Desktop', value: d.desktop },
		{ label: 'Tablet', value: d.tablet },
		{ label: 'Bot', value: d.bot },
		{ label: 'Unknown', value: d.unknown },
	]
		.filter((i) => i.value > 0)
		.sort((a, b) => b.value - a.value)
}

export function browserItems(b: BrowserBreakdown) {
	return [
		{ label: 'Chrome', value: b.chrome },
		{ label: 'Safari', value: b.safari },
		{ label: 'Firefox', value: b.firefox },
		{ label: 'Edge', value: b.edge },
		{ label: 'Opera', value: b.opera },
		{ label: 'Samsung', value: b.samsung },
		{ label: 'IE', value: b.ie },
		{ label: 'Other', value: b.other },
		{ label: 'Unknown', value: b.unknown },
	]
		.filter((i) => i.value > 0)
		.sort((a, b) => b.value - a.value)
}

export function referrerItems(r: ReferrerBreakdown) {
	return [
		{ label: 'Direct', value: r.direct },
		{ label: 'Google', value: r.google },
		{ label: 'Instagram', value: r.instagram },
		{ label: 'Facebook', value: r.facebook },
		{ label: 'Twitter/X', value: r.twitter },
		{ label: 'TikTok', value: r.tiktok },
		{ label: 'Discord', value: r.discord },
		{ label: 'LinkedIn', value: r.linkedin },
		{ label: 'WhatsApp', value: r.whatsapp },
		{ label: 'YouTube', value: r.youtube },
		{ label: 'Other', value: r.other },
	]
		.filter((i) => i.value > 0)
		.sort((a, b) => b.value - a.value)
}

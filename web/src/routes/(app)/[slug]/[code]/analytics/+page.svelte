<script lang="ts">
	import type { PageData } from './$types'
	import type { DeviceBreakdown, ReferrerBreakdown } from '$lib/types'
	import { BarChart, LineChart } from '$lib'

	let { data }: { data: PageData } = $props()

	const periods = [
		{ value: '7d', label: '7 days' },
		{ value: '30d', label: '30 days' },
		{ value: '90d', label: '90 days' },
	]

	function deviceItems(d: DeviceBreakdown) {
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

	function referrerItems(r: ReferrerBreakdown) {
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
</script>

<!-- breadcrumb -->
<div class="flex items-center gap-2 text-sm">
	<a href="/dashboard" class="text-muted-foreground hover:text-foreground">Projects</a>
	<span class="text-muted-foreground">/</span>
	<a href="/{data.slug}" class="text-muted-foreground hover:text-foreground">{data.slug}</a>
	<span class="text-muted-foreground">/</span>
	<span class="font-mono text-muted-foreground hover:text-foreground">
		<a href="/{data.slug}">{data.code}</a>
	</span>
	<span class="text-muted-foreground">/</span>
	<span class="font-semibold text-foreground">Analytics</span>
</div>

<!-- header + period selector -->
<div class="mt-4 flex items-center justify-end gap-4">
	<div class="flex rounded-md bg-muted p-1">
		{#each periods as p (p.value)}
			<a
				href="?period={p.value}"
				class="rounded px-3 py-1 text-sm transition-colors
					{data.period === p.value
					? 'bg-background font-medium text-foreground shadow-sm'
					: 'text-muted-foreground hover:text-foreground'}"
			>
				{p.label}
			</a>
		{/each}
	</div>
</div>

<!-- stat cards -->
<div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3">
	<div class="rounded-md bg-card p-4 shadow-sm">
		<p class="text-xs text-muted-foreground">Total Clicks</p>
		<p class="mt-1 text-2xl font-semibold text-foreground">
			{data.analytics.total_clicks.toLocaleString()}
		</p>
	</div>
	<div class="rounded-md bg-card p-4 shadow-sm">
		<p class="text-xs text-muted-foreground">Unique Visitors</p>
		<p class="mt-1 text-2xl font-semibold text-foreground">
			{data.analytics.unique_clicks.toLocaleString()}
		</p>
	</div>
	<div class="col-span-2 rounded-md bg-card p-4 shadow-sm sm:col-span-1">
		<p class="text-xs text-muted-foreground">Avg. per Day</p>
		<p class="mt-1 text-2xl font-semibold text-foreground">
			{data.analytics.over_time.length > 0
				? Math.round(
						data.analytics.total_clicks /
							Math.max(data.analytics.over_time.length, 1),
					).toLocaleString()
				: '0'}
		</p>
	</div>
</div>

<!-- clicks over time -->
<div class="mt-4 rounded-md bg-card p-4 shadow-sm">
	<h2 class="mb-3 text-sm font-semibold text-foreground">Clicks over time</h2>
	<LineChart data={data.analytics.over_time} />
</div>

<!-- device + referrer breakdowns -->
<div class="mt-4 grid grid-cols-1 gap-4 lg:grid-cols-2">
	<div class="rounded-md bg-card p-4 shadow-sm">
		<h2 class="mb-4 text-sm font-semibold text-foreground">Devices</h2>
		<BarChart items={deviceItems(data.analytics.devices)} />
	</div>
	<div class="rounded-md bg-card p-4 shadow-sm">
		<h2 class="mb-4 text-sm font-semibold text-foreground">Referrers</h2>
		<BarChart items={referrerItems(data.analytics.referrers)} />
	</div>
</div>

<!-- back link -->
<div class="mt-6">
	<a href="/{data.slug}/analytics" class="text-sm text-muted-foreground hover:text-foreground">
		← Back to project analytics
	</a>
</div>

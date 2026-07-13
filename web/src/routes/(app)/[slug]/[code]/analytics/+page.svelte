<script lang="ts">
	import { untrack } from 'svelte'
	import type { PageData } from './$types'
	import { LineChart } from '$lib/components/charts'
	import { PeriodSelector, AnalyticsBreakdowns, StatCard } from '$lib/components/analytics'
	import { deviceItems, referrerItems, browserItems } from '$lib/analytics'
	import { useSSE } from '$lib/state/sse.svelte'

	let { data }: { data: PageData } = $props()

	const devices = $derived(deviceItems(data.analytics.devices))
	const referrers = $derived(referrerItems(data.analytics.referrers))
	const browsers = $derived(browserItems(data.analytics.browsers))

	let liveClicks = $state(untrack(() => data.analytics.total_clicks))
	$effect(() => { liveClicks = data.analytics.total_clicks })

	const avgPerDay = $derived(
		data.analytics.over_time.length > 0
			? Math.round(liveClicks / Math.max(data.analytics.over_time.length, 1))
			: 0
	)

	const sse = useSSE(() => `/api/${data.slug}/${data.code}/stream`, () => { liveClicks++ })
</script>

<div class="sticky top-0 z-10 bg-background border-b border-border px-4 h-11 flex items-center justify-between">
	<div class="flex items-center gap-1.5 text-sm min-w-0">
		{#if sse.connected}
			<span class="flex items-center gap-1.5 text-xs text-success shrink-0">
				<span class="w-1.5 h-1.5 rounded-full bg-success"></span>Live
			</span>
		{/if}
	</div>
	<PeriodSelector current={data.period} />
</div>

<!-- Stats strip -->
<div class="grid grid-cols-3 border-b border-border">
	<StatCard label="Total Clicks" value={liveClicks} color="success" />
	<StatCard label="Unique Visitors" value={data.analytics.unique_clicks} />
	<StatCard label="Avg / Day" value={avgPerDay} bordered={false} />
</div>

<!-- Chart -->
<div class="px-6 py-5 border-b border-border">
	<p class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground mb-4">Clicks over time</p>
	<LineChart data={data.analytics.over_time} />
</div>

<!-- Breakdowns -->
<div class="px-6 py-5 border-b border-border">
	<AnalyticsBreakdowns {devices} {browsers} {referrers} />
</div>

<div class="px-6 py-4">
	<a href="/{data.slug}/analytics" class="text-sm text-muted-foreground hover:text-foreground transition-colors">
		← Back to project analytics
	</a>
</div>

<script lang="ts">
	import { onMount, untrack } from 'svelte'
	import type { PageData } from './$types'
	import { LineChart, StatBox, PeriodSelector, AnalyticsBreakdowns } from '$lib'
	import { deviceItems, referrerItems, browserItems } from '$lib/analytics'

	let { data }: { data: PageData } = $props()

	const devices = $derived(deviceItems(data.analytics.devices))
	const referrers = $derived(referrerItems(data.analytics.referrers))
	const browsers = $derived(browserItems(data.analytics.browsers))

	let liveClicks = $state(untrack(() => data.analytics.total_clicks))
	let connected = $state(false)

	$effect(() => {
		liveClicks = data.analytics.total_clicks
	})

	const avgPerDay = $derived(
		data.analytics.over_time.length > 0
			? Math.round(liveClicks / Math.max(data.analytics.over_time.length, 1))
			: 0,
	)

	onMount(() => {
		const es = new EventSource(`/api/${data.slug}/${data.code}/stream`)
		es.onopen = () => (connected = true)
		es.onmessage = () => { liveClicks++ }
		es.onerror = () => { connected = false }
		return () => es.close()
	})
</script>

<div class="mx-auto max-w-6xl space-y-4">
	<div class="tui-panel">
		<div class="tui-panel-header justify-between">
			<div class="flex items-center gap-2">
				<a
					href="/dashboard"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
				>
					projects
				</a>
				<span class="text-muted-foreground">/</span>
				<a
					href="/{data.slug}"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
				>
					{data.slug}
				</a>
				<span class="text-muted-foreground">/</span>
				<span class="font-mono text-[10px] text-muted-foreground">{data.code}</span>
				<span class="text-muted-foreground">/</span>
				<span>▌ analytics</span>
			</div>
			<div class="flex items-center gap-3">
				{#if connected}
					<span class="font-mono text-[10px] uppercase tracking-wider text-tui-green">● live</span>
				{/if}
				<PeriodSelector current={data.period} />
			</div>
		</div>

		<div class="grid grid-cols-1 divide-x divide-border border-b border-border sm:grid-cols-3">
			<StatBox label="Total Clicks"    value={liveClicks}  color="green" />
			<StatBox label="Unique Visitors" value={data.analytics.unique_clicks} color="cyan"  />
			<StatBox label="Avg / Day"       value={avgPerDay} />
		</div>
	</div>

	<div class="tui-panel">
		<div class="tui-panel-header">▌ clicks over time</div>
		<div class="p-4">
			<LineChart data={data.analytics.over_time} />
		</div>
	</div>

	<AnalyticsBreakdowns {devices} {browsers} {referrers} />

	<div>
		<a
			href="/{data.slug}/analytics"
			class="font-mono text-xs text-muted-foreground hover:text-accent transition-colors"
		>
			← back to project analytics
		</a>
	</div>
</div>

<script lang="ts">
	import type { PageData } from './$types'
	import { LineChart, StatBox, PeriodSelector, AnalyticsBreakdowns } from '$lib'
	import { deviceItems, referrerItems, browserItems } from '$lib/analytics'

	let { data }: { data: PageData } = $props()

	const devices = $derived(deviceItems(data.analytics.devices))
	const referrers = $derived(referrerItems(data.analytics.referrers))
	const browsers = $derived(browserItems(data.analytics.browsers))

	const avgPerDay = $derived(
		data.analytics.over_time.length > 0
			? Math.round(
					data.analytics.total_clicks / Math.max(data.analytics.over_time.length, 1),
				)
			: 0,
	)
</script>

<div class="mx-auto max-w-5xl space-y-4">
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
			<PeriodSelector current={data.period} />
		</div>

		<div class="grid grid-cols-1 divide-x divide-border border-b border-border sm:grid-cols-3">
			<StatBox label="Total Clicks"    value={data.analytics.total_clicks}  color="green" />
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

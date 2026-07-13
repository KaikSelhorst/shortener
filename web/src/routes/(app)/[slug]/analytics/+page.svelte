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
		(data.analytics.over_time ?? []).length > 0
			? Math.round(liveClicks / Math.max((data.analytics.over_time ?? []).length, 1))
			: 0
	)

	const sse = useSSE(() => `/api/${data.slug}/stream`, () => { liveClicks++ })
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
<div class="grid grid-cols-2 sm:grid-cols-4 border-b border-border">
	<StatCard label="Total Clicks" value={liveClicks} color="success" />
	<StatCard label="Unique Visitors" value={data.analytics.unique_clicks} />
	<StatCard label="Links" value={(data.analytics.top_links ?? []).length} />
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

<!-- Top links -->
{#if data.analytics.top_links.length > 0}
	<div class="border-b border-border px-4 h-10 flex items-center">
		<span class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">Top Links</span>
	</div>
	<table class="w-full text-sm">
		<thead class="border-b border-border">
			<tr>
				<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Short code</th>
				<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">URL</th>
				<th class="px-4 py-1.5 text-right text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Clicks</th>
			</tr>
		</thead>
		<tbody>
			{#each data.analytics.top_links as link (link.short_code)}
				<tr class="hover:bg-secondary/40 transition-colors">
					<td class="px-4 py-1.5">
						<a href="/{data.slug}/{link.short_code}/analytics" class="text-sm text-foreground hover:underline underline-offset-4">
							{link.short_code}
						</a>
					</td>
					<td class="px-4 py-1.5">
						<a href={link.original_url} target="_blank" rel="noopener noreferrer"
							class="block max-w-xs truncate text-sm text-muted-foreground hover:text-foreground transition-colors"
							title={link.title ?? link.original_url}>
							{link.title ?? link.original_url}
						</a>
					</td>
					<td class="px-4 py-1.5 text-right text-sm tabular-nums">{link.total_clicks.toLocaleString()}</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/if}

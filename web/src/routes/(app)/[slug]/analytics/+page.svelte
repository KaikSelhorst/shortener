<script lang="ts">
	import { untrack } from 'svelte'
	import type { PageData } from './$types'
	import { LineChart, StatBox, PeriodSelector, AnalyticsBreakdowns } from '$lib'
	import { deviceItems, referrerItems, browserItems } from '$lib/analytics'
	import { useSSE } from '$lib/sse.svelte'

	let { data }: { data: PageData } = $props()

	const devices = $derived(deviceItems(data.analytics.devices))
	const referrers = $derived(referrerItems(data.analytics.referrers))
	const browsers = $derived(browserItems(data.analytics.browsers))

	let liveClicks = $state(untrack(() => data.analytics.total_clicks))

	$effect(() => {
		liveClicks = data.analytics.total_clicks
	})

	const avgPerDay = $derived(
		(data.analytics.over_time ?? []).length > 0
			? Math.round(liveClicks / Math.max((data.analytics.over_time ?? []).length, 1))
			: 0,
	)

	const sse = useSSE(() => `/api/${data.slug}/stream`, () => { liveClicks++ })
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
				<span>▌ analytics</span>
			</div>
			<div class="flex items-center gap-3">
				{#if sse.connected}
					<span class="font-mono text-[10px] uppercase tracking-wider text-tui-green">● live</span>
				{/if}
				<PeriodSelector current={data.period} />
			</div>
		</div>

		<div class="grid grid-cols-2 divide-x divide-border border-b border-border sm:grid-cols-4">
			<StatBox label="Total Clicks"    value={liveClicks} color="green" />
			<StatBox label="Unique Visitors" value={data.analytics.unique_clicks} color="cyan"  />
			<StatBox label="Links"           value={(data.analytics.top_links ?? []).length} />
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

	{#if data.analytics.top_links.length > 0}
		<div class="tui-panel">
			<div class="tui-panel-header">▌ top links</div>
			<table class="w-full">
				<thead class="border-b border-border">
					<tr>
						<th class="px-4 py-2.5 text-left font-normal tui-label">Short code</th>
						<th class="px-4 py-2.5 text-left font-normal tui-label">URL</th>
						<th class="px-4 py-2.5 text-right font-normal tui-label">Clicks</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					{#each data.analytics.top_links as link (link.short_code)}
						<tr class="transition-colors hover:bg-primary/5">
							<td class="px-4 py-2.5">
								<a
									href="/{data.slug}/{link.short_code}/analytics"
									class="font-mono text-xs text-accent hover:underline"
								>
									{link.short_code}
								</a>
							</td>
							<td class="px-4 py-2.5">
								<a
									href={link.original_url}
									target="_blank"
									rel="noopener noreferrer"
									class="block max-w-xs truncate font-mono text-xs text-muted-foreground hover:text-foreground transition-colors"
									title={link.title ?? link.original_url}
								>
									{link.title ?? link.original_url}
								</a>
							</td>
							<td class="px-4 py-2.5 text-right font-mono text-xs text-foreground">
								{link.total_clicks.toLocaleString()}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

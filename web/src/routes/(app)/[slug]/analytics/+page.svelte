<script lang="ts">
	import type { PageData } from './$types'
	import { BarChart, LineChart, Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '$lib'
	import { periods, deviceItems, referrerItems } from '$lib/analytics'

	let { data }: { data: PageData } = $props()

	const devices = $derived(deviceItems(data.analytics.devices))
	const referrers = $derived(referrerItems(data.analytics.referrers))
</script>

<!-- breadcrumb -->
<div class="flex items-center gap-2 text-sm">
	<a href="/dashboard" class="text-muted-foreground hover:text-foreground">Projects</a>
	<span class="text-muted-foreground">/</span>
	<a href="/{data.slug}" class="text-muted-foreground hover:text-foreground">{data.slug}</a>
	<span class="text-muted-foreground">/</span>
	<span class="font-semibold text-foreground">Analytics</span>
</div>

<!-- period selector -->
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
<div class="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-4">
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
	<div class="rounded-md bg-card p-4 shadow-sm">
		<p class="text-xs text-muted-foreground">Links</p>
		<p class="mt-1 text-2xl font-semibold text-foreground">
			{(data.analytics.top_links ?? []).length}
		</p>
	</div>
	<div class="rounded-md bg-card p-4 shadow-sm">
		<p class="text-xs text-muted-foreground">Avg. per Day</p>
		<p class="mt-1 text-2xl font-semibold text-foreground">
			{(data.analytics.over_time ?? []).length > 0
				? Math.round(
						data.analytics.total_clicks /
							Math.max((data.analytics.over_time ?? []).length, 1),
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
		<BarChart items={devices} />
	</div>
	<div class="rounded-md bg-card p-4 shadow-sm">
		<h2 class="mb-4 text-sm font-semibold text-foreground">Referrers</h2>
		<BarChart items={referrers} />
	</div>
</div>

<!-- top links -->
{#if data.analytics.top_links.length > 0}
	<div class="mt-4 rounded-md bg-card shadow-sm">
		<div class="px-4 pt-4">
			<h2 class="text-sm font-semibold text-foreground">Top links</h2>
		</div>
		<div class="mt-3">
			<Table>
				<TableHead>
					<TableRow>
						<TableHeader>Short code</TableHeader>
						<TableHeader>URL</TableHeader>
						<TableHeader class="text-right">Clicks</TableHeader>
					</TableRow>
				</TableHead>
				<TableBody>
					{#each data.analytics.top_links as link (link.short_code)}
						<TableRow>
							<TableCell>
								<a
									href="/{data.slug}/{link.short_code}/analytics"
									class="font-mono text-sm hover:underline"
								>
									{link.short_code}
								</a>
							</TableCell>
							<TableCell>
								<a
									href={link.original_url}
									target="_blank"
									rel="noopener noreferrer"
									class="block max-w-xs truncate text-muted-foreground hover:text-foreground"
									title={link.original_url}
								>
									{link.title ?? link.original_url}
								</a>
							</TableCell>
							<TableCell class="text-right font-mono text-sm">
								{link.total_clicks.toLocaleString()}
							</TableCell>
						</TableRow>
					{/each}
				</TableBody>
			</Table>
		</div>
	</div>
{/if}

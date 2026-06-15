<script lang="ts">
	import type { PageData } from './$types'
	import {
		Badge,
		Table,
		TableHead,
		TableBody,
		TableRow,
		TableHeader,
		TableCell,
	} from '$lib'

	let { data }: { data: PageData } = $props()

	function extractShortCode(d: { event: string; payload: Record<string, unknown> }): string | null {
		if (d.event === 'link.clicked') return (d.payload.short_code as string) ?? null
		const link = d.payload.link as Record<string, unknown> | undefined
		return (link?.short_code as string) ?? null
	}
</script>

<div class="mx-auto max-w-6xl">
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
				<a
					href="/{data.slug}/webhooks"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
				>
					webhooks
				</a>
				<span class="text-muted-foreground">/</span>
				<span class="min-w-0 truncate max-w-[40%] md:max-w-none" title={data.webhook.url}>▌ {data.webhook.url}</span>
			</div>
		</div>

		<Table>
			<TableHead>
				<TableRow>
					<TableHeader>Event</TableHeader>
					<TableHeader>Link</TableHeader>
					<TableHeader>Status</TableHeader>
					<TableHeader>HTTP</TableHeader>
					<TableHeader>Attempts</TableHeader>
					<TableHeader>Date</TableHeader>
				</TableRow>
			</TableHead>
			<TableBody>
				{#each data.deliveries as d (d.id)}
					<TableRow>
						<TableCell>
							<span class="font-mono text-xs">{d.event}</span>
						</TableCell>
						<TableCell>
							{#if extractShortCode(d)}
								<a
									href="/{data.slug}/{extractShortCode(d)}"
									class="font-mono text-xs text-muted-foreground hover:text-foreground transition-colors"
								>/{extractShortCode(d)}</a>
							{:else}
								<span class="text-muted-foreground">—</span>
							{/if}
						</TableCell>
						<TableCell>
							{#if d.status === 'delivered'}
								<Badge variant="success">delivered</Badge>
							{:else if d.status === 'failed'}
								<Badge variant="error">failed</Badge>
							{:else if d.status === 'processing'}
								<Badge variant="solid">processing</Badge>
							{:else}
								<Badge>pending</Badge>
							{/if}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{d.response_status ?? '—'}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{d.attempts}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{new Date(d.created_at).toLocaleString()}
						</TableCell>
					</TableRow>
				{:else}
					<TableRow>
						<TableCell colspan={6} class="py-14 text-center text-muted-foreground">
							-- no deliveries yet --
						</TableCell>
					</TableRow>
				{/each}
			</TableBody>
		</Table>

		{#if data.page > 1 || data.hasMore}
			<div class="flex items-center justify-between border-t border-border px-4 py-2">
				{#if data.page > 1}
					<a
						href="?page={data.page - 1}"
						class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
					>
						← prev
					</a>
				{:else}
					<span class="font-mono text-[10px] uppercase tracking-wider opacity-30">← prev</span>
				{/if}
				<span class="font-mono text-[10px] text-muted-foreground">page {data.page}</span>
				{#if data.hasMore}
					<a
						href="?page={data.page + 1}"
						class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
					>
						next →
					</a>
				{:else}
					<span class="font-mono text-[10px] uppercase tracking-wider opacity-30">next →</span>
				{/if}
			</div>
		{/if}
	</div>
</div>

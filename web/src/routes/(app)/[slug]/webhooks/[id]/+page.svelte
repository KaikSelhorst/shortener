<script lang="ts">
	import type { PageData } from './$types'
	import { formatDateTime } from '$lib/format'
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

<div class="sticky top-0 z-10 bg-background border-b border-border px-4 h-11 flex items-center">
	<div class="flex items-center gap-1.5 text-sm min-w-0">
		<a href="/dashboard" class="text-muted-foreground hover:text-foreground transition-colors shrink-0">Projects</a>
		<span class="text-border shrink-0">/</span>
		<a href="/{data.slug}" class="text-muted-foreground hover:text-foreground transition-colors shrink-0">{data.slug}</a>
		<span class="text-border shrink-0">/</span>
		<a href="/{data.slug}/webhooks" class="text-muted-foreground hover:text-foreground transition-colors shrink-0">Webhooks</a>
		<span class="text-border shrink-0">/</span>
		<span class="font-medium text-foreground truncate" title={data.webhook.url}>{data.webhook.url}</span>
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
				<TableCell>{d.event}</TableCell>
				<TableCell>
					{@const code = extractShortCode(d)}
					{#if code}
						<a href="/{data.slug}/{code}/analytics" class="text-muted-foreground hover:text-foreground transition-colors">/{code}</a>
					{:else}
						<span class="text-muted-foreground">—</span>
					{/if}
				</TableCell>
				<TableCell>
					{#if d.status === 'delivered'}<Badge variant="success">Delivered</Badge>
					{:else if d.status === 'failed'}<Badge variant="error">Failed</Badge>
					{:else if d.status === 'processing'}<Badge variant="solid">Processing</Badge>
					{:else}<Badge>Pending</Badge>
					{/if}
				</TableCell>
				<TableCell class="text-muted-foreground">{d.response_status ?? '—'}</TableCell>
				<TableCell class="text-muted-foreground">{d.attempts}</TableCell>
				<TableCell class="text-muted-foreground">{formatDateTime(d.created_at)}</TableCell>
			</TableRow>
		{:else}
			<TableRow>
				<TableCell colspan={6} class="py-16 text-center text-muted-foreground">No deliveries yet.</TableCell>
			</TableRow>
		{/each}
	</TableBody>
</Table>

{#if data.page > 1 || data.hasMore}
	<div class="flex items-center justify-between border-t border-border px-4 py-2.5">
		{#if data.page > 1}
			<a href="?page={data.page - 1}" class="text-sm text-muted-foreground hover:text-foreground transition-colors">← Prev</a>
		{:else}
			<span class="text-sm text-muted-foreground/30">← Prev</span>
		{/if}
		<span class="text-sm text-muted-foreground">Page {data.page}</span>
		{#if data.hasMore}
			<a href="?page={data.page + 1}" class="text-sm text-muted-foreground hover:text-foreground transition-colors">Next →</a>
		{:else}
			<span class="text-sm text-muted-foreground/30">Next →</span>
		{/if}
	</div>
{/if}

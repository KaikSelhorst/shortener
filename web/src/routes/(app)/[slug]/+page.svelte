<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import type { Link } from '$lib/types'
	import { enhance } from '$app/forms'
	import { useSSE } from '$lib/state/sse.svelte'
	import { useQrcode } from '$lib/state/qrcode.svelte'
	import { formatDate, toDatetimeLocal, parseUtm } from '$lib/format'
	import {
		Badge,
		Button,
		Dialog,
		Table,
		TableHead,
		TableBody,
		TableRow,
		TableHeader,
		TableCell,
	} from '$lib/components/ui'
	import { QRCode, LinkFormFields } from '$lib/components/links'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	let createOpen = $state(false)
	let createError = $state<string | null>(null)
	let createUtmExpanded = $state(false)

	$effect(() => {
		if (!createOpen) { createError = null; createUtmExpanded = false }
	})

	const handleCreate: SubmitFunction = () => {
		createError = null
		return async ({ result, update }) => {
			if (result.type === 'failure') {
				const d = result.data as { error?: string }
				createError = d?.error ?? 'Failed to create link'
			} else {
				await update()
				createOpen = false
			}
		}
	}

	let editLink = $state<Link | null>(null)
	let editOpen = $state(false)
	let editError = $state<string | null>(null)
	let editUtmExpanded = $state(false)

	$effect(() => {
		if (!editOpen) { editLink = null; editError = null; editUtmExpanded = false }
	})

	const handleEdit: SubmitFunction = () => {
		editError = null
		return async ({ result, update }) => {
			if (result.type === 'failure') {
				const d = result.data as { error?: string }
				editError = d?.error ?? 'Failed to update link'
			} else {
				await update()
				editOpen = false
			}
		}
	}

	let pendingCode = $state<string | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingCode = null
	})

	const qr = useQrcode()

	let sessionClicks = $state<Record<string, number>>({})

	const now = Date.now()

	const sse = useSSE(() => `/api/${data.slug}/stream`, (e) => {
		const evt = JSON.parse(e.data) as { short_code: string }
		sessionClicks[evt.short_code] = (sessionClicks[evt.short_code] ?? 0) + 1
	})
</script>

<div class="sticky top-0 z-10 bg-background border-b border-border px-4 h-11 flex items-center justify-between">
	<div class="flex items-center gap-1.5 text-sm min-w-0">
		{#if sse.connected}
			<span class="flex items-center gap-1.5 text-xs text-success shrink-0">
				<span class="w-1.5 h-1.5 rounded-full bg-success"></span>Live
			</span>
		{/if}
	</div>
	<Button size="sm" onclick={() => (createOpen = true)}>+ Link</Button>
</div>

{#if form?.error}
	<div class="border-b border-border px-4 py-2">
		<p class="text-xs text-destructive">{form.error}</p>
	</div>
{/if}

<Table>
	<TableHead>
		<TableRow>
			<TableHeader>Original URL</TableHeader>
			<TableHeader>Short URL</TableHeader>
			<TableHeader>Status</TableHeader>
			<TableHeader>Clicks</TableHeader>
			<TableHeader>Created</TableHeader>
			<TableHeader>Expires</TableHeader>
			<TableHeader class="text-right"></TableHeader>
		</TableRow>
	</TableHead>
	<TableBody>
		{#each data.links.data as link (link.id)}
			{@const expired = link.expires_at != null && new Date(link.expires_at).getTime() < now}
			{@const totalClicks = link.total_clicks + (sessionClicks[link.short_code] ?? 0)}
			{@const maxed = link.max_clicks != null && totalClicks >= link.max_clicks}
			<TableRow>
				<TableCell>
					<a href={link.original_url} target="_blank" rel="noopener noreferrer"
						class="block max-w-xs truncate text-muted-foreground hover:text-foreground transition-colors" title={link.original_url}>
						{link.original_url}
					</a>
				</TableCell>
				<TableCell>
					<a href={link.short_url} target="_blank" rel="noopener noreferrer" class="text-foreground hover:underline underline-offset-4">
						{link.short_url}
					</a>
				</TableCell>
				<TableCell>
					{#if expired || maxed}
						<Badge variant="error">{maxed ? 'maxed' : 'expired'}</Badge>
					{:else}
						<Badge variant="success">active</Badge>
					{/if}
				</TableCell>
				<TableCell>
					<span class="tabular-nums {sessionClicks[link.short_code] ? 'text-success' : 'text-muted-foreground'}">
						{totalClicks.toLocaleString()}{link.max_clicks != null ? ` / ${link.max_clicks.toLocaleString()}` : ''}
					</span>
				</TableCell>
				<TableCell class="text-muted-foreground">{formatDate(link.created_at)}</TableCell>
				<TableCell class="text-muted-foreground">{formatDate(link.expires_at)}</TableCell>
				<TableCell class="text-right">
					<div class="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
						<Button variant="ghost" size="sm" onclick={() => {
							const utm = parseUtm(link.original_url)
							editUtmExpanded = !!(utm.source || utm.medium || utm.campaign || utm.term || utm.content)
							editLink = link; editOpen = true
						}}>Edit</Button>
						<Button variant="ghost" size="sm" onclick={() => qr.show(link.short_url, link.short_code)}>QR</Button>
						<Button variant="ghost" size="sm" href="/{data.slug}/{link.short_code}/analytics">Analytics</Button>
						<Button variant="ghost-destructive" size="sm" onclick={() => { pendingCode = link.short_code; confirmOpen = true }}>Delete</Button>
					</div>
				</TableCell>
			</TableRow>
		{:else}
			<TableRow>
				<TableCell colspan={7} class="py-16 text-center text-muted-foreground">No links yet.</TableCell>
			</TableRow>
		{/each}
	</TableBody>
</Table>

{#if data.links.prev_cursor || data.links.next_cursor}
	<div class="flex items-center justify-between border-t border-border px-4 py-2.5">
		<div>{#if data.links.prev_cursor}<Button variant="outline" size="sm" href="?cursor={data.links.prev_cursor}">← Prev</Button>{/if}</div>
		<div>{#if data.links.next_cursor}<Button variant="outline" size="sm" href="?cursor={data.links.next_cursor}">Next →</Button>{/if}</div>
	</div>
{/if}

<Dialog bind:open={createOpen} title="New link" size="md">
	{#snippet children()}
		<form
			id="create-link-form"
			method="POST"
			action="?/create"
			use:enhance={handleCreate}
			class="flex flex-col gap-4"
		>
			<LinkFormFields mode="create" bind:utmExpanded={createUtmExpanded} />
			{#if createError}
				<p class="text-sm text-destructive">{createError}</p>
			{/if}
		</form>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (createOpen = false)}>Cancel</Button>
		<Button type="submit" form="create-link-form">Create</Button>
	{/snippet}
</Dialog>

<Dialog bind:open={editOpen} title="Edit link" size="md">
	{#snippet children()}
		{#if editLink}
			{@const utm = parseUtm(editLink.original_url)}
			<form
				id="edit-link-form"
				method="POST"
				action="?/update"
				use:enhance={handleEdit}
				class="flex flex-col gap-4"
			>
				<input type="hidden" name="code" value={editLink.short_code} />
				<LinkFormFields
					mode="edit"
					bind:utmExpanded={editUtmExpanded}
					url={utm.baseUrl}
					title={editLink.title ?? ''}
					expiresAt={toDatetimeLocal(editLink.expires_at)}
					maxClicks={editLink.max_clicks ?? ''}
					utm={{ source: utm.source, medium: utm.medium, campaign: utm.campaign, term: utm.term, content: utm.content }}
				/>
				{#if editError}
					<p class="text-sm text-destructive">{editError}</p>
				{/if}
			</form>
		{/if}
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (editOpen = false)}>Cancel</Button>
		<Button type="submit" form="edit-link-form">Save</Button>
	{/snippet}
</Dialog>

<Dialog bind:open={qr.open} title="QR Code">
	{#snippet children()}
		{#if qr.url}
			<div class="flex flex-col items-center gap-4">
				<div bind:this={qr.container}>
					<QRCode data={qr.url} moduleSize={6} />
				</div>
				<p class="text-xs text-muted-foreground break-all text-center">{qr.url}</p>
			</div>
		{/if}
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={qr.downloadSVG}>↓ SVG</Button>
		<Button variant="outline" onclick={qr.downloadPNG}>↓ PNG</Button>
		<Button onclick={() => (qr.open = false)}>Close</Button>
	{/snippet}
</Dialog>

<Dialog
	bind:open={confirmOpen}
	title="Delete link"
	description="This will permanently delete the shortened link. This action cannot be undone."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>Cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="code" value={pendingCode} />
			<Button type="submit" variant="destructive">Delete</Button>
		</form>
	{/snippet}
</Dialog>

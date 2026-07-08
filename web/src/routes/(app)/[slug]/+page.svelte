<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import type { Link } from '$lib/types'
	import { enhance } from '$app/forms'
	import { useSSE } from '$lib/sse.svelte'
	import { useQrcode } from '$lib/qrcode.svelte'
	import { formatDate, toDatetimeLocal, parseUtm } from '$lib/format'
	import {
		Badge,
		Button,
		Input,
		Dialog,
		QRCode,
		Table,
		TableHead,
		TableBody,
		TableRow,
		TableHeader,
		TableCell,
	} from '$lib'

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
		<a href="/dashboard" class="text-muted-foreground hover:text-foreground transition-colors shrink-0">Projects</a>
		<span class="text-border shrink-0">/</span>
		<span class="font-medium text-foreground truncate">{data.slug}</span>
		{#if sse.connected}
			<span class="flex items-center gap-1.5 text-xs text-success ml-2 shrink-0">
				<span class="w-1.5 h-1.5 rounded-full bg-success"></span>Live
			</span>
		{/if}
	</div>
	<div class="flex items-center gap-2 shrink-0">
		<a href="/{data.slug}/analytics" class="text-sm text-muted-foreground hover:text-foreground transition-colors">Analytics</a>
		<a href="/{data.slug}/webhooks" class="text-sm text-muted-foreground hover:text-foreground transition-colors">Webhooks</a>
		<Button size="sm" onclick={() => (createOpen = true)}>+ Link</Button>
	</div>
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
			<Input name="url" type="url" label="URL" placeholder="https://example.com" required />
			<Input name="title" type="text" label="Title" placeholder="optional title" />
			<Input
				name="custom_code"
				type="text"
				label="Short code"
				placeholder="my-link (leave blank to auto-generate)"
				pattern="[a-zA-Z0-9_\-]+"
				minlength={3}
				maxlength={50}
			/>
			<Input name="expires_at" type="datetime-local" label="Expires at" />
			<Input name="max_clicks" type="number" label="Max clicks" placeholder="unlimited" min={1} />
			<button
				type="button"
				onclick={() => (createUtmExpanded = !createUtmExpanded)}
				class="text-sm text-muted-foreground hover:text-foreground transition-colors text-left"
			>
				{createUtmExpanded ? '▾' : '▸'} UTM parameters
			</button>
			{#if createUtmExpanded}
				<div class="flex flex-col gap-4 border-l border-border pl-3">
					<Input name="utm_source" type="text" label="Source" placeholder="newsletter, twitter, google" />
					<Input name="utm_medium" type="text" label="Medium" placeholder="email, social, cpc" />
					<Input name="utm_campaign" type="text" label="Campaign" placeholder="spring_sale" />
					<Input name="utm_term" type="text" label="Term" placeholder="keywords (optional)" />
					<Input name="utm_content" type="text" label="Content" placeholder="variant_a (optional)" />
				</div>
			{/if}
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
				<Input name="url" type="url" label="URL" required value={utm.baseUrl} />
				<Input name="title" type="text" label="Title" placeholder="optional title" value={editLink.title ?? ''} />
				<Input name="expires_at" type="datetime-local" label="Expires at" value={toDatetimeLocal(editLink.expires_at)} />
				<Input name="max_clicks" type="number" label="Max clicks" placeholder="unlimited" min={1} value={editLink.max_clicks ?? ''} />
				<button
					type="button"
					onclick={() => (editUtmExpanded = !editUtmExpanded)}
					class="text-sm text-muted-foreground hover:text-foreground transition-colors text-left"
				>
					{editUtmExpanded ? '▾' : '▸'} UTM parameters
				</button>
				{#if editUtmExpanded}
					<div class="flex flex-col gap-4 border-l border-border pl-3">
						<Input name="utm_source" type="text" label="Source" placeholder="newsletter, twitter, google" value={utm.source} />
						<Input name="utm_medium" type="text" label="Medium" placeholder="email, social, cpc" value={utm.medium} />
						<Input name="utm_campaign" type="text" label="Campaign" placeholder="spring_sale" value={utm.campaign} />
						<Input name="utm_term" type="text" label="Term" placeholder="keywords (optional)" value={utm.term} />
						<Input name="utm_content" type="text" label="Content" placeholder="variant_a (optional)" value={utm.content} />
					</div>
				{/if}
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

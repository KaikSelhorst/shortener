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

<div class="mx-auto max-w-7xl">
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
				<span>▌ {data.slug}</span>
			</div>
			<div class="flex items-center gap-3">
				{#if sse.connected}
					<span class="font-mono text-[10px] uppercase tracking-wider text-tui-green">● live</span>
				{/if}
				<a
					href="/{data.slug}/webhooks"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-accent transition-colors"
				>
					webhooks
				</a>
				<a
					href="/{data.slug}/analytics"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-accent transition-colors"
				>
					analytics
				</a>
				<Button size="sm" onclick={() => (createOpen = true)}>+ link</Button>
			</div>
		</div>

		{#if form?.error}
			<div class="border-b border-border px-4 py-2">
				<p class="font-mono text-xs text-destructive">{form.error}</p>
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
							<a
								href={link.original_url}
								target="_blank"
								rel="noopener noreferrer"
								class="block max-w-xs truncate text-muted-foreground hover:text-foreground transition-colors"
								title={link.original_url}
							>
								{link.original_url}
							</a>
						</TableCell>
						<TableCell>
							<a
								href={link.short_url}
								target="_blank"
								rel="noopener noreferrer"
								class="text-accent hover:underline"
							>
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
							<span class="font-mono text-xs tabular-nums {sessionClicks[link.short_code] ? 'text-tui-green' : 'text-muted-foreground'}">
								{totalClicks.toLocaleString()}{link.max_clicks != null ? ` / ${link.max_clicks.toLocaleString()}` : ''}
							</span>
						</TableCell>
						<TableCell class="text-muted-foreground">
							{formatDate(link.created_at)}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{formatDate(link.expires_at)}
						</TableCell>
						<TableCell class="text-right">
							<div class="flex items-center justify-end gap-2">
								<Button
									variant="ghost"
									size="sm"
									onclick={() => {
										const utm = parseUtm(link.original_url)
										editUtmExpanded = !!(utm.source || utm.medium || utm.campaign || utm.term || utm.content)
										editLink = link
										editOpen = true
									}}
								>
									edit
								</Button>
								<Button
									variant="ghost"
									size="sm"
									onclick={() => qr.show(link.short_url, link.short_code)}
								>
									qr
								</Button>
								<Button
									variant="ghost"
									size="sm"
									href="/{data.slug}/{link.short_code}/analytics"
								>
									analytics
								</Button>
								<Button
									variant="ghost-destructive"
									size="sm"
									onclick={() => {
										pendingCode = link.short_code
										confirmOpen = true
									}}
								>
									delete
								</Button>
							</div>
						</TableCell>
					</TableRow>
				{:else}
					<TableRow>
						<TableCell colspan={7} class="py-14 text-center text-muted-foreground">
							-- no links yet --
						</TableCell>
					</TableRow>
				{/each}
			</TableBody>
		</Table>

		{#if data.links.prev_cursor || data.links.next_cursor}
			<div class="flex items-center justify-between border-t border-border px-4 py-2.5">
				<div>
					{#if data.links.prev_cursor}
						<Button variant="outline" size="sm" href="?cursor={data.links.prev_cursor}">
							← prev
						</Button>
					{/if}
				</div>
				<div>
					{#if data.links.next_cursor}
						<Button variant="outline" size="sm" href="?cursor={data.links.next_cursor}">
							next →
						</Button>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>

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
				class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors text-left"
			>
				{createUtmExpanded ? '▾' : '▸'} utm parameters
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
				<p class="font-mono text-xs text-destructive">{createError}</p>
			{/if}
		</form>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (createOpen = false)}>cancel</Button>
		<Button type="submit" form="create-link-form">create</Button>
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
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors text-left"
				>
					{editUtmExpanded ? '▾' : '▸'} utm parameters
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
					<p class="font-mono text-xs text-destructive">{editError}</p>
				{/if}
			</form>
		{/if}
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (editOpen = false)}>cancel</Button>
		<Button type="submit" form="edit-link-form">save</Button>
	{/snippet}
</Dialog>

<Dialog bind:open={qr.open} title="QR Code">
	{#snippet children()}
		{#if qr.url}
			<div class="flex flex-col items-center gap-4">
				<div bind:this={qr.container}>
					<QRCode data={qr.url} moduleSize={6} />
				</div>
				<p class="font-mono text-[10px] text-muted-foreground break-all text-center">{qr.url}</p>
			</div>
		{/if}
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={qr.downloadSVG}>↓ svg</Button>
		<Button variant="outline" onclick={qr.downloadPNG}>↓ png</Button>
		<Button onclick={() => (qr.open = false)}>close</Button>
	{/snippet}
</Dialog>

<Dialog
	bind:open={confirmOpen}
	title="Delete link"
	description="This will permanently delete the shortened link. This action cannot be undone."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="code" value={pendingCode} />
			<Button type="submit" variant="destructive">confirm delete</Button>
		</form>
	{/snippet}
</Dialog>

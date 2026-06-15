<script lang="ts">
	import { onMount } from 'svelte'
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import { enhance } from '$app/forms'
	import {
		Badge,
		Button,
		Input,
		Dialog,
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

	$effect(() => {
		if (!createOpen) createError = null
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

	let pendingCode = $state<string | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingCode = null
	})

	let sessionClicks = $state<Record<string, number>>({})
	let connected = $state(false)

	const linkClicks = $derived(
		Object.fromEntries(
			data.links.data.map((l) => [l.short_code, l.total_clicks + (sessionClicks[l.short_code] ?? 0)]),
		),
	)

	onMount(() => {
		const es = new EventSource(`/api/${data.slug}/stream`)
		es.onopen = () => (connected = true)
		es.onmessage = (e) => {
			const evt = JSON.parse(e.data) as { short_code: string }
			sessionClicks[evt.short_code] = (sessionClicks[evt.short_code] ?? 0) + 1
		}
		es.onerror = () => (connected = false)
		return () => es.close()
	})
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
				<span>▌ {data.slug}</span>
			</div>
			<div class="flex items-center gap-3">
				{#if connected}
					<span class="font-mono text-[10px] uppercase tracking-wider text-tui-green">● live</span>
				{/if}
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
					{@const expired = link.expires_at != null && new Date(link.expires_at) < new Date()}
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
							{#if expired}
								<Badge variant="error">expired</Badge>
							{:else}
								<Badge variant="success">active</Badge>
							{/if}
						</TableCell>
						<TableCell>
							<span class="font-mono text-xs tabular-nums {sessionClicks[link.short_code] ? 'text-tui-green' : 'text-muted-foreground'}">
								{linkClicks[link.short_code].toLocaleString()}
							</span>
						</TableCell>
						<TableCell class="text-muted-foreground">
							{new Date(link.created_at).toLocaleDateString()}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{link.expires_at ? new Date(link.expires_at).toLocaleDateString() : '—'}
						</TableCell>
						<TableCell class="text-right">
							<div class="flex items-center justify-end gap-2">
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
			<Input name="expires_at" type="datetime-local" label="Expires at" />
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

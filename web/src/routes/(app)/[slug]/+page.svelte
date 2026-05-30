<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import { enhance } from '$app/forms'
	import {
		Button,
		Input,
		Dialog,
		Table,
		TableHead,
		TableBody,
		TableRow,
		TableHeader,
		TableCell
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
</script>

<div class="flex items-center justify-between">
	<div class="flex items-center gap-2">
		<a href="/dashboard" class="text-sm text-muted-foreground hover:text-foreground">Projects</a>
		<span class="text-muted-foreground">/</span>
		<h1 class="text-sm font-semibold text-foreground">{data.slug}</h1>
	</div>
	<div class="flex items-center gap-2">
		<Button variant="outline" size="sm" href="/{data.slug}/analytics">Analytics</Button>
		<Button size="sm" onclick={() => (createOpen = true)}>New link</Button>
	</div>
</div>

{#if form?.error}
	<p class="mt-4 rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">{form.error}</p>
{/if}

<div class="mt-6">
	<Table>
		<TableHead>
			<TableRow>
				<TableHeader>Original URL</TableHeader>
				<TableHeader>Short URL</TableHeader>
				<TableHeader>Created</TableHeader>
				<TableHeader>Expires at</TableHeader>
				<TableHeader class="text-right"></TableHeader>
				<TableHeader class="text-right"></TableHeader>
			</TableRow>
		</TableHead>
		<TableBody>
			{#each data.links.data as link (link.id)}
				<TableRow>
					<TableCell>
						<a
							href={link.original_url}
							target="_blank"
							rel="noopener noreferrer"
							class="block max-w-xs truncate text-muted-foreground hover:text-foreground"
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
							class="hover:underline"
						>
							{link.short_url}
						</a>
					</TableCell>
					<TableCell class="text-muted-foreground">
						{new Date(link.created_at).toLocaleDateString()}
					</TableCell>
					<TableCell class="text-muted-foreground">
						{link.expires_at ? new Date(link.expires_at).toLocaleDateString() : '—'}
					</TableCell>
					<TableCell class="text-right">
						<Button
							variant="ghost"
							size="sm"
							href="/{data.slug}/{link.short_code}/analytics"
						>
							Analytics
						</Button>
					</TableCell>
					<TableCell class="text-right">
						<Button
							variant="ghost-destructive"
							size="sm"
							onclick={() => { pendingCode = link.short_code; confirmOpen = true }}
						>
							Delete
						</Button>
					</TableCell>
				</TableRow>
			{:else}
				<TableRow>
					<TableCell colspan={6} class="py-10 text-center text-muted-foreground">
						No links yet. Click "New link" to add one.
					</TableCell>
				</TableRow>
			{/each}
		</TableBody>
	</Table>
</div>

{#if data.links.prev_cursor || data.links.next_cursor}
	<div class="mt-4 flex items-center justify-between">
		{#if data.links.prev_cursor}
			<Button variant="outline" size="sm" href="?cursor={data.links.prev_cursor}">
				← Previous
			</Button>
		{:else}
			<span></span>
		{/if}

		{#if data.links.next_cursor}
			<Button variant="outline" size="sm" href="?cursor={data.links.next_cursor}">
				Next →
			</Button>
		{/if}
	</div>
{/if}

<Dialog bind:open={createOpen} title="New link" size="md">
	{#snippet children()}
		<form
			id="create-link-form"
			method="POST"
			action="?/create"
			use:enhance={handleCreate}
			class="flex flex-col gap-3"
		>
			<Input name="url" type="url" label="URL" placeholder="https://example.com" required />
			<Input name="title" type="text" label="Title" placeholder="Title (optional)" />
			<Input name="expires_at" type="datetime-local" label="Expires at" />
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

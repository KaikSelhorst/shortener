<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import type { CreateWebhookResponse } from '$lib/types'
	import { enhance } from '$app/forms'
	import { formatDate } from '$lib/format'
	import {
		Badge,
		Button,
		CopyInput,
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

	// Create dialog
	let createOpen = $state(false)
	let createError = $state<string | null>(null)
	let createdWebhook = $state<CreateWebhookResponse | null>(null)
	let secretOpen = $state(false)

	$effect(() => {
		if (!createOpen) createError = null
	})

	const handleCreate: SubmitFunction = () => {
		createError = null
		return async ({ result, update }) => {
			if (result.type === 'failure') {
				createError = (result.data as { createError?: string })?.createError ?? 'Failed to create webhook'
			} else if (result.type === 'success') {
				createdWebhook = (result.data as { created?: CreateWebhookResponse })?.created ?? null
				createOpen = false
				secretOpen = true
				await update()
			}
		}
	}

	// Delete dialog
	let pendingDeleteId = $state<string | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingDeleteId = null
	})

	const ALL_EVENTS = ['link.clicked', 'link.created', 'link.updated', 'link.deleted']

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
				<a
					href="/{data.slug}"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
				>
					{data.slug}
				</a>
				<span class="text-muted-foreground">/</span>
				<span>▌ webhooks</span>
			</div>
			<div class="flex items-center gap-3">
				<a
					href="/{data.slug}/analytics"
					class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground hover:text-accent transition-colors"
				>
					analytics
				</a>
				<Button size="sm" onclick={() => (createOpen = true)}>+ webhook</Button>
			</div>
		</div>

		<Table>
			<TableHead>
				<TableRow>
					<TableHeader>URL</TableHeader>
					<TableHeader>Events</TableHeader>
					<TableHeader>Status</TableHeader>
					<TableHeader>Created</TableHeader>
					<TableHeader class="text-right"></TableHeader>
				</TableRow>
			</TableHead>
			<TableBody>
				{#each data.webhooks as wh (wh.id)}
					<TableRow>
						<TableCell>
							<span class="block max-w-xs truncate font-mono text-xs text-muted-foreground" title={wh.url}>
								{wh.url}
							</span>
						</TableCell>
						<TableCell>
							<div class="flex flex-wrap gap-1">
								{#each wh.events as evt}
									<Badge variant="solid">{evt}</Badge>
								{/each}
							</div>
						</TableCell>
						<TableCell>
							{#if wh.enabled}
								<Badge variant="success">active</Badge>
							{:else}
								<Badge variant="error">disabled</Badge>
							{/if}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{formatDate(wh.created_at)}
						</TableCell>
						<TableCell class="text-right">
							<div class="flex items-center justify-end gap-2">
								<Button variant="ghost" size="sm" href="/{data.slug}/webhooks/{wh.id}">
									deliveries
								</Button>
								<Button
									variant="ghost-destructive"
									size="sm"
									aria-label="Delete webhook {wh.url}"
									onclick={() => {
										pendingDeleteId = wh.id
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
						<TableCell colspan={5} class="py-14 text-center text-muted-foreground">
							-- no webhooks yet --
						</TableCell>
					</TableRow>
				{/each}
			</TableBody>
		</Table>
	</div>
</div>

<!-- Create webhook dialog -->
<Dialog bind:open={createOpen} title="New webhook" size="md">
	{#snippet children()}
		<form
			id="create-webhook-form"
			method="POST"
			action="?/create"
			use:enhance={handleCreate}
			class="flex flex-col gap-4"
		>
			<Input name="url" type="url" label="Endpoint URL" placeholder="https://your-app.com/webhook" required />

			<fieldset class="flex flex-col gap-2 border-0 p-0 m-0">
				<legend class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground mb-2">Events</legend>
				<div class="grid grid-cols-2 gap-2">
					{#each ALL_EVENTS as evt}
						<label class="flex items-center gap-2 cursor-pointer">
							<input
								type="checkbox"
								name="events"
								value={evt}
							/>
							<span class="font-mono text-xs text-foreground">{evt}</span>
						</label>
					{/each}
				</div>
			</fieldset>

			{#if createError}
				<p class="font-mono text-xs text-destructive">{createError}</p>
			{/if}
		</form>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (createOpen = false)}>cancel</Button>
		<Button type="submit" form="create-webhook-form">create</Button>
	{/snippet}
</Dialog>

<!-- Secret reveal dialog -->
<Dialog bind:open={secretOpen} title="Save your webhook secret">
	{#snippet children()}
		<div class="flex flex-col gap-3">
			<p class="font-mono text-xs text-muted-foreground">
				This secret will not be shown again. Use it to verify the
				<span class="text-foreground">X-Webhook-Signature</span> header on incoming requests.
			</p>
			<CopyInput value={createdWebhook?.secret ?? ''} />
		</div>
	{/snippet}
	{#snippet footer()}
		<Button onclick={() => (secretOpen = false)}>I saved it</Button>
	{/snippet}
</Dialog>

<!-- Delete confirmation dialog -->
<Dialog
	bind:open={confirmOpen}
	title="Delete webhook"
	description="This will permanently delete the webhook and all delivery history."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="id" value={pendingDeleteId} />
			<Button type="submit" variant="destructive">confirm delete</Button>
		</form>
	{/snippet}
</Dialog>

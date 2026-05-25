<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import type { CreateApiKeyResponse } from '$lib/types'
	import { enhance } from '$app/forms'
	import { Badge, Button, Input, Dialog, Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '$lib'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	const ALL_SCOPES = [
		{ value: '*', label: 'Full access' },
		{ value: 'links:create', label: 'Links — Create' },
		{ value: 'links:read', label: 'Links — Read' },
		{ value: 'links:update', label: 'Links — Update' },
		{ value: 'links:delete', label: 'Links — Delete' },
		{ value: 'projects:create', label: 'Projects — Create' },
		{ value: 'projects:read', label: 'Projects — Read' },
		{ value: 'projects:update', label: 'Projects — Update' },
		{ value: 'projects:delete', label: 'Projects — Delete' },
	]

	let createOpen = $state(false)
	let createError = $state<string | null>(null)
	let createdKey = $state<CreateApiKeyResponse | null>(null)
	let tokenCopied = $state(false)

	$effect(() => {
		if (!createOpen) {
			createError = null
			tokenCopied = false
		}
	})

	const handleCreate: SubmitFunction = () => {
		createError = null
		return async ({ result, update }) => {
			if (result.type === 'failure') {
				const d = result.data as { createError?: string }
				createError = d?.createError ?? 'Failed to create key'
			} else if (result.type === 'success') {
				const d = result.data as { created?: CreateApiKeyResponse }
				if (d?.created) {
					createdKey = d.created
					createOpen = false
				}
				await update()
			}
		}
	}

	async function copyToken() {
		if (!createdKey) return
		await navigator.clipboard.writeText(createdKey.token)
		tokenCopied = true
	}

	let pendingDeleteId = $state<number | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingDeleteId = null
	})

</script>

<div class="flex items-center justify-between">
	<h1 class="text-lg font-semibold text-foreground">Settings</h1>
</div>

<div class="mt-8">
	<div class="flex items-center justify-between">
		<div>
			<h2 class="text-sm font-semibold text-foreground">API Keys</h2>
			<p class="mt-0.5 text-sm text-muted-foreground">Manage programmatic access to your account.</p>
		</div>
		<Button size="sm" onclick={() => (createOpen = true)}>New API Key</Button>
	</div>

	{#if form?.deleteError}
		<p class="mt-4 rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">{form.deleteError}</p>
	{/if}

	<div class="mt-4">
		<Table>
			<TableHead>
				<TableRow>
					<TableHeader>Name</TableHeader>
					<TableHeader>Prefix</TableHeader>
					<TableHeader>Scopes</TableHeader>
					<TableHeader>Project</TableHeader>
					<TableHeader>Last used</TableHeader>
					<TableHeader class="text-right"></TableHeader>
				</TableRow>
			</TableHead>
			<TableBody>
				{#each data.apiKeys as key (key.id)}
					<TableRow>
						<TableCell class="font-medium">{key.name}</TableCell>
						<TableCell>
							<span class="font-mono text-xs text-muted-foreground">{key.key_prefix}...</span>
						</TableCell>
						<TableCell>
							<div class="flex flex-wrap gap-1">
								{#if key.scopes.includes('*')}
									<Badge variant="solid">*</Badge>
								{:else}
									{#each key.scopes as scope (scope)}
										<Badge>{scope}</Badge>
									{/each}
								{/if}
							</div>
						</TableCell>
						<TableCell class="text-muted-foreground">
							{#if key.project_id}
								{data.projects.find(p => p.id === key.project_id)?.name ?? `#${key.project_id}`}
							{:else}
								<span class="text-xs">Global</span>
							{/if}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{key.last_used_at ? new Date(key.last_used_at).toLocaleDateString() : '—'}
						</TableCell>
						<TableCell class="text-right">
							<Button
								variant="ghost-destructive"
								size="sm"
								onclick={() => { pendingDeleteId = key.id; confirmOpen = true }}
							>
								Revoke
							</Button>
						</TableCell>
					</TableRow>
				{:else}
					<TableRow>
						<TableCell colspan={6} class="py-10 text-center text-muted-foreground">
							No API keys yet. Click "New API Key" to create one.
						</TableCell>
					</TableRow>
				{/each}
			</TableBody>
		</Table>
	</div>
</div>

<!-- Create API Key modal -->
<Dialog bind:open={createOpen} title="New API Key" size="md">
	{#snippet children()}
		<form
			id="create-key-form"
			method="POST"
			action="?/create"
			use:enhance={handleCreate}
			class="flex flex-col gap-4"
		>
			<Input name="name" type="text" label="Name" placeholder="e.g. CI/CD pipeline" required />

			<div>
				<p class="mb-2 text-sm font-medium text-foreground">Scopes</p>
				<div class="grid grid-cols-2 gap-1.5">
					{#each ALL_SCOPES as scope}
						<label class="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-sm text-foreground hover:bg-muted">
							<input type="checkbox" name="scopes" value={scope.value} class="accent-primary" />
							{scope.label}
						</label>
					{/each}
				</div>
			</div>

			{#if data.projects.length > 0}
				<div>
					<label for="project_id" class="mb-1 block text-sm font-medium text-foreground">
						Restrict to project <span class="font-normal text-muted-foreground">(optional)</span>
					</label>
					<select
						id="project_id"
						name="project_id"
						class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-1 focus:ring-ring"
					>
						<option value="">All projects (global)</option>
						{#each data.projects as project}
							<option value={project.id}>{project.name}</option>
						{/each}
					</select>
				</div>
			{/if}

			{#if createError}
				<p class="text-sm text-destructive">{createError}</p>
			{/if}
		</form>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (createOpen = false)}>Cancel</Button>
		<Button type="submit" form="create-key-form">Create</Button>
	{/snippet}
</Dialog>

<!-- Show token once modal -->
{#if createdKey}
	<Dialog
		open={true}
		title="Save your API Key"
		description="This key will only be shown once. Copy it now and store it somewhere safe."
	>
		{#snippet children()}
			<div class="mt-2 flex items-center gap-2 rounded-md border border-border bg-muted px-3 py-2">
				<code class="flex-1 break-all font-mono text-xs text-foreground">{createdKey?.token}</code>
				<Button variant="outline" size="sm" onclick={copyToken}>
					{tokenCopied ? 'Copied!' : 'Copy'}
				</Button>
			</div>
		{/snippet}
		{#snippet footer()}
			<Button onclick={() => (createdKey = null)}>Done</Button>
		{/snippet}
	</Dialog>
{/if}

<!-- Revoke confirm modal -->
<Dialog
	bind:open={confirmOpen}
	title="Revoke API Key"
	description="This will permanently revoke the key. Any integrations using it will stop working immediately."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>Cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="id" value={pendingDeleteId} />
			<Button type="submit" variant="destructive">Revoke</Button>
		</form>
	{/snippet}
</Dialog>

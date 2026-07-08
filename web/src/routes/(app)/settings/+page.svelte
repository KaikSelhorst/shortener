<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import type { CreateApiKeyResponse, TOTPSetupResponse } from '$lib/types'
	import { enhance } from '$app/forms'
	import { formatDate } from '$lib/format'
	import {
		Badge,
		Button,
		CopyInput,
		Input,
		Dialog,
		QRCode,
	} from '$lib'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	// ── TOTP ──────────────────────────────────────────────────────────────────
	const totpEnabled = $derived(data.me.totp_enabled)
	let totpSetupData = $state<TOTPSetupResponse | null>(null)
	let totpSetupOpen = $state(false)
	let totpDisableOpen = $state(false)
	let totpError = $state<string | null>(null)
	$effect(() => {
		if (form && 'totpSetup' in form && form.totpSetup) {
			totpSetupData = form.totpSetup as TOTPSetupResponse
		}
		if (form && 'totpEnabled' in form && form.totpEnabled) {
			totpSetupOpen = false
			totpSetupData = null
		}
		if (form && 'totpDisabled' in form && form.totpDisabled) {
			totpDisableOpen = false
		}
		if (form && 'totpError' in form) {
			totpError = (form.totpError as string) ?? null
		} else {
			totpError = null
		}
	})

	$effect(() => {
		if (!totpSetupOpen) {
			totpSetupData = null
			totpError = null
		}
	})

	$effect(() => {
		if (!totpDisableOpen) totpError = null
	})
	// ─────────────────────────────────────────────────────────────────────────

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
		{ value: 'webhooks:read', label: 'Webhooks — Read' },
		{ value: 'webhooks:create', label: 'Webhooks — Create' },
		{ value: 'webhooks:delete', label: 'Webhooks — Delete' },
	]

	let createOpen = $state(false)
	let createError = $state<string | null>(null)
	let createdKey = $state<CreateApiKeyResponse | null>(null)
	let tokenOpen = $state(false)
	$effect(() => {
		if (!createOpen) createError = null
	})

	$effect(() => {
		if (!tokenOpen) createdKey = null
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
					tokenOpen = true
					createOpen = false
				}
				await update()
			}
		}
	}

	let pendingDeleteId = $state<number | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingDeleteId = null
	})
</script>

<div class="sticky top-0 z-10 bg-background border-b border-border px-4 h-11 flex items-center">
	<span class="text-sm font-medium text-foreground">Settings</span>
</div>

<!-- API Keys -->
<div class="border-b border-border px-4 h-10 flex items-center justify-between shrink-0">
	<span class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">API Keys</span>
	<Button size="sm" onclick={() => (createOpen = true)}>+ New key</Button>
</div>

{#if form?.deleteError}
	<div class="border-b border-border px-4 py-2">
		<p class="text-xs text-destructive">{form.deleteError}</p>
	</div>
{/if}

<table class="w-full text-sm">
	<thead class="border-b border-border">
		<tr>
			<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Name</th>
			<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Prefix</th>
			<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Scopes</th>
			<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Project</th>
			<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Last used</th>
			<th class="px-4 py-1.5"></th>
		</tr>
	</thead>
	<tbody>
		{#each data.apiKeys as key (key.id)}
			<tr class="group hover:bg-secondary/40 transition-colors">
				<td class="px-4 py-1.5 font-medium text-sm">{key.name}</td>
				<td class="px-4 py-1.5 text-muted-foreground font-mono text-xs">{key.key_prefix}…</td>
				<td class="px-4 py-1.5">
					<div class="flex flex-wrap gap-1">
						{#if key.scopes.includes('*')}
							<Badge variant="solid">*</Badge>
						{:else}
							{#each key.scopes as scope}<Badge>{scope}</Badge>{/each}
						{/if}
					</div>
				</td>
				<td class="px-4 py-1.5 text-sm text-muted-foreground">
					{#if key.project_id}
						{data.projects.find((p) => p.id === key.project_id)?.name ?? `#${key.project_id}`}
					{:else}
						Global
					{/if}
				</td>
				<td class="px-4 py-1.5 text-sm text-muted-foreground">{formatDate(key.last_used_at)}</td>
				<td class="px-4 py-1.5 text-right">
					<div class="opacity-0 group-hover:opacity-100 transition-opacity">
						<Button variant="ghost-destructive" size="sm" onclick={() => { pendingDeleteId = key.id; confirmOpen = true }}>Revoke</Button>
					</div>
				</td>
			</tr>
		{:else}
			<tr>
				<td colspan={6} class="px-4 py-12 text-center text-sm text-muted-foreground">No API keys yet.</td>
			</tr>
		{/each}
	</tbody>
</table>

<!-- Two-factor auth -->
<div class="border-t border-border px-4 py-4 flex items-center justify-between">
	<div>
		<p class="text-sm font-medium text-foreground">Two-factor authentication</p>
		<p class="text-xs text-muted-foreground mt-0.5">
			{#if totpEnabled}Active — your account requires an authenticator code at sign-in.
			{:else}Add an extra layer of security with a TOTP authenticator app.{/if}
		</p>
	</div>
	{#if totpEnabled}
		<div class="flex items-center gap-3">
			<span class="flex items-center gap-1.5 text-xs text-success"><span class="w-1.5 h-1.5 rounded-full bg-success"></span>Active</span>
			<Button variant="outline" size="sm" onclick={() => (totpDisableOpen = true)}>Disable</Button>
		</div>
	{:else}
		<form method="POST" action="?/totpSetup" use:enhance={() => { return async ({ result, update }) => { await update(); if (result.type === 'success') totpSetupOpen = true } }}>
			<Button type="submit" size="sm">Enable</Button>
		</form>
	{/if}
</div>

<!-- TOTP setup modal -->
<Dialog bind:open={totpSetupOpen} title="Enable two-factor authentication" size="md">
	{#snippet children()}
		<div class="flex flex-col gap-4">
			{#if totpSetupData}
				<p class="text-sm text-muted-foreground">
					Scan the QR code with your authenticator app, then enter the 6-digit code to confirm.
				</p>

				<div class="flex justify-center">
					<div class="border border-border bg-white p-2">
						<QRCode data={totpSetupData.uri} moduleSize={4} border={4} />
					</div>
				</div>

				<details class="text-sm">
					<summary class="cursor-pointer text-muted-foreground hover:text-foreground transition-colors">
						Can't scan? Enter key manually
					</summary>
					<div class="mt-2">
						<CopyInput value={totpSetupData.secret} />
					</div>
				</details>

				{#if totpError}
					<p class="text-sm text-destructive">{totpError}</p>
				{/if}

				<form
					id="totp-confirm-form"
					method="POST"
					action="?/totpConfirm"
					use:enhance={() => ({ update }) => update()}
				>
					<Input
						name="code"
						type="text"
						label="Authenticator code"
						inputmode="numeric"
						pattern="[0-9]*"
						maxlength={6}
						autocomplete="one-time-code"
						required
					/>
				</form>
			{/if}
		</div>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (totpSetupOpen = false)}>Cancel</Button>
		<Button type="submit" form="totp-confirm-form">Confirm</Button>
	{/snippet}
</Dialog>

<!-- TOTP disable modal -->
<Dialog
	bind:open={totpDisableOpen}
	title="Disable two-factor authentication"
	description="Enter your current authenticator code to confirm."
>
	{#snippet children()}
		<form
			id="totp-disable-form"
			method="POST"
			action="?/totpDisable"
			use:enhance={() => ({ update }) => update()}
			class="flex flex-col gap-4"
		>
			{#if totpError}
				<p class="text-sm text-destructive">{totpError}</p>
			{/if}
			<Input
				name="code"
				type="text"
				label="Authenticator code"
				inputmode="numeric"
				pattern="[0-9]*"
				maxlength={6}
				autocomplete="one-time-code"
				required
			/>
		</form>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (totpDisableOpen = false)}>Cancel</Button>
		<Button type="submit" form="totp-disable-form" variant="destructive">Disable</Button>
	{/snippet}
</Dialog>

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
				<p class="tui-label mb-2">Scopes</p>
				<div class="grid grid-cols-2 gap-1">
					{#each ALL_SCOPES as scope}
						<label class="flex cursor-pointer items-center gap-2 px-2 py-1.5 text-sm text-foreground hover:bg-card transition-colors">
							<input type="checkbox" name="scopes" value={scope.value} />
							{scope.label}
						</label>
					{/each}
				</div>
			</div>

			{#if data.projects.length > 0}
				<div>
					<label for="project_id" class="tui-label mb-1.5 block">
						Restrict to project <span class="normal-case text-muted-foreground">(optional)</span>
					</label>
					<select id="project_id" name="project_id" class="w-full">
						<option value="">all projects (global)</option>
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
<Dialog
	bind:open={tokenOpen}
	title="Save your API Key"
	description="This key will only be shown once. Copy it now and store it somewhere safe."
>
	{#snippet children()}
		<CopyInput value={createdKey?.token ?? ''} />
	{/snippet}
	{#snippet footer()}
		<Button onclick={() => (tokenOpen = false)}>Done</Button>
	{/snippet}
</Dialog>

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

<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import type { SubmitFunction } from '@sveltejs/kit'
	import type { CreateApiKeyResponse, TOTPSetupResponse } from '$lib/types'
	import { enhance } from '$app/forms'
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

	// ── TOTP ──────────────────────────────────────────────────────────────────
	const totpEnabled = $derived(data.me.totp_enabled)
	let totpSetupData = $state<TOTPSetupResponse | null>(null)
	let totpSetupOpen = $state(false)
	let totpDisableOpen = $state(false)
	let totpError = $state<string | null>(null)
	let totpSecretCopied = $state(false)

	async function copyTotpSecret() {
		if (!totpSetupData) return
		await navigator.clipboard.writeText(totpSetupData.secret)
		totpSecretCopied = true
		setTimeout(() => (totpSecretCopied = false), 2000)
	}

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
	]

	let createOpen = $state(false)
	let createError = $state<string | null>(null)
	let createdKey = $state<CreateApiKeyResponse | null>(null)
	let tokenOpen = $state(false)
	let tokenCopied = $state(false)

	$effect(() => {
		if (!createOpen) createError = null
	})

	$effect(() => {
		if (!tokenOpen) {
			createdKey = null
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
					tokenOpen = true
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

<div class="mx-auto max-w-6xl space-y-6">
	<!-- API Keys -->
	<div class="tui-panel">
		<div class="tui-panel-header justify-between">
			<span>▌ api keys</span>
			<Button size="sm" onclick={() => (createOpen = true)}>+ new key</Button>
		</div>

		{#if form?.deleteError}
			<div class="border-b border-border px-4 py-2">
				<p class="font-mono text-xs text-destructive">{form.deleteError}</p>
			</div>
		{/if}

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
							<span class="text-muted-foreground">{key.key_prefix}...</span>
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
								{data.projects.find((p) => p.id === key.project_id)?.name ?? `#${key.project_id}`}
							{:else}
								<span>global</span>
							{/if}
						</TableCell>
						<TableCell class="text-muted-foreground">
							{key.last_used_at ? new Date(key.last_used_at).toLocaleDateString() : '—'}
						</TableCell>
						<TableCell class="text-right">
							<Button
								variant="ghost-destructive"
								size="sm"
								onclick={() => {
									pendingDeleteId = key.id
									confirmOpen = true
								}}
							>
								revoke
							</Button>
						</TableCell>
					</TableRow>
				{:else}
					<TableRow>
						<TableCell colspan={6} class="py-12 text-center text-muted-foreground">
							-- no api keys yet --
						</TableCell>
					</TableRow>
				{/each}
			</TableBody>
		</Table>
	</div>

	<!-- Two-factor auth -->
	<div class="tui-panel">
		<div class="tui-panel-header justify-between">
			<span>▌ two-factor auth</span>
			{#if totpEnabled}
				<div class="flex items-center gap-3">
					<span class="font-mono text-[10px] text-primary">● active</span>
					<Button variant="outline" size="sm" onclick={() => (totpDisableOpen = true)}>
						disable
					</Button>
				</div>
			{:else}
				<form
					method="POST"
					action="?/totpSetup"
					use:enhance={() => {
						return async ({ result, update }) => {
							await update()
							if (result.type === 'success') totpSetupOpen = true
						}
					}}
				>
					<Button type="submit" size="sm">enable</Button>
				</form>
			{/if}
		</div>
		<div class="px-4 py-3">
			<p class="font-mono text-xs text-muted-foreground">
				{#if totpEnabled}
					totp is active — your account requires an authenticator code at sign-in.
				{:else}
					add an extra layer of security with a totp authenticator app.
				{/if}
			</p>
		</div>
	</div>
</div>

<!-- TOTP setup modal -->
<Dialog bind:open={totpSetupOpen} title="Enable two-factor authentication" size="md">
	{#snippet children()}
		<div class="flex flex-col gap-4">
			{#if totpSetupData}
				<p class="font-mono text-xs text-muted-foreground">
					Scan the QR code with your authenticator app, then enter the 6-digit code to confirm.
				</p>

				<div class="flex justify-center">
					<div class="border border-border bg-white p-2">
						<QRCode data={totpSetupData.uri} moduleSize={4} border={4} />
					</div>
				</div>

				<details class="font-mono text-xs">
					<summary
						class="cursor-pointer text-muted-foreground hover:text-foreground transition-colors"
					>
						[expand] can't scan? enter key manually
					</summary>
					<div class="mt-2 flex items-center gap-2 border border-border bg-card px-3 py-2">
						<code class="flex-1 break-all text-[10px] text-foreground">
							{totpSetupData.secret}
						</code>
						<Button variant="outline" size="sm" onclick={copyTotpSecret}>
							{totpSecretCopied ? 'copied!' : 'copy'}
						</Button>
					</div>
				</details>

				{#if totpError}
					<p class="font-mono text-xs text-destructive">{totpError}</p>
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
		<Button variant="outline" onclick={() => (totpSetupOpen = false)}>cancel</Button>
		<Button type="submit" form="totp-confirm-form">confirm</Button>
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
				<p class="font-mono text-xs text-destructive">{totpError}</p>
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
		<Button variant="outline" onclick={() => (totpDisableOpen = false)}>cancel</Button>
		<Button type="submit" form="totp-disable-form" variant="destructive">disable</Button>
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
						<label
							class="flex cursor-pointer items-center gap-2 px-2 py-1.5 font-mono text-xs text-foreground hover:bg-card transition-colors"
						>
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
				<p class="font-mono text-xs text-destructive">{createError}</p>
			{/if}
		</form>
	{/snippet}
	{#snippet footer()}
		<Button variant="outline" onclick={() => (createOpen = false)}>cancel</Button>
		<Button type="submit" form="create-key-form">create</Button>
	{/snippet}
</Dialog>

<!-- Show token once modal -->
<Dialog
	bind:open={tokenOpen}
	title="Save your API Key"
	description="This key will only be shown once. Copy it now and store it somewhere safe."
>
	{#snippet children()}
		<div class="flex items-center gap-2 border border-border bg-card px-3 py-2">
			<code class="flex-1 break-all font-mono text-[10px] text-foreground">{createdKey?.token}</code>
			<Button variant="outline" size="sm" onclick={copyToken}>
				{tokenCopied ? 'copied!' : 'copy'}
			</Button>
		</div>
	{/snippet}
	{#snippet footer()}
		<Button onclick={() => (tokenOpen = false)}>done</Button>
	{/snippet}
</Dialog>

<!-- Revoke confirm modal -->
<Dialog
	bind:open={confirmOpen}
	title="Revoke API Key"
	description="This will permanently revoke the key. Any integrations using it will stop working immediately."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="id" value={pendingDeleteId} />
			<Button type="submit" variant="destructive">revoke</Button>
		</form>
	{/snippet}
</Dialog>

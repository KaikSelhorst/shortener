<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'

	let { form }: { form: ActionData } = $props()

	let session = $derived(form && 'session' in form ? (form.session as string) : null)
</script>

<div class="tui-label mb-5 text-accent">
	{#if session}mfa verification{:else}sign in{/if}
</div>

{#if form?.error}
	<p class="mb-4 font-mono text-xs text-destructive">{form.error}</p>
{/if}

{#if session}
	<form method="POST" action="?/totp" class="flex flex-col gap-5">
		<input type="hidden" name="session" value={session} />
		<p class="font-mono text-xs text-muted-foreground">
			Enter the 6-digit code from your authenticator app.
		</p>
		<Input
			label="Code"
			name="code"
			type="text"
			inputmode="numeric"
			pattern="[0-9]*"
			maxlength={6}
			autocomplete="one-time-code"
			required
		/>
		<Button type="submit" class="mt-1 w-full">verify</Button>
	</form>
{:else}
	<form method="POST" action="?/credentials" class="flex flex-col gap-5">
		<Input label="Email" name="email" type="email" required autocomplete="email" />
		<Input
			label="Password"
			name="password"
			type="password"
			required
			autocomplete="current-password"
		/>
		<Button type="submit" class="mt-1 w-full">sign in</Button>
		<p class="text-center font-mono text-xs text-muted-foreground">
			no account? <a href="/register" class="text-accent hover:underline">register</a>
		</p>
	</form>
{/if}

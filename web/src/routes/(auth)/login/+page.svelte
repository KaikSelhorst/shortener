<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'
	let { form }: { form: ActionData } = $props()
	let session = $derived(form && 'session' in form ? (form.session as string) : null)
</script>

{#if session}
	<form method="POST" action="?/totp" class="flex flex-col gap-3">
		<input type="hidden" name="session" value={session} />
		<div class="mb-1">
			<h1 class="text-sm font-semibold text-foreground">Two-factor authentication</h1>
			<p class="mt-1 text-sm text-muted-foreground">Enter the 6-digit code from your authenticator app.</p>
		</div>
		{#if form?.error}<p class="text-sm text-destructive">{form.error}</p>{/if}
		<Input label="Code" name="code" type="text" inputmode="numeric" pattern="[0-9]*" maxlength={6} autocomplete="one-time-code" required />
		<Button type="submit" class="w-full mt-1">Verify</Button>
	</form>
{:else}
	<form method="POST" action="?/credentials" class="flex flex-col gap-3">
		<div class="mb-1">
			<h1 class="text-sm font-semibold text-foreground">Sign in</h1>
			<p class="mt-1 text-sm text-muted-foreground">Enter your credentials to continue.</p>
		</div>
		{#if form?.error}<p class="text-sm text-destructive">{form.error}</p>{/if}
		<Input label="Email" name="email" type="email" required autocomplete="email" />
		<Input label="Password" name="password" type="password" required autocomplete="current-password" />
		<Button type="submit" class="w-full mt-1">Sign in</Button>
		<p class="text-sm text-muted-foreground">
			No account? <a href="/register" class="text-foreground hover:underline underline-offset-4">Register</a>
		</p>
	</form>
{/if}

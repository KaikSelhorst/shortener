<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'

	let { form }: { form: ActionData } = $props()

	// After step 1 succeeds with TOTP required, the server returns { session }.
	// We keep session in local state so step 2 can carry it as a hidden field.
	let session = $derived(form && 'session' in form ? (form.session as string) : null)
</script>

<h1 class="mb-6 text-xl font-semibold text-foreground">Sign in</h1>

{#if form?.error}
	<p class="mb-4 rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">{form.error}</p>
{/if}

{#if session}
	<!-- Step 2: TOTP code -->
	<form method="POST" action="?/totp" class="flex flex-col gap-4">
		<input type="hidden" name="session" value={session} />

		<div>
			<p class="mb-1 text-sm font-medium text-foreground">Authenticator code</p>
			<p class="mb-3 text-sm text-muted-foreground">
				Open your authenticator app and enter the 6-digit code.
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
		</div>

		<Button type="submit" class="mt-2 w-full">Verify</Button>
	</form>
{:else}
	<!-- Step 1: credentials -->
	<form method="POST" action="?/credentials" class="flex flex-col gap-4">
		<Input label="Email" name="email" type="email" required autocomplete="email" />

		<Input
			label="Password"
			name="password"
			type="password"
			required
			autocomplete="current-password"
		/>

		<Button type="submit" class="mt-2 w-full">Sign in</Button>

		<p class="text-center text-sm text-muted-foreground">
			No account? <a href="/register" class="text-foreground underline">Register</a>
		</p>
	</form>
{/if}

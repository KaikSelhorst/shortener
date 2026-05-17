<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'

	let { form }: { form: ActionData } = $props()
</script>

<h1 class="mb-6 text-xl font-semibold text-foreground">Sign in</h1>

<form method="POST" class="flex flex-col gap-4">
	{#if form?.error}
		<p class="rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">{form.error}</p>
	{/if}

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

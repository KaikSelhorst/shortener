<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'

	let { form }: { form: ActionData } = $props()
</script>

<h1 class="mb-6 text-xl font-semibold text-foreground">Create account</h1>

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
		minlength={8}
		autocomplete="new-password"
	/>

	<Button type="submit" class="mt-2 w-full">Create account</Button>

	<p class="text-center text-sm text-muted-foreground">
		Already have an account? <a href="/login" class="text-foreground underline">Sign in</a>
	</p>
</form>

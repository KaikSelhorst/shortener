<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'

	let { form }: { form: ActionData } = $props()
</script>

<div class="tui-label mb-5 text-accent">create account</div>

<form method="POST" class="flex flex-col gap-5">
	{#if form?.error}
		<p class="font-mono text-xs text-destructive">{form.error}</p>
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

	<Button type="submit" class="mt-1 w-full">create account</Button>

	<p class="text-center font-mono text-xs text-muted-foreground">
		already have an account? <a href="/login" class="text-accent hover:underline">sign in</a>
	</p>
</form>

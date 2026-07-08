<script lang="ts">
	import type { ActionData } from './$types'
	import { Button, Input } from '$lib'
	let { form }: { form: ActionData } = $props()
</script>

<form method="POST" class="flex flex-col gap-3">
	<div class="mb-1">
		<h1 class="text-sm font-semibold text-foreground">Create account</h1>
		<p class="mt-1 text-sm text-muted-foreground">Fill in the details below to get started.</p>
	</div>
	{#if form?.error}<p class="text-sm text-destructive">{form.error}</p>{/if}
	<Input label="Email" name="email" type="email" required autocomplete="email" />
	<Input label="Password" name="password" type="password" required minlength={8} autocomplete="new-password" />
	<Button type="submit" class="w-full mt-1">Create account</Button>
	<p class="text-sm text-muted-foreground">
		Already have an account? <a href="/login" class="text-foreground hover:underline underline-offset-4">Sign in</a>
	</p>
</form>

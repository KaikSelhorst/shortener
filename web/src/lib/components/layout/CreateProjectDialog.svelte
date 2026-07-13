<script lang="ts">
	import { goto, invalidate } from '$app/navigation'
	import { Button, Input, Dialog } from '$lib/components/ui'

	interface Props {
		open: boolean
		dismissable?: boolean
	}

	let { open = $bindable(), dismissable = true }: Props = $props()

	let error = $state<string | null>(null)
	let submitting = $state(false)

	$effect(() => {
		if (!open) error = null
	})

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault()
		if (submitting) return

		const form = e.currentTarget as HTMLFormElement
		const name = (new FormData(form).get('name') as string)?.trim()
		if (!name) return

		submitting = true
		error = null

		try {
			const res = await fetch('/api/projects', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ name }),
			})
			const body = await res.json()

			if (!res.ok) {
				error = body?.error ?? 'Failed to create project'
				return
			}

			form.reset()
			open = false
			await invalidate('app:projects')
			await goto(`/${body.slug}`)
		} finally {
			submitting = false
		}
	}
</script>

<Dialog bind:open title="New project" {dismissable}>
	{#snippet children()}
		<form id="create-project-form" onsubmit={handleSubmit} class="flex flex-col gap-4">
			<Input
				name="name"
				type="text"
				label="Name"
				placeholder="e.g. Marketing"
				required
			/>
			{#if error}
				<p class="text-sm text-destructive">{error}</p>
			{/if}
		</form>
	{/snippet}
	{#snippet footer()}
		{#if dismissable}
			<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
		{/if}
		<Button type="submit" form="create-project-form" disabled={submitting}>Create</Button>
	{/snippet}
</Dialog>

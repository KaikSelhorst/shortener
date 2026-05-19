<script lang="ts">
	import type { Snippet } from 'svelte'

	interface Props {
		open: boolean
		title: string
		description?: string
		size?: 'sm' | 'md'
		children?: Snippet
		footer?: Snippet
	}

	let { open = $bindable(), title, description, size = 'sm', children, footer }: Props = $props()

	let dialog = $state<HTMLDialogElement>()
	let inner = $state<HTMLDivElement>()

	const widths: Record<'sm' | 'md', string> = { sm: 'max-w-sm', md: 'max-w-md' }
	const dialogCls = $derived(
		`m-auto w-full ${widths[size]} rounded-md border border-border bg-card shadow-lg backdrop:bg-black/50`
	)

	$effect(() => {
		if (!dialog) return
		if (open) dialog.showModal()
		else dialog.close()
	})
</script>

<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<dialog
	bind:this={dialog}
	onclose={() => (open = false)}
	onmousedown={(e) => { if (!inner?.contains(e.target as Node)) open = false }}
	class={dialogCls}
>
	<div bind:this={inner} class="p-6">
		<h2 class="text-base font-semibold text-foreground">{title}</h2>

		{#if description}
			<p class="mt-2 text-sm text-muted-foreground">{description}</p>
		{/if}

		{#if children}
			<div class="mt-4">{@render children()}</div>
		{/if}

		{#if footer}
			<div class="mt-6 flex justify-end gap-2">{@render footer()}</div>
		{/if}
	</div>
</dialog>

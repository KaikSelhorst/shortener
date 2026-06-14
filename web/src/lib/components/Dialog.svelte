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
	class="m-auto w-full {widths[size]} border border-border bg-background p-0 shadow-none backdrop:bg-black/75 open:flex open:flex-col"
>
	<div bind:this={inner} class="flex flex-col">
		<div class="flex items-center justify-between border-b border-border px-4 py-2.5">
			<span class="font-mono text-[10px] uppercase tracking-widest text-accent">▌ {title}</span>
			<button
				onclick={() => (open = false)}
				class="font-mono text-[10px] text-muted-foreground hover:text-foreground transition-colors"
				aria-label="Close"
			>[×]</button>
		</div>

		<div class="p-5">
			{#if description}
				<p class="mb-4 font-mono text-xs text-muted-foreground">{description}</p>
			{/if}

			{#if children}
				<div>{@render children()}</div>
			{/if}

			{#if footer}
				<div class="mt-5 flex justify-end gap-2 border-t border-border pt-4">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
</dialog>

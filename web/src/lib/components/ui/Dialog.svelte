<script lang="ts">
	import type { Snippet } from 'svelte'

	interface Props {
		open: boolean
		title: string
		description?: string
		size?: 'sm' | 'md'
		dismissable?: boolean
		children?: Snippet
		footer?: Snippet
	}

	let { open = $bindable(), title, description, size = 'sm', dismissable = true, children, footer }: Props = $props()

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
	oncancel={(e) => { if (!dismissable) e.preventDefault() }}
	onmousedown={(e) => { if (dismissable && !inner?.contains(e.target as Node)) open = false }}
	class="m-auto w-full {widths[size]} border border-border bg-card rounded-xl p-0 shadow-2xl backdrop:bg-black/60 open:flex open:flex-col"
>
	<div bind:this={inner} class="flex flex-col">
		<div class="flex items-center justify-between border-b border-border px-4 py-3">
			<span class="text-sm font-semibold text-foreground">{title}</span>
			{#if dismissable}
				<button
					onclick={() => (open = false)}
					class="flex items-center justify-center w-6 h-6 rounded-md text-muted-foreground hover:text-foreground hover:bg-secondary transition-colors"
					aria-label="Close"
				>
					<svg width="14" height="14" viewBox="0 0 14 14" fill="none">
						<path d="M1 1l12 12M13 1L1 13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
					</svg>
				</button>
			{/if}
		</div>

		<div class="p-5">
			{#if description}
				<p class="mb-4 text-sm text-muted-foreground">{description}</p>
			{/if}

			{#if children}
				<div>{@render children()}</div>
			{/if}
		</div>

		{#if footer}
			<div class="flex justify-end gap-2 border-t border-border px-5 py-4">
				{@render footer()}
			</div>
		{/if}
	</div>
</dialog>

<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements'

	interface Props extends HTMLInputAttributes {
		label?: string
		error?: string
	}

	let { label, error, class: className, ...rest }: Props = $props()

	const base =
		'flex w-full rounded-md border-2 bg-background px-3 text-sm text-foreground ' +
		'placeholder:text-muted-foreground ' +
		'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ' +
		'disabled:cursor-not-allowed disabled:opacity-50'

	const state = $derived(error ? 'border-destructive focus-visible:ring-destructive' : 'border-input')

	const cls = $derived([base, state, 'h-9', className].filter(Boolean).join(' '))
</script>

{#if label}
	<label class="flex flex-col gap-1.5">
		<span class="text-sm font-medium text-foreground">{label}</span>
		<input class={cls} aria-invalid={error ? true : undefined} {...rest} />
		{#if error}
			<span class="text-xs text-destructive">{error}</span>
		{/if}
	</label>
{:else}
	<input class={cls} aria-invalid={error ? true : undefined} {...rest} />
	{#if error}
		<span class="text-xs text-destructive">{error}</span>
	{/if}
{/if}

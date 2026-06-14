<script lang="ts">
	import type { HTMLTextareaAttributes } from 'svelte/elements'

	interface Props extends HTMLTextareaAttributes {
		label?: string
		error?: string
	}

	let { label, error, class: className, ...rest }: Props = $props()

	const base =
		'flex w-full bg-transparent border border-input px-3 py-2 text-xs text-foreground ' +
		'font-mono placeholder:text-muted-foreground/40 ' +
		'focus:outline-none focus:border-primary ' +
		'disabled:opacity-30 disabled:cursor-not-allowed ' +
		'transition-colors resize-none'

	const borderState = $derived(error ? 'border-destructive focus:border-destructive' : '')
	const cls = $derived([base, borderState, className].filter(Boolean).join(' '))
</script>

{#if label}
	<label class="flex flex-col gap-1.5">
		<span class="tui-label">{label}</span>
		<textarea class={cls} aria-invalid={error ? true : undefined} {...rest}></textarea>
		{#if error}
			<span class="font-mono text-[10px] text-destructive">{error}</span>
		{/if}
	</label>
{:else}
	<textarea class={cls} aria-invalid={error ? true : undefined} {...rest}></textarea>
	{#if error}
		<span class="font-mono text-[10px] text-destructive">{error}</span>
	{/if}
{/if}

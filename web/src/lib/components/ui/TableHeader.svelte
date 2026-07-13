<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLThAttributes } from 'svelte/elements'

	interface Props extends HTMLThAttributes {
		children?: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<th
	class={[
		'px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider whitespace-nowrap',
		className,
	]
		.filter(Boolean)
		.join(' ')}
	{...rest}
>
	{#if children}{@render children()}{/if}
</th>

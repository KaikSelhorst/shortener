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
		'px-4 py-2.5 text-left text-[10px] font-normal uppercase tracking-[0.15em] text-muted-foreground',
		className,
	]
		.filter(Boolean)
		.join(' ')}
	{...rest}
>
	{#if children}{@render children()}{/if}
</th>

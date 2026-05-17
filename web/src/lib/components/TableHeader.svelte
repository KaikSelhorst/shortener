<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLThAttributes } from 'svelte/elements'

	interface Props extends HTMLThAttributes {
		children?: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<th
	class={['px-4 py-3 text-left font-medium text-muted-foreground', className]
		.filter(Boolean)
		.join(' ')}
	{...rest}
>
	{#if children}{@render children()}{/if}
</th>

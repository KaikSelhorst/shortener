<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLTdAttributes } from 'svelte/elements'

	interface Props extends HTMLTdAttributes {
		children: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<td
	class={['px-4 py-3 text-foreground', className].filter(Boolean).join(' ')}
	{...rest}
>
	{@render children()}
</td>

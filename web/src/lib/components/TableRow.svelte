<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLAttributes } from 'svelte/elements'

	interface Props extends HTMLAttributes<HTMLTableRowElement> {
		children: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<tr class={['group transition-colors hover:bg-secondary/40', className].filter(Boolean).join(' ')} {...rest}>
	{@render children()}
</tr>

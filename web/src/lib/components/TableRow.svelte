<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLAttributes } from 'svelte/elements'

	interface Props extends HTMLAttributes<HTMLTableRowElement> {
		children: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<tr class={['hover:bg-muted/50 transition-colors', className].filter(Boolean).join(' ')} {...rest}>
	{@render children()}
</tr>

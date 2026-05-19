<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLAttributes } from 'svelte/elements'

	interface Props extends HTMLAttributes<HTMLTableSectionElement> {
		children: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<tbody class={['divide-y divide-border', className].filter(Boolean).join(' ')} {...rest}>
	{@render children()}
</tbody>

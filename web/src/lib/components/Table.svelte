<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLTableAttributes } from 'svelte/elements'

	interface Props extends HTMLTableAttributes {
		children: Snippet
	}

	let { children, class: className, ...rest }: Props = $props()
</script>

<div class="w-full overflow-auto border-t border-border">
	<table class={['w-full', className].filter(Boolean).join(' ')} {...rest}>
		{@render children()}
	</table>
</div>

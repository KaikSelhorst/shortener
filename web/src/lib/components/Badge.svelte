<script lang="ts">
	import type { Snippet } from 'svelte'

	type Variant = 'default' | 'solid'

	interface Props {
		variant?: Variant
		class?: string
		children: Snippet
	}

	let { variant = 'default', class: className, children }: Props = $props()

	const variants: Record<Variant, string> = {
		default: 'bg-secondary text-secondary-foreground ring-1 ring-inset ring-border',
		solid:   'bg-primary text-primary-foreground',
	}

	const cls = $derived(
		[
			'inline-flex items-center rounded-sm px-1.5 py-0.5 font-mono text-xs font-medium',
			variants[variant],
			className,
		]
			.filter(Boolean)
			.join(' '),
	)
</script>

<span class={cls}>{@render children()}</span>

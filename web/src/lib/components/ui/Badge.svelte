<script lang="ts">
	import type { Snippet } from 'svelte'

	type Variant = 'default' | 'solid' | 'success' | 'error'

	interface Props {
		variant?: Variant
		class?: string
		children: Snippet
	}

	let { variant = 'default', class: className, children }: Props = $props()

	const variants: Record<Variant, string> = {
		default: 'border border-border text-muted-foreground',
		solid:   'border border-border bg-secondary text-foreground',
		success: 'border border-success/30 text-success bg-success/10',
		error:   'border border-destructive/30 text-destructive bg-destructive/10',
	}

	const cls = $derived(
		[
			'inline-flex items-center rounded-md px-2 py-0.5 text-[11px] font-medium leading-none',
			variants[variant],
			className,
		]
			.filter(Boolean)
			.join(' '),
	)
</script>

<span class={cls}>{@render children()}</span>

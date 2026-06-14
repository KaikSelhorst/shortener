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
		solid:   'border border-primary text-primary bg-primary/10',
		success: 'border border-tui-green/30 text-tui-green bg-tui-green/10',
		error:   'border border-destructive/30 text-destructive bg-destructive/10',
	}

	const cls = $derived(
		[
			'inline-flex items-center rounded-[2px] px-[0.45rem] py-[0.2rem] font-mono text-[10px] leading-none',
			variants[variant],
			className,
		]
			.filter(Boolean)
			.join(' '),
	)
</script>

<span class={cls}>{@render children()}</span>

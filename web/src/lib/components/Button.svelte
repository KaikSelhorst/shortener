<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLButtonAttributes } from 'svelte/elements'

	type Variant = 'primary' | 'secondary' | 'destructive' | 'ghost' | 'ghost-destructive' | 'outline'
	type Size = 'sm' | 'md'

	interface Props extends HTMLButtonAttributes {
		variant?: Variant
		size?: Size
		href?: string
		children: Snippet
	}

	let {
		variant = 'primary',
		size = 'md',
		type = 'button',
		href,
		class: className,
		children,
		...rest
	}: Props = $props()

	const base =
		'inline-flex cursor-pointer items-center justify-center font-mono uppercase tracking-wider ' +
		'border transition-colors ' +
		'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring ' +
		'disabled:pointer-events-none disabled:opacity-30'

	const variants: Record<Variant, string> = {
		primary:             'border-primary text-primary bg-transparent hover:bg-primary/10',
		secondary:           'border-border text-muted-foreground hover:border-accent hover:text-accent',
		destructive:         'border-destructive text-destructive bg-transparent hover:bg-destructive/10',
		ghost:               'border-transparent text-muted-foreground hover:text-foreground hover:border-border',
		'ghost-destructive': 'border-transparent text-destructive hover:border-destructive/40',
		outline:             'border-border text-muted-foreground hover:border-accent hover:text-accent',
	}

	const sizes: Record<Size, string> = {
		sm: 'h-7 gap-1 px-2.5 text-[10px]',
		md: 'h-8 gap-1.5 px-3 text-xs',
	}

	const cls = $derived([base, variants[variant], sizes[size], className].filter(Boolean).join(' '))
</script>

{#if href}
	<a {href} class={cls}>
		{@render children()}
	</a>
{:else}
	<button {type} class={cls} {...rest}>
		{@render children()}
	</button>
{/if}

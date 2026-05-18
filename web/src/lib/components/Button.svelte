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
		'inline-flex cursor-pointer items-center justify-center font-medium transition-colors ' +
		'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring ' +
		'disabled:pointer-events-none disabled:opacity-50'

	const variants: Record<Variant, string> = {
		primary: 'bg-primary text-primary-foreground hover:opacity-90',
		secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
		destructive: 'bg-destructive text-destructive-foreground hover:opacity-90',
		ghost: 'text-foreground hover:bg-muted hover:text-foreground',
		'ghost-destructive': 'text-destructive hover:bg-destructive/10 hover:text-destructive',
		outline: 'border-2 border-input bg-background text-foreground hover:bg-muted',
	}

	const sizes: Record<Size, string> = {
		sm: 'h-8 gap-1.5 px-3 text-xs rounded-sm',
		md: 'h-9 gap-2 px-4 text-sm rounded-md',
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

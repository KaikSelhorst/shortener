<script lang="ts">
	import type { Snippet } from 'svelte'
	import type { HTMLButtonAttributes } from 'svelte/elements'

	type Variant = 'default' | 'secondary' | 'outline' | 'ghost' | 'destructive' | 'ghost-destructive'
	type Size = 'sm' | 'md'

	interface Props extends HTMLButtonAttributes {
		variant?: Variant
		size?: Size
		href?: string
		children: Snippet
	}

	let {
		variant = 'default',
		size = 'md',
		type = 'button',
		href,
		class: className,
		children,
		...rest
	}: Props = $props()

	const base =
		'inline-flex cursor-pointer items-center justify-center font-medium rounded ' +
		'transition-colors select-none ' +
		'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background ' +
		'disabled:pointer-events-none disabled:opacity-40'

	const variants: Record<Variant, string> = {
		default:            'bg-primary text-primary-foreground hover:bg-primary/90',
		secondary:          'bg-secondary text-secondary-foreground border border-border hover:bg-secondary/70',
		outline:            'border border-border bg-transparent text-foreground hover:bg-secondary',
		ghost:              'border border-transparent bg-transparent text-foreground hover:bg-secondary',
		destructive:        'border border-destructive/40 text-destructive bg-destructive/5 hover:bg-destructive/10',
		'ghost-destructive':'border border-transparent bg-transparent text-destructive hover:bg-destructive/10',
	}

	const sizes: Record<Size, string> = {
		sm: 'h-7 gap-1 px-2.5 text-[12px]',
		md: 'h-8 gap-1.5 px-3 text-[13px]',
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

<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements'

	interface Props extends HTMLInputAttributes {
		label?: string
		error?: string
	}

	let { label, error, class: className, type = 'text', ...rest }: Props = $props()

	const isPassword = $derived(type === 'password')
	let showPassword = $state(false)
	const inputType = $derived(isPassword ? (showPassword ? 'text' : 'password') : type)

	const base =
		'flex w-full bg-transparent border-b border-input px-0 py-1.5 text-xs text-foreground ' +
		'font-mono placeholder:text-muted-foreground/40 ' +
		'focus:outline-none focus:border-primary ' +
		'disabled:opacity-30 disabled:cursor-not-allowed ' +
		'transition-colors'

	const borderState = $derived(error ? 'border-destructive focus:border-destructive' : '')
	const cls = $derived([base, borderState, isPassword ? 'pr-14' : '', className].filter(Boolean).join(' '))
</script>

{#if label}
	<label class="flex flex-col gap-1.5">
		<span class="tui-label">{label}</span>
		<div class="relative">
			<input class={cls} type={inputType} aria-invalid={error ? true : undefined} {...rest} />
			{#if isPassword}
				<button
					type="button"
					tabindex="-1"
					onclick={() => (showPassword = !showPassword)}
					class="absolute right-0 top-1/2 -translate-y-1/2 font-mono text-[9px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
					aria-label={showPassword ? 'Hide password' : 'Show password'}
				>
					[{showPassword ? 'hide' : 'show'}]
				</button>
			{/if}
		</div>
		{#if error}
			<span class="font-mono text-[10px] text-destructive">{error}</span>
		{/if}
	</label>
{:else}
	<div class="relative">
		<input class={cls} type={inputType} aria-invalid={error ? true : undefined} {...rest} />
		{#if isPassword}
			<button
				type="button"
				tabindex="-1"
				onclick={() => (showPassword = !showPassword)}
				class="absolute right-0 top-1/2 -translate-y-1/2 font-mono text-[9px] uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				[{showPassword ? 'hide' : 'show'}]
			</button>
		{/if}
	</div>
	{#if error}
		<span class="font-mono text-[10px] text-destructive">{error}</span>
	{/if}
{/if}

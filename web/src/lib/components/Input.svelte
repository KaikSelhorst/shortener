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
		'flex w-full rounded-md border bg-background px-3 text-sm text-foreground ' +
		'placeholder:text-muted-foreground ' +
		'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring ' +
		'disabled:cursor-not-allowed disabled:opacity-50'

	const borderState = $derived(error ? 'border-destructive focus-visible:ring-destructive' : 'border-input')

	const cls = $derived([base, borderState, 'h-9', isPassword ? 'pr-9' : '', className].filter(Boolean).join(' '))
</script>

{#if label}
	<label class="flex flex-col gap-1.5">
		<span class="text-sm font-medium text-foreground">{label}</span>
		<div class="relative">
			<input class={cls} type={inputType} aria-invalid={error ? true : undefined} {...rest} />
			{#if isPassword}
				<button
					type="button"
					tabindex="-1"
					onclick={() => (showPassword = !showPassword)}
					class="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
					aria-label={showPassword ? 'Hide password' : 'Show password'}
				>
					{#if showPassword}
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
							<path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
							<line x1="1" y1="1" x2="23" y2="23" />
						</svg>
					{:else}
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
							<circle cx="12" cy="12" r="3" />
						</svg>
					{/if}
				</button>
			{/if}
		</div>
		{#if error}
			<span class="text-xs text-destructive">{error}</span>
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
				class="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				{#if showPassword}
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
						<path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
						<line x1="1" y1="1" x2="23" y2="23" />
					</svg>
				{:else}
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
						<circle cx="12" cy="12" r="3" />
					</svg>
				{/if}
			</button>
		{/if}
	</div>
	{#if error}
		<span class="text-xs text-destructive">{error}</span>
	{/if}
{/if}

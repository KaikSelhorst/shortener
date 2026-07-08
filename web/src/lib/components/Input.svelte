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
		'flex w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm text-foreground ' +
		'placeholder:text-muted-foreground/50 ' +
		'focus:outline-none focus:border-ring focus:ring-1 focus:ring-ring ' +
		'disabled:opacity-40 disabled:cursor-not-allowed ' +
		'transition-colors h-9'

	const borderState = $derived(error ? 'border-destructive focus:border-destructive focus:ring-destructive' : '')
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
					class="absolute right-2.5 top-1/2 -translate-y-1/2 text-[11px] font-medium text-muted-foreground hover:text-foreground transition-colors"
					aria-label={showPassword ? 'Hide password' : 'Show password'}
				>
					{showPassword ? 'hide' : 'show'}
				</button>
			{/if}
		</div>
		{#if error}
			<span class="text-[12px] text-destructive">{error}</span>
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
				class="absolute right-2.5 top-1/2 -translate-y-1/2 text-[11px] font-medium text-muted-foreground hover:text-foreground transition-colors"
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				{showPassword ? 'hide' : 'show'}
			</button>
		{/if}
	</div>
	{#if error}
		<span class="text-[12px] text-destructive">{error}</span>
	{/if}
{/if}

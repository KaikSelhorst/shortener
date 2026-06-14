<script lang="ts">
	import { periods } from '$lib/analytics'

	interface Props {
		current: string
	}

	let { current }: Props = $props()
</script>

<div class="flex">
	{#each periods as p (p.value)}
		<a
			href="?period={p.value}"
			class="border border-transparent px-3 py-0.5 font-mono text-[10px] uppercase tracking-wider transition-colors
			{current === p.value
				? 'border-primary text-primary'
				: 'text-muted-foreground hover:text-foreground'}"
		>
			{p.label}
		</a>
	{/each}
</div>

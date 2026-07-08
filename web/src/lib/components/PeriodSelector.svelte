<script lang="ts">
	import { periods } from '$lib/analytics'
	interface Props { current: string }
	let { current }: Props = $props()
</script>

<div class="flex rounded-md border border-border overflow-hidden">
	{#each periods as p (p.value)}
		<a href="?period={p.value}"
			class="px-3 py-1 text-xs font-medium transition-colors border-r border-border last:border-r-0
				{current === p.value ? 'bg-secondary text-foreground' : 'text-muted-foreground hover:text-foreground hover:bg-secondary/50'}">
			{p.label}
		</a>
	{/each}
</div>

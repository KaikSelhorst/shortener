<script lang="ts">
	interface Item {
		label: string
		value: number
	}

	interface Props {
		items: Item[]
	}

	let { items }: Props = $props()

	const max = $derived(Math.max(...items.map((i) => i.value), 1))
	const total = $derived(items.reduce((s, i) => s + i.value, 0))
	const BAR_WIDTH = 24

	function bar(pct: number): string {
		const filled = Math.round(pct * BAR_WIDTH)
		return '█'.repeat(filled) + '░'.repeat(BAR_WIDTH - filled)
	}
</script>

<div class="flex flex-col gap-2">
	{#each items as item (item.label)}
		{@const pct = item.value / max}
		<div class="flex items-baseline gap-3 font-mono text-xs">
			<span class="w-20 shrink-0 truncate text-right text-muted-foreground" title={item.label}>
				{item.label}
			</span>
			<span class="shrink-0 text-[11px] tracking-tighter" style="color: color-mix(in srgb, var(--success) {Math.round(40 + pct * 60)}%, var(--muted-foreground))">
				{bar(pct)}
			</span>
			<div class="flex gap-2 shrink-0 text-muted-foreground">
				<span class="text-foreground">{item.value.toLocaleString()}</span>
				{#if total > 0}
					<span>({Math.round((item.value / total) * 100)}%)</span>
				{/if}
			</div>
		</div>
	{:else}
		<p class="py-6 text-center font-mono text-xs text-muted-foreground">-- no data --</p>
	{/each}
</div>

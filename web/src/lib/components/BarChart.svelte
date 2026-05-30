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
</script>

<div class="flex flex-col gap-2.5">
	{#each items as item (item.label)}
		<div class="flex items-center gap-3 text-sm">
			<span class="w-24 shrink-0 truncate text-right text-xs text-muted-foreground"
				>{item.label}</span
			>
			<div class="relative h-5 flex-1 overflow-hidden rounded-sm bg-muted">
				<div
					class="h-full rounded-sm bg-primary/70 transition-all duration-500"
					style="width: {(item.value / max) * 100}%"
				></div>
			</div>
			<div class="flex w-20 shrink-0 justify-end gap-1.5">
				<span class="font-mono text-xs text-foreground">{item.value.toLocaleString()}</span>
				{#if total > 0}
					<span class="text-xs text-muted-foreground"
						>({Math.round((item.value / total) * 100)}%)</span
					>
				{/if}
			</div>
		</div>
	{:else}
		<p class="py-4 text-center text-sm text-muted-foreground">No data</p>
	{/each}
</div>

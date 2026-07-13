<script lang="ts">
	import BarChart from '../charts/BarChart.svelte'

	interface Item { label: string; value: number }
	interface Props { devices: Item[]; browsers: Item[]; referrers: Item[] }

	let { devices, browsers, referrers }: Props = $props()
</script>

<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
	<div class="rounded-lg border border-border overflow-hidden">
		<div class="border-b border-border px-4 py-3">
			<h2 class="text-sm font-medium text-foreground">Devices</h2>
		</div>
		<div class="p-4"><BarChart items={devices} /></div>
	</div>
	<div class="rounded-lg border border-border overflow-hidden">
		<div class="border-b border-border px-4 py-3">
			<h2 class="text-sm font-medium text-foreground">Browsers</h2>
		</div>
		<div class="p-4"><BarChart items={browsers} /></div>
	</div>
	<div class="rounded-lg border border-border overflow-hidden lg:col-span-2">
		<div class="border-b border-border px-4 py-3">
			<h2 class="text-sm font-medium text-foreground">Referrers</h2>
		</div>
		<div class="p-4"><BarChart items={referrers} /></div>
	</div>
</div>

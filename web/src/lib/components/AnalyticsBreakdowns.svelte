<script lang="ts">
	import BarChart from './BarChart.svelte'

	interface Item {
		label: string
		value: number
	}

	interface Props {
		devices: Item[]
		browsers: Item[]
		referrers: Item[]
	}

	let { devices, browsers, referrers }: Props = $props()
</script>

<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
	<div class="tui-panel">
		<div class="tui-panel-header">▌ devices</div>
		<div class="p-4"><BarChart items={devices} /></div>
	</div>
	<div class="tui-panel">
		<div class="tui-panel-header">▌ browsers</div>
		<div class="p-4"><BarChart items={browsers} /></div>
	</div>
	<div class="tui-panel lg:col-span-2">
		<div class="tui-panel-header">▌ referrers</div>
		<div class="p-4"><BarChart items={referrers} /></div>
	</div>
</div>

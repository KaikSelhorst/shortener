<script lang="ts">
	import type { ClicksOverTime } from '$lib/types'

	interface Props {
		data: ClicksOverTime[]
		height?: number
	}

	let { data, height = 140 }: Props = $props()

	const W = 600
	const padTop = 8
	const padBottom = 24
	const padLeft = 4
	const padRight = 4

	const chartH = $derived(height - padTop - padBottom)
	const chartW = $derived(W - padLeft - padRight)

	const counts = $derived(data.map((d) => d.count))
	const maxVal = $derived(Math.max(...counts, 1))

	const pts = $derived(
		data.map((d, i) => ({
			x: padLeft + (data.length > 1 ? (i / (data.length - 1)) * chartW : chartW / 2),
			y: padTop + (1 - d.count / maxVal) * chartH,
			date: d.date,
			count: d.count,
		})),
	)

	const linePath = $derived(
		pts.length < 2 ? '' : `M ${pts.map((p) => `${p.x},${p.y}`).join(' L ')}`,
	)

	const areaPath = $derived(
		pts.length < 2
			? ''
			: `M ${pts[0].x},${padTop + chartH} L ${pts.map((p) => `${p.x},${p.y}`).join(' L ')} L ${pts[pts.length - 1].x},${padTop + chartH} Z`,
	)

	const labelStep = $derived(Math.ceil(data.length / 5))

	function fmtDate(iso: string) {
		return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
	}
</script>

{#if data.length === 0}
	<div
		class="flex items-center justify-center text-sm text-muted-foreground"
		style="height: {height}px"
	>
		No data for this period
	</div>
{:else}
	<svg
		viewBox="0 0 {W} {height}"
		class="w-full overflow-visible"
		preserveAspectRatio="none"
		aria-hidden="true"
	>
		<defs>
			<linearGradient id="lc-fill" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" style="stop-color: var(--primary); stop-opacity: 0.25" />
				<stop offset="100%" style="stop-color: var(--primary); stop-opacity: 0" />
			</linearGradient>
		</defs>

		<!-- area fill -->
		<path d={areaPath} fill="url(#lc-fill)" />

		<!-- line — currentColor inherits from the parent text-primary class -->
		<g class="text-primary">
			<path
				d={linePath}
				fill="none"
				stroke="currentColor"
				stroke-width="1.5"
				stroke-linejoin="round"
			/>
			{#each pts as p (p.date)}
				<circle cx={p.x} cy={p.y} r="2" fill="currentColor" />
			{/each}
		</g>

		<!-- x-axis labels — style= so var() is resolved as a CSS property -->
		{#each pts as p, i (p.date)}
			{#if i % labelStep === 0 || i === pts.length - 1}
				<text
					x={p.x}
					y={height - 4}
					text-anchor="middle"
					font-size="9"
					style="fill: var(--muted-foreground)"
				>
					{fmtDate(p.date)}
				</text>
			{/if}
		{/each}
	</svg>
{/if}

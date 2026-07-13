<script lang="ts">
	import type { ClicksOverTime } from '$lib/types'

	interface Props {
		data: ClicksOverTime[]
		height?: number
	}

	let { data, height = 240 }: Props = $props()

	const W = 600
	const padTop = 10
	const padBottom = 26
	const padLeft = 4
	const padRight = 4

	const chartH = $derived(height - padTop - padBottom)
	const chartW = $derived(W - padLeft - padRight)

	const maxVal = $derived(data.reduce((m, d) => Math.max(m, d.count), 1))

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

	const dateFmt = new Intl.DateTimeFormat(undefined, { month: 'short', day: 'numeric' })
	function fmtDate(iso: string) {
		return dateFmt.format(new Date(iso))
	}

	// tooltip state
	type Pt = (typeof pts)[number]
	let hovered = $state<Pt | null>(null)

	const TW = 110  // tooltip width in SVG units
	const TH = 36   // tooltip height in SVG units
	const PAD = 8   // gap between point and tooltip

	const tooltipX = $derived(
		hovered
			? hovered.x + TW + PAD > W
				? hovered.x - TW - PAD
				: hovered.x + PAD
			: 0,
	)
	const tooltipY = $derived(
		hovered
			? hovered.y - TH - PAD < padTop
				? hovered.y + PAD
				: hovered.y - TH - PAD
			: 0,
	)
</script>

{#if data.length === 0}
	<div class="flex items-center justify-center font-mono text-xs text-muted-foreground" style="height:{height}px">
		-- no data for this period --
	</div>
{:else}
	<svg
		viewBox="0 0 {W} {height}"
		style="height: {height}px"
		class="w-full overflow-visible"
		preserveAspectRatio="none"
		aria-hidden="true"
		onmouseleave={() => (hovered = null)}
	>
		<defs>
			<linearGradient id="line-chart-area" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%"   stop-color="var(--success)" stop-opacity="0.18" />
				<stop offset="100%" stop-color="var(--success)" stop-opacity="0"    />
			</linearGradient>
		</defs>

		<!-- horizontal grid lines -->
		{#each [0, 0.25, 0.5, 0.75, 1] as frac}
			<line
				x1={padLeft} y1={padTop + frac * chartH}
				x2={W - padRight} y2={padTop + frac * chartH}
				stroke="var(--border)" stroke-width="1"
			/>
		{/each}

		<path d={areaPath} fill="url(#line-chart-area)" />

		<path d={linePath} fill="none" stroke="var(--success)" stroke-width="1.5" stroke-linejoin="round" />

		<!-- invisible wider hit area + visible dot per point -->
		{#each pts as p (p.date)}
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<g
				onmouseenter={() => (hovered = p)}
				style="cursor: crosshair"
			>
				<!-- large invisible hit area -->
				<circle cx={p.x} cy={p.y} r="10" fill="transparent" />
				<!-- visible dot, grows on hover -->
				<circle
					cx={p.x} cy={p.y}
					r={hovered?.date === p.date ? 4 : 2.5}
					fill="var(--success)"
					stroke={hovered?.date === p.date ? 'var(--background)' : 'none'}
					stroke-width="1.5"
					style="transition: r 0.1s"
				/>
			</g>
		{/each}

		<!-- x-axis labels -->
		{#each pts as p, i (p.date)}
			{#if i % labelStep === 0 || i === pts.length - 1}
				<text
					x={p.x} y={height - 5}
					text-anchor={i === 0 ? 'start' : i === pts.length - 1 ? 'end' : 'middle'}
					font-size="9"
					fill={hovered?.date === p.date ? 'var(--info)' : 'var(--muted-foreground)'}
				>
					{fmtDate(p.date)}
				</text>
			{/if}
		{/each}

		<!-- tooltip -->
		{#if hovered}
			<g pointer-events="none">
				<rect
					x={tooltipX} y={tooltipY}
					width={TW} height={TH}
					fill="var(--card)"
					stroke="var(--border)"
					stroke-width="1"
				/>
				<!-- top line: date -->
				<text
					x={tooltipX + 8} y={tooltipY + 13}
					font-size="9"
					fill="var(--muted-foreground)"
				>
					{fmtDate(hovered.date)}
				</text>
				<!-- bottom line: count -->
				<text
					x={tooltipX + 8} y={tooltipY + 27}
					font-size="10"
					fill="var(--success)"
					font-weight="700"
				>
					{hovered.count.toLocaleString()} clicks
				</text>
			</g>
		{/if}
	</svg>
{/if}

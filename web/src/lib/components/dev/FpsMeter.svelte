<script lang="ts">
	import { dev } from '$app/environment'

	const HISTORY_SIZE = 10

	let fps = $state(0)
	let minFps = $state(0)
	let avgFps = $state(0)

	let heapUsed = $state<number | null>(null)
	let heapLimit = $state<number | null>(null)

	let lcp = $state<number | null>(null)
	let longTasks = $state(0)

	let open = $state(false)
	let container = $state<HTMLDivElement>()

	type PerfWithMemory = Performance & {
		memory?: { usedJSHeapSize: number; jsHeapSizeLimit: number }
	}

	function formatMB(bytes: number) {
		return (bytes / 1024 / 1024).toFixed(1) + ' MB'
	}

	function formatMs(ms: number) {
		return ms < 1000 ? Math.round(ms) + ' ms' : (ms / 1000).toFixed(2) + ' s'
	}

	function fpsColor(v: number) {
		return v < 30 ? 'text-red-400' : v < 50 ? 'text-yellow-400' : 'text-green-400'
	}

	function lcpColor(v: number | null) {
		if (v === null) return ''
		if (v > 4000) return 'text-red-400'
		if (v > 2500) return 'text-yellow-400'
		return 'text-green-400'
	}

	function longTasksColor(v: number) {
		return v > 5 ? 'text-red-400' : v > 0 ? 'text-yellow-400' : ''
	}

	$effect(() => {
		if (!dev) return

		const history: number[] = []
		let frames = 0
		let lastTime = performance.now()
		let raf: number

		function tick(now: number) {
			frames++
			if (now - lastTime >= 1000) {
				const current = Math.round((frames * 1000) / (now - lastTime))
				fps = current
				history.push(current)
				if (history.length > HISTORY_SIZE) history.shift()
				minFps = Math.min(...history)
				avgFps = Math.round(history.reduce((a, b) => a + b, 0) / history.length)
				frames = 0
				lastTime = now

				const mem = (performance as PerfWithMemory).memory
				if (mem) {
					heapUsed = mem.usedJSHeapSize
					heapLimit = mem.jsHeapSizeLimit
				}
			}
			raf = requestAnimationFrame(tick)
		}

		raf = requestAnimationFrame(tick)
		return () => cancelAnimationFrame(raf)
	})

	$effect(() => {
		if (!dev) return
		if (!('PerformanceObserver' in window)) return

		const supported = PerformanceObserver.supportedEntryTypes
		const observers: PerformanceObserver[] = []

		if (supported.includes('largest-contentful-paint')) {
			const obs = new PerformanceObserver((list) => {
				const entries = list.getEntries()
				lcp = entries[entries.length - 1].startTime
			})
			obs.observe({ type: 'largest-contentful-paint', buffered: true })
			observers.push(obs)
		}

		if (supported.includes('longtask')) {
			const obs = new PerformanceObserver((list) => {
				longTasks += list.getEntries().length
			})
			obs.observe({ type: 'longtask' })
			observers.push(obs)
		}

		return () => observers.forEach((o) => o.disconnect())
	})
</script>

<svelte:window
	onclick={(e) => {
		if (open && container && !container.contains(e.target as Node)) open = false
	}}
/>

{#if dev}
	<div bind:this={container} class="fixed bottom-3 right-3 z-50">
		{#if open}
			<div
				class="absolute bottom-full right-0 mb-2 w-52 rounded border border-white/10 bg-black/90 p-3 font-mono text-xs text-white backdrop-blur-sm"
			>
				<p class="mb-2 text-[10px] uppercase tracking-wider text-white/40">performance</p>
				<dl class="flex flex-col gap-1.5">
					<div class="flex justify-between">
						<dt class="text-white/50">fps now</dt>
						<dd class={fpsColor(fps)}>{fps}</dd>
					</div>
					<div class="flex justify-between">
						<dt class="text-white/50">fps avg (10s)</dt>
						<dd>{avgFps}</dd>
					</div>
					<div class="flex justify-between">
						<dt class="text-white/50">fps min (10s)</dt>
						<dd class={minFps < 30 ? 'text-red-400' : minFps < 50 ? 'text-yellow-400' : ''}>{minFps}</dd>
					</div>

					<div class="my-0.5 border-t border-white/10"></div>

					<div class="flex justify-between">
						<dt class="text-white/50">heap used</dt>
						<dd>{heapUsed !== null ? formatMB(heapUsed) : '—'}</dd>
					</div>
					<div class="flex justify-between">
						<dt class="text-white/50">heap limit</dt>
						<dd>{heapUsed !== null ? formatMB(heapLimit!) : '—'}</dd>
					</div>

					<div class="my-0.5 border-t border-white/10"></div>

					<div class="flex justify-between">
						<dt class="text-white/50">lcp</dt>
						<dd class={lcpColor(lcp)}>{lcp !== null ? formatMs(lcp) : '—'}</dd>
					</div>
					<div class="flex justify-between">
						<dt class="text-white/50">long tasks</dt>
						<dd class={longTasksColor(longTasks)}>{longTasks}</dd>
					</div>
				</dl>
			</div>
		{/if}

		<button
			onclick={() => (open = !open)}
			class="min-w-[3.5rem] rounded bg-black/80 px-2 py-1 text-center font-mono text-xs tabular-nums text-white backdrop-blur-sm transition-colors hover:bg-black/90 {open ? 'ring-1 ring-white/20' : ''}"
		>
			{fps} fps
		</button>
	</div>
{/if}

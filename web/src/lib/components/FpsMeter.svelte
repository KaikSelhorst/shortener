<script lang="ts">
	import { dev } from '$app/environment'

	let fps = $state(0)

	$effect(() => {
		if (!dev) return

		let frames = 0
		let lastTime = performance.now()
		let raf: number

		function tick(now: number) {
			frames++
			if (now - lastTime >= 1000) {
				fps = Math.round((frames * 1000) / (now - lastTime))
				frames = 0
				lastTime = now
			}
			raf = requestAnimationFrame(tick)
		}

		raf = requestAnimationFrame(tick)

		return () => cancelAnimationFrame(raf)
	})
</script>

{#if dev}
	<div
		class="fixed bottom-3 right-3 z-50 min-w-[3.5rem] rounded bg-black/80 px-2 py-1 text-center font-mono text-xs tabular-nums text-white backdrop-blur-sm"
	>
		{fps} fps
	</div>
{/if}

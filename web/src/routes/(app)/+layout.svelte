<script lang="ts">
	import UserMenu from '$lib/components/UserMenu.svelte'

	let { children } = $props()

	let time = $state('')

	$effect(() => {
		function tick() {
			const now = new Date()
			time =
				String(now.getHours()).padStart(2, '0') +
				':' +
				String(now.getMinutes()).padStart(2, '0') +
				':' +
				String(now.getSeconds()).padStart(2, '0')
		}

		let id: ReturnType<typeof setInterval> | undefined

		function start() {
			tick()
			id = setInterval(tick, 1000)
		}

		function stop() {
			clearInterval(id)
			id = undefined
		}

		function onVisibilityChange() {
			document.hidden ? stop() : start()
		}

		document.addEventListener('visibilitychange', onVisibilityChange)
		start()

		return () => {
			stop()
			document.removeEventListener('visibilitychange', onVisibilityChange)
		}
	})
</script>

<div class="flex min-h-screen flex-col bg-background">
	<header class="border-b border-border bg-card">
		<div class="flex items-center px-4 py-2.5">
			<div class="ml-auto flex items-center gap-4">
				<span class="font-mono text-[11px] tabular-nums text-accent">{time}</span>
				<div class="h-3 w-px bg-border"></div>
				<UserMenu />
			</div>
		</div>
	</header>

	<main class="flex-1 px-4 py-6">
		{@render children()}
	</main>

	<footer class="border-t border-border px-4 py-1.5">
		<div class="flex items-center justify-between">
			<span class="font-mono text-[9px] uppercase tracking-widest text-muted-foreground">
				shortener — url management system
			</span>
			<span class="font-mono text-[9px] text-muted-foreground">sys:ok</span>
		</div>
	</footer>
</div>

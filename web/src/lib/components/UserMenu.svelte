<script lang="ts">
	let open = $state(false)

	function clickOutside(node: HTMLElement) {
		function handle(e: MouseEvent) {
			if (!node.contains(e.target as Node)) open = false
		}
		document.addEventListener('click', handle, true)
		return { destroy() { document.removeEventListener('click', handle, true) } }
	}
</script>

<div class="relative" use:clickOutside>
	<button
		onclick={() => (open = !open)}
		aria-expanded={open}
		class="font-mono text-[10px] uppercase tracking-widest text-muted-foreground hover:text-primary transition-colors"
	>
		[menu {open ? '▲' : '▼'}]
	</button>

	{#if open}
		<div class="absolute right-0 top-full z-50 mt-1 w-40 border border-border bg-background">
			<a
				href="/settings"
				onclick={() => (open = false)}
				class="flex w-full items-center gap-2 px-3 py-2 font-mono text-xs text-foreground hover:bg-card hover:text-primary transition-colors"
			>
				▶ settings
			</a>
			<a
				href="/docs"
				onclick={() => (open = false)}
				class="flex w-full items-center gap-2 px-3 py-2 font-mono text-xs text-foreground hover:bg-card hover:text-primary transition-colors"
			>
				▶ api docs
			</a>
			<div class="border-t border-border"></div>
			<form method="POST" action="/api/auth/logout">
				<button
					type="submit"
					class="flex w-full items-center gap-2 px-3 py-2 font-mono text-xs text-destructive hover:bg-card transition-colors"
				>
					▶ sign out
				</button>
			</form>
		</div>
	{/if}
</div>

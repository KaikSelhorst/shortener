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
		class="flex items-center gap-1.5 rounded-md px-2 py-1 text-sm text-muted-foreground hover:text-foreground hover:bg-secondary transition-colors"
	>
		<span>Account</span>
		<svg width="12" height="12" viewBox="0 0 12 12" fill="none" class="transition-transform {open ? 'rotate-180' : ''}">
			<path d="M2 4l4 4 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
		</svg>
	</button>

	{#if open}
		<div class="absolute right-0 top-full z-50 mt-1.5 w-44 rounded-lg border border-border bg-popover shadow-lg overflow-hidden">
			<div class="p-1">
				<a
					href="/settings"
					onclick={() => (open = false)}
					class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-foreground hover:bg-secondary transition-colors"
				>
					Settings
				</a>
				<a
					href="/docs"
					onclick={() => (open = false)}
					class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-foreground hover:bg-secondary transition-colors"
				>
					API Docs
				</a>
			</div>
			<div class="border-t border-border p-1">
				<form method="POST" action="/api/auth/logout">
					<button
						type="submit"
						class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-destructive hover:bg-destructive/10 transition-colors"
					>
						Sign out
					</button>
				</form>
			</div>
		</div>
	{/if}
</div>

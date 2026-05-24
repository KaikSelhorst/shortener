<script lang="ts">
	let open = $state(false)
	let dark = $state(false)

	$effect(() => {
		dark = document.documentElement.classList.contains('dark')
	})

	function toggleTheme() {
		dark = !dark
		document.documentElement.classList.toggle('dark', dark)
		localStorage.setItem('theme', dark ? 'dark' : 'light')
	}

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
		aria-label="User menu"
		class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			aria-hidden="true"
		>
			<circle cx="12" cy="8" r="4" />
			<path d="M4 20c0-4 3.6-7 8-7s8 3 8 7" />
		</svg>
	</button>

	{#if open}
		<div
			class="absolute right-0 top-full z-50 mt-1 w-44 rounded-md border border-border bg-card shadow-sm"
		>
			<button
				onclick={toggleTheme}
				class="flex w-full items-center gap-2.5 px-3 py-2 text-sm text-foreground transition-colors hover:bg-muted"
			>
				{#if dark}
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
						<circle cx="12" cy="12" r="4" />
						<path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
					</svg>
					Light mode
				{:else}
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
						<path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z" />
					</svg>
					Dark mode
				{/if}
			</button>

			<div class="border-t border-border"></div>

			<a
				href="/settings"
				class="flex w-full items-center gap-2.5 px-3 py-2 text-sm text-foreground transition-colors hover:bg-muted"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
					<path d="M12 20h9"/><path d="M16.376 3.622a1 1 0 0 1 3.002 3.002L7.368 18.635a2 2 0 0 1-.855.506l-2.872.838a.5.5 0 0 1-.62-.62l.838-2.872a2 2 0 0 1 .506-.854z"/>
				</svg>
				Settings
			</a>

			<div class="border-t border-border"></div>

			<form method="POST" action="/api/auth/logout">
				<button
					type="submit"
					class="flex w-full items-center gap-2.5 px-3 py-2 text-sm text-destructive transition-colors hover:bg-muted"
				>
					<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
						<polyline points="16 17 21 12 16 7" />
						<line x1="21" y1="12" x2="9" y2="12" />
					</svg>
					Sign out
				</button>
			</form>
		</div>
	{/if}
</div>

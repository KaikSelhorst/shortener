<script lang="ts">
	let dark = $state(false)

	$effect(() => {
		dark = document.documentElement.classList.contains('dark')
	})

	function toggle() {
		dark = !dark
		document.documentElement.classList.toggle('dark', dark)
		localStorage.setItem('theme', dark ? 'dark' : 'light')
	}
</script>

<button
	onclick={toggle}
	aria-label={dark ? 'Switch to light mode' : 'Switch to dark mode'}
	class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
>
	{#if dark}
		<!-- sun -->
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
			<circle cx="12" cy="12" r="4" />
			<path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41" />
		</svg>
	{:else}
		<!-- moon -->
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
			<path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z" />
		</svg>
	{/if}
</button>

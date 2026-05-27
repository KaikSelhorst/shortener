<script lang="ts">
	import UserMenu from '$lib/components/UserMenu.svelte'

	let { children } = $props()
</script>

<div class="flex min-h-screen flex-col">
	<header class="border-b border-border bg-background">
		<div class="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
			<a href="/" class="text-sm font-semibold text-foreground hover:text-muted-foreground transition-colors">
				Shortener
			</a>
			<UserMenu />
		</div>
	</header>

	<main class="mx-auto w-full max-w-5xl flex-1 px-4 py-8">
		{@render children()}
	</main>
</div>

<script lang="ts">
	let { children } = $props()
</script>

<main class="flex min-h-screen items-center justify-center bg-background">
	<div class="w-full max-w-sm rounded-md border border-border bg-card p-8">
		{@render children()}
	</div>
</main>

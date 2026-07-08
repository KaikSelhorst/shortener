<script lang="ts">
	let { children } = $props()
</script>

<div class="flex min-h-screen flex-col items-center justify-center bg-background px-4">
	<div class="w-full max-w-sm">
		<div class="mb-8">
			<a href="/" class="text-[13px] font-semibold text-foreground tracking-tight">Shortener</a>
		</div>
		{@render children()}
	</div>
</div>

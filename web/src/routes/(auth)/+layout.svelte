<script lang="ts">
	let { children } = $props()
</script>

<div class="flex min-h-screen flex-col items-center justify-center bg-background px-4">
	<div class="w-full max-w-sm">
		{@render children()}
	</div>
</div>

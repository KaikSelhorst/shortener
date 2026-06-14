<script lang="ts">
	let { children } = $props()
</script>

<div class="flex min-h-screen items-center justify-center bg-background px-4">
	<div class="w-full max-w-sm">
		<div class="border border-border bg-card p-6">
			{@render children()}
		</div>
		<div class="border-x border-b border-border bg-background px-3 py-1">
			<div class="flex items-center justify-between">
				<span class="font-mono text-[9px] uppercase tracking-widest text-muted-foreground">
					secure connection
				</span>
				<span class="font-mono text-[9px] text-primary">●</span>
			</div>
		</div>
	</div>
</div>

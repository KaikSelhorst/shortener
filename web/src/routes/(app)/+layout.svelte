<script lang="ts">
	import { page } from '$app/stores'
	let { children } = $props()

	const navItems = [
		{ href: '/dashboard', label: 'Dashboard' },
		{ href: '/settings', label: 'Settings' },
		{ href: '/docs', label: 'API Docs' },
	]

	function isActive(href: string, pathname: string): boolean {
		if (href === '/dashboard') return pathname === '/dashboard'
		return pathname === href || pathname.startsWith(href + '/')
	}
</script>

<div class="flex h-screen overflow-hidden bg-background">
	<aside class="w-44 shrink-0 border-r border-border flex flex-col">
		<div class="h-11 flex items-center px-5 border-b border-border shrink-0">
			<a href="/dashboard" class="text-[13px] font-semibold text-foreground tracking-tight">Shortener</a>
		</div>

		<nav class="flex-1 overflow-y-auto py-2">
			{#each navItems as item}
				{@const active = isActive(item.href, $page.url.pathname)}
				<a
					href={item.href}
					class="flex items-center px-5 py-1.5 text-sm transition-colors border-l-2
						{active
							? 'border-foreground text-foreground'
							: 'border-transparent text-muted-foreground hover:text-foreground'}"
				>
					{item.label}
				</a>
			{/each}
		</nav>

		<div class="border-t border-border shrink-0">
			<form method="POST" action="/api/auth/logout">
				<button
					type="submit"
					class="w-full flex items-center px-5 py-3 text-sm text-muted-foreground hover:text-destructive transition-colors"
				>
					Sign out
				</button>
			</form>
		</div>
	</aside>

	<div class="flex-1 min-h-0 overflow-auto">
		{@render children()}
	</div>
</div>

<script lang="ts">
	import { page } from '$app/stores'
	import type { LayoutData } from './$types'
	import { ProjectSelector, Breadcrumbs, CreateProjectDialog } from '$lib/components/layout'
	import { Button } from '$lib/components/ui'

	let { data, children }: { data: LayoutData; children: import('svelte').Snippet } = $props()

	const actualSlug = $derived($page.params.slug ?? '')
	// Falls back to the first project so the project tabs always point somewhere.
	const activeSlug = $derived(actualSlug || (data.projects[0]?.slug ?? ''))
	const noProjects = $derived(data.projects.length === 0)

	let manualCreateOpen = $state(false)
	$effect(() => {
		if (noProjects) manualCreateOpen = true
	})

	const appLinks = [
		{ href: '/settings', label: 'Settings' },
		{ href: '/docs', label: 'API Docs' },
	]

	const projectLinks = $derived(
		activeSlug
			? [
					{ key: 'links', href: `/${activeSlug}`, label: 'Links' },
					{ key: 'analytics', href: `/${activeSlug}/analytics`, label: 'Analytics' },
					{ key: 'webhooks', href: `/${activeSlug}/webhooks`, label: 'Webhooks' },
				]
			: [],
	)

	const activeProjectSection = $derived.by(() => {
		if (!actualSlug) return null
		const rest = $page.url.pathname.slice(`/${actualSlug}`.length)
		if (rest.startsWith('/webhooks')) return 'webhooks'
		if (rest === '/analytics' || rest.endsWith('/analytics')) return 'analytics'
		return 'links'
	})

	const linkClass = (active: boolean) =>
		`mx-2 px-3 py-1.5 rounded text-sm transition-colors ${
			active ? 'bg-secondary text-foreground' : 'text-muted-foreground hover:text-foreground hover:bg-secondary/50'
		}`
</script>

<div class="flex h-screen overflow-hidden bg-background">
	<aside class="w-56 shrink-0 border-r border-border flex flex-col">
		<div class="px-3 py-3 border-b border-border flex flex-col gap-2">
			<ProjectSelector projects={data.projects} />
			<Button variant="outline" size="sm" class="w-full" onclick={() => (manualCreateOpen = true)}>
				+ New project
			</Button>
		</div>

		<nav class="flex-1 overflow-y-auto py-2 flex flex-col gap-0.5">
			{#if projectLinks.length > 0}
				{#each projectLinks as link (link.key)}
					<a href={link.href} class={linkClass(activeProjectSection === link.key)}>{link.label}</a>
				{/each}
				<div class="mt-4 mb-1 px-5 text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">
					App
				</div>
			{/if}

			{#each appLinks as link (link.href)}
				<a href={link.href} class={linkClass($page.url.pathname === link.href)}>{link.label}</a>
			{/each}
		</nav>

		<CreateProjectDialog bind:open={manualCreateOpen} dismissable={!noProjects} />

		<div class="border-t border-border">
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

	<div class="flex-1 min-w-0 flex flex-col">
		<div class="h-11 shrink-0 border-b border-border px-4 flex items-center">
			<Breadcrumbs />
		</div>
		<div class="flex-1 min-h-0 overflow-y-auto">
			{@render children()}
		</div>
	</div>
</div>

<script lang="ts">
	import { page } from '$app/stores'

	interface Crumb {
		label: string
		href?: string
	}

	const crumbs = $derived.by((): Crumb[] => {
		const id = $page.route.id ?? ''
		const params = $page.params
		const data = $page.data as { webhook?: { url: string } }

		let list: Crumb[] = []

		if (id === '/(app)/settings') {
			list = [{ label: 'Settings' }]
		} else if (id === '/(app)/docs') {
			list = [{ label: 'API Docs' }]
		} else if (id.startsWith('/(app)/[slug]')) {
			const slug = params.slug ?? ''
			list = [{ label: slug, href: `/${slug}` }]

			if (id === '/(app)/[slug]/analytics') {
				list.push({ label: 'Analytics' })
			} else if (id === '/(app)/[slug]/[code]/analytics') {
				list.push({ label: params.code ?? '' }, { label: 'Analytics' })
			} else if (id === '/(app)/[slug]/webhooks') {
				list.push({ label: 'Webhooks' })
			} else if (id === '/(app)/[slug]/webhooks/[id]') {
				list.push(
					{ label: 'Webhooks', href: `/${slug}/webhooks` },
					{ label: data.webhook?.url ?? params.id ?? '' },
				)
			}
		}

		if (list.length > 0) list[list.length - 1] = { label: list[list.length - 1].label }
		return list
	})
</script>

<div class="flex items-center gap-1.5 text-sm min-w-0">
	{#each crumbs as crumb, i (i)}
		{#if i > 0}<span class="text-border shrink-0">/</span>{/if}
		{#if crumb.href}
			<a href={crumb.href} class="text-muted-foreground hover:text-foreground transition-colors truncate shrink-0">
				{crumb.label}
			</a>
		{:else}
			<span class="font-medium text-foreground truncate">{crumb.label}</span>
		{/if}
	{/each}
</div>

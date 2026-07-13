<script lang="ts">
	import { page } from '$app/stores'
	import { goto } from '$app/navigation'
	import type { Project } from '$lib/types'

	interface Props {
		projects: Project[]
	}

	let { projects }: Props = $props()

	// The route's actual slug (empty outside a project, e.g. /settings).
	const actualSlug = $derived($page.params.slug ?? '')
	// What the select shows: falls back to the first project so it's never stuck on a placeholder.
	const displaySlug = $derived(actualSlug || (projects[0]?.slug ?? ''))

	function sectionFor(pathname: string, slug: string): string {
		if (!slug) return ''
		const rest = pathname.slice(`/${slug}`.length)
		if (rest.startsWith('/webhooks')) return '/webhooks'
		if (rest === '/analytics') return '/analytics'
		return ''
	}

	function handleChange(e: Event) {
		const newSlug = (e.currentTarget as HTMLSelectElement).value
		if (!newSlug || newSlug === actualSlug) return
		const section = sectionFor($page.url.pathname, actualSlug)
		goto(`/${newSlug}${section}`)
	}
</script>

<select
	value={displaySlug}
	onchange={handleChange}
	class="h-8 w-full text-[13px]"
	aria-label="Select project"
>
	{#if projects.length === 0}
		<option value="" disabled selected>No projects</option>
	{/if}
	{#each projects as project (project.id)}
		<option value={project.slug}>{project.name}</option>
	{/each}
</select>

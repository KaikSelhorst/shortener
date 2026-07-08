<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import { Button, Input, Dialog } from '$lib'
	import { useSSE } from '$lib/sse.svelte'
	import { formatDate } from '$lib/format'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	let pendingSlug = $state<string | null>(null)
	let confirmOpen = $state(false)
	$effect(() => { if (!confirmOpen) pendingSlug = null })

	let sessionClicks = $state<Record<number, number>>({})
	const sse = useSSE(() => '/api/stream', (e) => {
		const evt = JSON.parse(e.data) as { project_id: number }
		sessionClicks[evt.project_id] = (sessionClicks[evt.project_id] ?? 0) + 1
	})
</script>

<div class="sticky top-0 z-10 bg-background border-b border-border px-4 h-11 flex items-center justify-between">
	<div class="flex items-center gap-2">
		<span class="text-sm font-medium text-foreground">Projects</span>
		<span class="text-xs text-muted-foreground">{data.projects.length}</span>
	</div>
	{#if sse.connected}
		<span class="flex items-center gap-1.5 text-xs text-success">
			<span class="w-1.5 h-1.5 rounded-full bg-success"></span>Live
		</span>
	{/if}
</div>

<div class="border-b border-border px-4 py-2.5">
	<form method="POST" action="?/create" class="flex items-center gap-2">
		<Input name="name" type="text" placeholder="New project name…" class="max-w-xs h-8 text-sm" required />
		<Button type="submit" size="sm">Create</Button>
	</form>
	{#if form?.error}
		<p class="mt-2 text-xs text-destructive">{form.error}</p>
	{/if}
</div>

{#if data.projects.length === 0}
	<div class="flex items-center justify-center py-20 text-sm text-muted-foreground">
		No projects yet — create one above.
	</div>
{:else}
	<table class="w-full text-sm">
		<thead class="border-b border-border">
			<tr>
				<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Name</th>
				<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Slug</th>
				<th class="px-4 py-1.5 text-left text-[11px] font-semibold text-muted-foreground uppercase tracking-wider">Created</th>
				<th class="px-4 py-1.5 text-right text-[11px] font-semibold text-muted-foreground uppercase tracking-wider"></th>
			</tr>
		</thead>
		<tbody>
			{#each data.projects as project (project.id)}
				<tr class="group hover:bg-secondary/40 transition-colors">
					<td class="px-4 py-1.5">
						<div class="flex items-center gap-2">
							<a href="/{project.slug}" class="text-sm font-medium text-foreground hover:underline underline-offset-4">
								{project.name}
							</a>
							{#if sessionClicks[project.id]}
								<span class="text-xs text-success">+{sessionClicks[project.id]}</span>
							{/if}
						</div>
					</td>
					<td class="px-4 py-1.5 text-sm text-muted-foreground">{project.slug}</td>
					<td class="px-4 py-1.5 text-sm text-muted-foreground">{formatDate(project.created_at)}</td>
					<td class="px-4 py-1.5 text-right">
						<div class="flex items-center justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
							<Button variant="ghost" size="sm" href="/{project.slug}/analytics">Analytics</Button>
							<Button variant="ghost-destructive" size="sm" onclick={() => { pendingSlug = project.slug; confirmOpen = true }}>Delete</Button>
						</div>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
{/if}

<Dialog bind:open={confirmOpen} title="Delete project" description="This will permanently delete the project and all its links. This action cannot be undone.">
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>Cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="slug" value={pendingSlug} />
			<Button type="submit" variant="destructive">Delete</Button>
		</form>
	{/snippet}
</Dialog>

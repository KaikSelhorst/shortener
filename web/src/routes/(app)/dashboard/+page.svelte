<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import { Button, Input, Dialog } from '$lib'
	import { useSSE } from '$lib/sse.svelte'
	import { formatDate } from '$lib/format'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	let pendingSlug = $state<string | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingSlug = null
	})

	let sessionClicks = $state<Record<number, number>>({})

	const sse = useSSE(() => '/api/stream', (e) => {
		const evt = JSON.parse(e.data) as { project_id: number }
		sessionClicks[evt.project_id] = (sessionClicks[evt.project_id] ?? 0) + 1
	})
</script>

<div class="mx-auto max-w-6xl">
	<div class="tui-panel">
		<div class="tui-panel-header justify-between">
			<span>▌ projects</span>
			<div class="flex items-center gap-3">
				{#if sse.connected}
					<span class="font-mono text-[10px] uppercase tracking-wider text-tui-green">● live</span>
				{/if}
				<span class="text-muted-foreground">{data.projects.length} total</span>
			</div>
		</div>

		<div class="border-b border-border bg-background px-4 py-3">
			<form method="POST" action="?/create" class="flex items-end gap-3">
				<div class="flex-1">
					<Input name="name" type="text" placeholder="new project name..." required />
				</div>
				<Button type="submit" size="sm">+ create</Button>
			</form>
			{#if form?.error}
				<p class="mt-2 font-mono text-xs text-destructive">{form.error}</p>
			{/if}
		</div>

		{#if data.projects.length === 0}
			<div class="px-4 py-14 text-center font-mono text-xs text-muted-foreground">
				-- no projects yet --
			</div>
		{:else}
			<table class="w-full">
				<thead class="border-b border-border">
					<tr>
						<th class="px-4 py-2.5 text-left font-normal tui-label">Name</th>
						<th class="px-4 py-2.5 text-left font-normal tui-label">Slug</th>
						<th class="px-4 py-2.5 text-left font-normal tui-label">Created</th>
						<th class="px-4 py-2.5 text-right font-normal tui-label"></th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border">
					{#each data.projects as project (project.id)}
						<tr class="group transition-colors hover:bg-primary/5">
							<td class="px-4 py-3">
								<div class="flex items-center gap-2">
									<a href="/{project.slug}" class="font-mono text-xs text-primary hover:underline">
										{project.name}
									</a>
									{#if sessionClicks[project.id]}
										<span class="font-mono text-[10px] text-tui-green">
											+{sessionClicks[project.id]}
										</span>
									{/if}
								</div>
							</td>
							<td class="px-4 py-3">
								<span class="font-mono text-xs text-muted-foreground">{project.slug}</span>
							</td>
							<td class="px-4 py-3">
								<span class="font-mono text-xs text-muted-foreground">
									{formatDate(project.created_at)}
								</span>
							</td>
							<td class="px-4 py-3 text-right">
								<div
									class="flex items-center justify-end gap-3 opacity-0 transition-opacity group-hover:opacity-100"
								>
									<a
										href="/{project.slug}/analytics"
										class="font-mono text-[10px] uppercase tracking-wider text-accent hover:underline"
									>
										analytics
									</a>
									<span class="text-border">|</span>
									<button
										onclick={() => {
											pendingSlug = project.slug
											confirmOpen = true
										}}
										class="font-mono text-[10px] uppercase tracking-wider text-destructive hover:underline"
									>
										delete
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>

<Dialog
	bind:open={confirmOpen}
	title="Delete project"
	description="This will permanently delete the project and all its links. This action cannot be undone."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (confirmOpen = false)}>cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="slug" value={pendingSlug} />
			<Button type="submit" variant="destructive">confirm delete</Button>
		</form>
	{/snippet}
</Dialog>

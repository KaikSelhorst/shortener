<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import { Button, Input, Dialog } from '$lib'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	let pendingSlug = $state<string | null>(null)
	let confirmOpen = $derived(pendingSlug !== null)
</script>

<h1 class="text-lg font-semibold text-foreground">Projects</h1>

{#if form?.error}
	<p class="mt-4 rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">{form.error}</p>
{/if}

<form method="POST" action="?/create" class="mt-4 grid grid-cols-[1fr_auto] gap-2">
	<Input name="name" type="text" placeholder="New project name" required />
	<Button type="submit">Create</Button>
</form>

{#if data.projects.length === 0}
	<p class="mt-10 text-center text-sm text-muted-foreground">No projects yet. Create one above.</p>
{:else}
	<ul class="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
		{#each data.projects as project (project.id)}
			<li class="group flex flex-col justify-between rounded-md border-2 border-border bg-card p-5">
				<div>
					<p class="text-xs text-muted-foreground">
						{new Date(project.created_at).toLocaleDateString()}
					</p>
					<a
						href="/{project.slug}"
						class="mt-1 block text-base font-semibold text-card-foreground hover:underline"
					>
						{project.name}
					</a>
					<p class="mt-0.5 font-mono text-xs text-muted-foreground">{project.slug}</p>
				</div>

				<div class="mt-4 flex items-center justify-between border-t-2 border-border pt-3">
					<a href="/{project.slug}" class="text-xs font-medium text-foreground hover:underline">
						View links →
					</a>
					<Button
						variant="ghost-destructive"
						size="sm"
						onclick={() => (pendingSlug = project.slug)}
					>
						Delete
					</Button>
				</div>
			</li>
		{/each}
	</ul>
{/if}

<Dialog
	bind:open={confirmOpen}
	title="Delete project"
	description="This will permanently delete the project and all its links. This action cannot be undone."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (pendingSlug = null)}>Cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="slug" value={pendingSlug} />
			<Button type="submit" variant="destructive">Delete</Button>
		</form>
	{/snippet}
</Dialog>

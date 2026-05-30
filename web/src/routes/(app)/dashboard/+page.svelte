<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import { Button, Input, Dialog } from '$lib'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	let pendingSlug = $state<string | null>(null)
	let confirmOpen = $state(false)

	$effect(() => {
		if (!confirmOpen) pendingSlug = null
	})

	function avatarHue(str: string) {
		let hash = 0
		for (let i = 0; i < str.length; i++) hash = str.charCodeAt(i) + ((hash << 5) - hash)
		return Math.abs(hash) % 360
	}

	function initials(name: string) {
		return name.slice(0, 2).toUpperCase()
	}
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
			<li class="group flex flex-col gap-3 rounded-md bg-card p-4 shadow-sm">
				<div class="flex items-start justify-between gap-3">
					<div class="flex items-center gap-2.5">
						<div
							class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full text-xs font-semibold text-white"
							style="background-color: oklch(55% 0.15 {avatarHue(project.name)})"
						>
							{initials(project.name)}
						</div>
						<div>
							<a
								href="/{project.slug}"
								class="block text-sm font-semibold text-card-foreground hover:underline"
							>
								{project.name}
							</a>
							<p class="font-mono text-xs text-muted-foreground">{project.slug}</p>
						</div>
					</div>
					<span class="whitespace-nowrap text-xs text-muted-foreground">
						{new Date(project.created_at).toLocaleDateString()}
					</span>
				</div>

				<div class="flex items-center justify-between">
					<div class="flex items-center gap-3">
						<a href="/{project.slug}" class="text-xs text-muted-foreground hover:text-foreground">
							Links →
						</a>
						<a href="/{project.slug}/analytics" class="text-xs text-muted-foreground hover:text-foreground">
							Analytics →
						</a>
					</div>
					<Button
						variant="ghost-destructive"
						size="sm"
						onclick={() => { pendingSlug = project.slug; confirmOpen = true }}
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
		<Button variant="outline" onclick={() => (confirmOpen = false)}>Cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="slug" value={pendingSlug} />
			<Button type="submit" variant="destructive">Delete</Button>
		</form>
	{/snippet}
</Dialog>

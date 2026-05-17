<script lang="ts">
	import type { PageData, ActionData } from './$types'
	import { Button, Input, Dialog, Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '$lib'

	let { data, form }: { data: PageData; form: ActionData } = $props()

	let pendingCode = $state<string | null>(null)
	let confirmOpen = $derived(pendingCode !== null)
</script>

<div class="flex items-center gap-2">
	<a href="/dashboard" class="text-sm text-muted-foreground hover:text-foreground">Projects</a>
	<span class="text-muted-foreground">/</span>
	<h1 class="text-sm font-semibold text-foreground">{data.slug}</h1>
</div>

{#if form?.error}
	<p class="mt-4 rounded-md bg-destructive/10 px-4 py-2 text-sm text-destructive">{form.error}</p>
{/if}

<form method="POST" action="?/create" class="mt-4 grid grid-cols-[1fr_10rem_auto] gap-2">
	<Input name="url" type="url" placeholder="https://example.com" required />
	<Input name="title" type="text" placeholder="Title (optional)" />
	<Button type="submit">Shorten</Button>
</form>

<div class="mt-6">
	<Table>
		<TableHead>
			<TableRow>
				<TableHeader>Title</TableHeader>
				<TableHeader>Original URL</TableHeader>
				<TableHeader>Short URL</TableHeader>
				<TableHeader>Created</TableHeader>
				<TableHeader class="text-right"></TableHeader>
			</TableRow>
		</TableHead>
		<TableBody>
			{#each data.links.data as link (link.id)}
				<TableRow>
					<TableCell class="font-medium">{link.title ?? '—'}</TableCell>
					<TableCell class="max-w-xs">
						<a
							href={link.original_url}
							target="_blank"
							rel="noopener noreferrer"
							class="block truncate text-muted-foreground hover:text-foreground"
							title={link.original_url}
						>
							{link.original_url}
						</a>
					</TableCell>
					<TableCell>
						<a
							href={link.short_url}
							target="_blank"
							rel="noopener noreferrer"
							class="hover:underline"
						>
							{link.short_url}
						</a>
					</TableCell>
					<TableCell class="text-muted-foreground">
						{new Date(link.created_at).toLocaleDateString()}
					</TableCell>
					<TableCell class="text-right">
						<Button
							variant="ghost-destructive"
							size="sm"
							onclick={() => (pendingCode = link.short_code)}
						>
							Delete
						</Button>
					</TableCell>
				</TableRow>
			{:else}
				<TableRow>
					<TableCell colspan={5} class="py-10 text-center text-muted-foreground">
						No links yet. Add one above.
					</TableCell>
				</TableRow>
			{/each}
		</TableBody>
	</Table>
</div>

{#if data.links.next_cursor}
	<div class="mt-4 text-center">
		<a
			href="?cursor={data.links.next_cursor}"
			class="text-sm text-muted-foreground hover:text-foreground"
		>
			Load more →
		</a>
	</div>
{/if}

<Dialog
	bind:open={confirmOpen}
	title="Delete link"
	description="This will permanently delete the shortened link. This action cannot be undone."
>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (pendingCode = null)}>Cancel</Button>
		<form method="POST" action="?/delete">
			<input type="hidden" name="code" value={pendingCode} />
			<Button type="submit" variant="destructive">Delete</Button>
		</form>
	{/snippet}
</Dialog>

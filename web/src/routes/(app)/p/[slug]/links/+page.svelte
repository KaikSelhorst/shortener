<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import * as Table from "$lib/components/ui/table";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { Tooltip } from "$lib/components/ui/tooltip";
  import DeleteLinkButton from "$lib/components/delete-link-button.svelte";
  import type { PageProps } from "./$types";

  let { data, form }: PageProps = $props();

  let submitting = $state(false);
</script>

<div class="flex flex-col gap-4 p-6">
  <div class="flex items-center justify-end">
    <Modal.Root>
      <Modal.Trigger>
        {#snippet children(open)}
          <Button onclick={open}>New link</Button>
        {/snippet}
      </Modal.Trigger>
      <Modal.Content>
        {#snippet children(close)}
          <Modal.Header>
            <Modal.Title>New link</Modal.Title>
            <Modal.Description>Shorten a URL for this project.</Modal.Description>
          </Modal.Header>

          {#if form?.error}
            <Alert token="destructive" title="Couldn't create link" class="mt-4">{form.error}</Alert>
          {/if}

          <form
            class="mt-4 flex flex-col gap-4"
            method="POST"
            action="?/create"
            use:enhance={() => {
              submitting = true;
              return async ({ update, result }) => {
                await update();
                submitting = false;
                if (result.type === "success") close();
              };
            }}
          >
            <Field label="URL">
              <Input type="url" name="url" placeholder="https://example.com/very/long/path" required />
            </Field>
            <Field label="Title" description="Optional.">
              <Input name="title" placeholder="My link" />
            </Field>
            <Field label="Custom code" description="Optional, 3–50 characters.">
              <Input name="custom_code" placeholder="my-link" minlength={3} maxlength={50} />
            </Field>
            <Modal.Footer>
              <Button type="button" variant="outline" onclick={close}>Cancel</Button>
              <Button type="submit" disabled={submitting}>{submitting ? "Creating…" : "Create"}</Button>
            </Modal.Footer>
          </form>
        {/snippet}
      </Modal.Content>
    </Modal.Root>
  </div>

  {#if data.links.data.length === 0}
    <p class="text-sm text-muted-foreground">No links yet.</p>
  {:else}
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Short link</Table.Head>
          <Table.Head>Original URL</Table.Head>
          <Table.Head>Clicks</Table.Head>
          <Table.Head>Expires</Table.Head>
          <Table.Head>Created</Table.Head>
          <Table.Head>Actions</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each data.links.data as link (link.id)}
          <Table.Row>
            <Table.Cell class="max-w-32">
              <a
                href={link.short_url}
                target="_blank"
                rel="noopener noreferrer"
                class="block truncate font-medium text-foreground hover:underline"
              >
                {link.short_url}
              </a>
            </Table.Cell>
            <Table.Cell class="max-w-xs truncate text-muted-foreground">{link.original_url}</Table.Cell>
            <Table.Cell>{link.total_clicks}{link.max_clicks ? `/${link.max_clicks}` : ""}</Table.Cell>
            <Table.Cell class="text-muted-foreground">
              {link.expires_at ? new Date(link.expires_at).toLocaleDateString() : "—"}
            </Table.Cell>
            <Table.Cell class="text-muted-foreground">{new Date(link.created_at).toLocaleDateString()}</Table.Cell>
            <Table.Cell>
              <div class="flex items-center gap-1">
                <Tooltip label="View analytics">
                  <a
                    href="/p/{data.project.slug}/links/{link.short_code}"
                    aria-label="View analytics"
                    class="rounded-md p-1.5 text-muted-foreground hover:bg-accent hover:text-foreground"
                  >
                    <svg
                      class="size-4"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <line x1="18" y1="20" x2="18" y2="10" />
                      <line x1="12" y1="20" x2="12" y2="4" />
                      <line x1="6" y1="20" x2="6" y2="14" />
                    </svg>
                  </a>
                </Tooltip>
                <DeleteLinkButton code={link.short_code} />
              </div>
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>

    {#if data.links.prev_cursor || data.links.next_cursor}
      <div class="flex items-center justify-end gap-2">
        <Button
          variant="outline"
          size="sm"
          href={data.links.prev_cursor ? `?cursor=${encodeURIComponent(data.links.prev_cursor)}` : undefined}
          disabled={!data.links.prev_cursor}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          href={data.links.next_cursor ? `?cursor=${encodeURIComponent(data.links.next_cursor)}` : undefined}
          disabled={!data.links.next_cursor}
        >
          Next
        </Button>
      </div>
    {/if}
  {/if}
</div>

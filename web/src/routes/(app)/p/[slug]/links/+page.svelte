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
  import EditLinkButton from "$lib/components/edit-link-button.svelte";
  import { createFormSubmit } from "$lib/form-submit.svelte";
  import { formatDate } from "$lib/format";
  import { untrack } from "svelte";
  import type { PageProps } from "./$types";

  let { data, form }: PageProps = $props();

  let links = $state(untrack(() => data.links.data));
  $effect(() => {
    links = data.links.data;
  });

  $effect(() => {
    const source = new EventSource(`/stream/projects/${data.project.slug}/stream`);
    source.onmessage = (event) => {
      const evt: { short_code: string } = JSON.parse(event.data);
      const link = links.find((l) => l.short_code === evt.short_code);
      if (link) link.total_clicks += 1;
    };
    return () => source.close();
  });
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

          {@const formSubmit = createFormSubmit(close)}
          <form
            class="mt-4 flex flex-col gap-4"
            method="POST"
            action="?/create"
            use:enhance={(input) => {
              const expiresAt = input.formData.get("expires_at");
              if (typeof expiresAt === "string" && expiresAt) {
                input.formData.set("expires_at", new Date(expiresAt).toISOString());
              }
              return formSubmit.submit(input);
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
            <Field label="Max clicks" description="Optional. The link stops working after this many clicks.">
              <Input type="number" name="max_clicks" min="1" placeholder="Unlimited" />
            </Field>
            <Field label="Expires at" description="Optional, in your local time.">
              <Input type="datetime-local" name="expires_at" />
            </Field>
            <Modal.Footer>
              <Button type="button" variant="outline" onclick={close}>Cancel</Button>
              <Button type="submit" disabled={formSubmit.submitting}>
                {formSubmit.submitting ? "Creating…" : "Create"}
              </Button>
            </Modal.Footer>
          </form>
        {/snippet}
      </Modal.Content>
    </Modal.Root>
  </div>

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
      {#if links.length === 0}
        <Table.Row>
          <Table.Cell colspan={6} class="text-center text-muted-foreground">No links yet.</Table.Cell>
        </Table.Row>
      {:else}
        {#each links as link (link.id)}
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
            <Table.Cell class="text-muted-foreground">{formatDate(link.expires_at)}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{formatDate(link.created_at)}</Table.Cell>
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
                <EditLinkButton {link} />
                <DeleteLinkButton code={link.short_code} />
              </div>
            </Table.Cell>
          </Table.Row>
        {/each}
      {/if}
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
</div>

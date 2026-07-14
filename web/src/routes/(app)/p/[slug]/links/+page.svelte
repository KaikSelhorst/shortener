<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import * as Table from "$lib/components/ui/table";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import type { PageProps } from "./$types";

  let { data, form }: PageProps = $props();

  let submitting = $state(false);
</script>

<div class="mx-auto flex max-w-4xl flex-col gap-4 p-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-semibold text-foreground">Links</h1>
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
          <Table.Head>Created</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each data.links.data as link (link.id)}
          <Table.Row>
            <Table.Cell>
              <a
                href="/p/{data.project.slug}/links/{link.short_code}"
                class="font-medium text-foreground hover:underline"
              >
                {link.short_code}
              </a>
            </Table.Cell>
            <Table.Cell class="max-w-xs truncate text-muted-foreground">{link.original_url}</Table.Cell>
            <Table.Cell>{link.total_clicks}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{new Date(link.created_at).toLocaleDateString()}</Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  {/if}
</div>

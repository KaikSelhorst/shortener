<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Card from "$lib/components/ui/card";
  import * as Modal from "$lib/components/ui/modal";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { createFormSubmit } from "$lib/form-submit.svelte";
  import type { PageProps } from "./$types";

  let { data, form }: PageProps = $props();

  const renameSubmit = createFormSubmit();
</script>

<div class="flex flex-col gap-4 p-6">
  <Card.Root>
    <Card.Header>
      <Card.Title>Project details</Card.Title>
    </Card.Header>
    <Card.Content>
      {#if form?.error}
        <Alert token="destructive" title="Couldn't update project" class="mb-4">{form.error}</Alert>
      {/if}
      <form class="flex items-end gap-2" method="POST" action="?/rename" use:enhance={renameSubmit.submit}>
        <Field label="Name" class="flex-1">
          <Input name="name" value={data.project.name} required maxlength={100} />
        </Field>
        <Button type="submit" disabled={renameSubmit.submitting}>
          {renameSubmit.submitting ? "Saving…" : "Save"}
        </Button>
      </form>
    </Card.Content>
  </Card.Root>

  <Card.Root>
    <Card.Header>
      <Card.Title>Danger zone</Card.Title>
      <Card.Description>Deleting a project permanently removes its links, clicks, and webhooks.</Card.Description>
    </Card.Header>
    <Card.Content>
      <Modal.Root>
        <Modal.Trigger>
          {#snippet children(open)}
            <Button variant="destructive" onclick={open}>Delete project</Button>
          {/snippet}
        </Modal.Trigger>
        <Modal.Content>
          {#snippet children(close)}
            <Modal.Header>
              <Modal.Title>Delete "{data.project.name}"?</Modal.Title>
              <Modal.Description>
                This can't be undone. All links, clicks, and webhooks in this project will be permanently deleted.
              </Modal.Description>
            </Modal.Header>
            {@const deleteSubmit = createFormSubmit(close)}
            <form method="POST" action="?/delete" use:enhance={deleteSubmit.submit}>
              <Modal.Footer>
                <Button type="button" variant="outline" onclick={close}>Cancel</Button>
                <Button type="submit" variant="destructive" disabled={deleteSubmit.submitting}>
                  {deleteSubmit.submitting ? "Deleting…" : "Delete"}
                </Button>
              </Modal.Footer>
            </form>
          {/snippet}
        </Modal.Content>
      </Modal.Root>
    </Card.Content>
  </Card.Root>
</div>

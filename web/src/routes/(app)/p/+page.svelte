<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { createFormSubmit } from "$lib/form-submit.svelte";
  import type { PageProps } from "./$types";

  let { form }: PageProps = $props();

  const formSubmit = createFormSubmit();
</script>

<div class="flex min-h-screen flex-col items-center justify-center gap-4 p-4 text-center">
  <div class="flex flex-col gap-1">
    <h1 class="text-2xl font-semibold text-foreground">Create your first project</h1>
    <p class="text-sm text-muted-foreground">Projects group your short links, analytics, and webhooks.</p>
  </div>

  {#if form?.error}
    <Alert token="destructive" title="Couldn't create project">{form.error}</Alert>
  {/if}

  <Modal.Root>
    <Modal.Trigger>
      {#snippet children(open)}
        <Button onclick={open}>New project</Button>
      {/snippet}
    </Modal.Trigger>
    <Modal.Content>
      {#snippet children(close)}
        <Modal.Header>
          <Modal.Title>New project</Modal.Title>
          <Modal.Description>Give it a name — the URL slug is generated automatically.</Modal.Description>
        </Modal.Header>
        <form class="mt-4" method="POST" use:enhance={formSubmit.submit}>
          <Field label="Name">
            <Input name="name" placeholder="My project" required maxlength={100} />
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

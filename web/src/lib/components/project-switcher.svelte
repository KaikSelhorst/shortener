<script lang="ts">
  import { enhance } from "$app/forms";
  import { goto } from "$app/navigation";
  import * as Popover from "$lib/components/ui/popover";
  import * as Modal from "$lib/components/ui/modal";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import type { Project } from "$lib/types";

  interface Props {
    projects: Project[];
    active: Project;
  }

  let { projects, active }: Props = $props();

  let submitting = $state(false);
  let createError = $state<string | null>(null);
</script>

<Popover.Root class="w-full">
  <Popover.Trigger>
    {#snippet children(toggle)}
      <button
        type="button"
        onclick={toggle}
        class="flex h-14 w-full items-center justify-between border-b border-border px-4 text-sm font-medium text-foreground hover:bg-accent"
      >
        <span class="truncate">{active.name}</span>
        <svg
          class="size-4 shrink-0 text-muted-foreground"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M6 9l6 6 6-6" />
        </svg>
      </button>
    {/snippet}
  </Popover.Trigger>

  <Popover.Content class="w-full">
    {#snippet children(close)}
      <p class="px-2.5 py-2 text-xs font-medium text-muted-foreground">Projects</p>
      {#each projects as project (project.id)}
        <a
          href="/p/{project.slug}/links"
          onclick={close}
          class="block truncate rounded-md px-2.5 py-2 text-left text-sm hover:bg-accent {project.id === active.id
            ? 'text-foreground'
            : 'text-muted-foreground'}"
        >
          {project.name}
        </a>
      {/each}

      <div class="my-1 border-t border-border"></div>

      <Modal.Root>
        <Modal.Trigger>
          {#snippet children(open)}
            <button
              type="button"
              onclick={() => {
                close();
                open();
              }}
              class="flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm text-foreground hover:bg-accent"
            >
              + New project
            </button>
          {/snippet}
        </Modal.Trigger>
        <Modal.Content>
          {#snippet children(closeModal)}
            <Modal.Header>
              <Modal.Title>New project</Modal.Title>
              <Modal.Description>Give it a name — the URL slug is generated automatically.</Modal.Description>
            </Modal.Header>

            {#if createError}
              <Alert token="destructive" title="Couldn't create project" class="mt-4">{createError}</Alert>
            {/if}

            <form
              class="mt-4"
              method="POST"
              action="/p"
              use:enhance={() => {
                submitting = true;
                createError = null;
                return async ({ result }) => {
                  submitting = false;
                  if (result.type === "failure") {
                    createError = (result.data as { error?: string } | undefined)?.error ?? "Something went wrong.";
                  } else if (result.type === "redirect") {
                    closeModal();
                    await goto(result.location, { invalidateAll: true });
                  }
                };
              }}
            >
              <Field label="Name">
                <Input name="name" placeholder="My project" required maxlength={100} />
              </Field>
              <Modal.Footer>
                <Button type="button" variant="outline" onclick={closeModal}>Cancel</Button>
                <Button type="submit" disabled={submitting}>{submitting ? "Creating…" : "Create"}</Button>
              </Modal.Footer>
            </form>
          {/snippet}
        </Modal.Content>
      </Modal.Root>
    {/snippet}
  </Popover.Content>
</Popover.Root>

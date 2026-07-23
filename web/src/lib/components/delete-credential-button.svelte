<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import { Button } from "$lib/components/ui/button";
  import { Tooltip } from "$lib/components/ui/tooltip";
  import { createFormSubmit } from "$lib/form-submit.svelte";

  interface Props {
    id: number;
  }

  let { id }: Props = $props();
</script>

<Modal.Root>
  <Modal.Trigger>
    {#snippet children(open)}
      <Tooltip label="Delete credential">
        <button
          type="button"
          onclick={open}
          aria-label="Delete credential"
          class="rounded-md p-1.5 text-muted-foreground hover:bg-accent hover:text-destructive"
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
            <path d="M4 7h16" />
            <path d="M9 7V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v3" />
            <path d="M6 7l1 12a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-12" />
          </svg>
        </button>
      </Tooltip>
    {/snippet}
  </Modal.Trigger>
  <Modal.Content>
    {#snippet children(close)}
      <Modal.Header>
        <Modal.Title>Delete credential?</Modal.Title>
        <Modal.Description>This can't be undone. Any requests using this credential will start failing immediately.</Modal.Description>
      </Modal.Header>
      {@const formSubmit = createFormSubmit(close)}
      <form method="POST" action="?/delete" use:enhance={formSubmit.submit}>
        <input type="hidden" name="id" value={id} />
        <Modal.Footer>
          <Button type="button" variant="outline" onclick={close}>Cancel</Button>
          <Button type="submit" variant="destructive" disabled={formSubmit.submitting}>
            {formSubmit.submitting ? "Deleting…" : "Delete"}
          </Button>
        </Modal.Footer>
      </form>
    {/snippet}
  </Modal.Content>
</Modal.Root>

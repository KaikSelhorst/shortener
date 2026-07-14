<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import { Button } from "$lib/components/ui/button";
  import { Tooltip } from "$lib/components/ui/tooltip";

  interface Props {
    code: string;
  }

  let { code }: Props = $props();
  let submitting = $state(false);
</script>

<Modal.Root>
  <Modal.Trigger>
    {#snippet children(open)}
      <Tooltip label="Delete link">
        <button
          type="button"
          onclick={open}
          aria-label="Delete link"
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
        <Modal.Title>Delete link?</Modal.Title>
        <Modal.Description>This can't be undone. The short link will stop working immediately.</Modal.Description>
      </Modal.Header>
      <form
        method="POST"
        action="?/delete"
        use:enhance={() => {
          submitting = true;
          return async ({ update, result }) => {
            await update();
            submitting = false;
            if (result.type === "success") close();
          };
        }}
      >
        <input type="hidden" name="code" value={code} />
        <Modal.Footer>
          <Button type="button" variant="outline" onclick={close}>Cancel</Button>
          <Button type="submit" variant="destructive" disabled={submitting}>
            {submitting ? "Deleting…" : "Delete"}
          </Button>
        </Modal.Footer>
      </form>
    {/snippet}
  </Modal.Content>
</Modal.Root>

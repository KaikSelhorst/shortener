<script lang="ts">
  import * as Modal from "$lib/components/ui/modal";
  import { Button } from "$lib/components/ui/button";
</script>

<Modal.Root>
  <Modal.Trigger>
    {#snippet children(open)}
      <button
        type="button"
        onclick={open}
        class="w-full rounded-md px-3 py-2 text-left text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground"
      >
        Logout
      </button>
    {/snippet}
  </Modal.Trigger>
  <Modal.Content>
    {#snippet children(close)}
      <Modal.Header>
        <Modal.Title>Log out?</Modal.Title>
        <Modal.Description>You'll need to sign in again to access your projects.</Modal.Description>
      </Modal.Header>
      <Modal.Footer>
        <Button type="button" variant="outline" onclick={close}>Cancel</Button>
        <form method="POST" action="/logout">
          <Button type="submit">Log out</Button>
        </form>
      </Modal.Footer>
    {/snippet}
  </Modal.Content>
</Modal.Root>

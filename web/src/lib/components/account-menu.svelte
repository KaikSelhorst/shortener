<script lang="ts">
  import * as Popover from "$lib/components/ui/popover";
  import * as Modal from "$lib/components/ui/modal";
  import { Button } from "$lib/components/ui/button";

  let logoutModalOpen = $state(false);
</script>

<Popover.Root class="w-full">
  <Popover.Trigger>
    {#snippet children(toggle)}
      <button
        type="button"
        onclick={toggle}
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground"
      >
        <svg
          class="size-4 shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
          <circle cx="12" cy="7" r="4" />
        </svg>
        Account
      </button>
    {/snippet}
  </Popover.Trigger>

  <Popover.Content class="bottom-full mb-1 w-full">
    {#snippet children(close)}
      <a
        href="/settings/credentials"
        onclick={close}
        class="flex items-center gap-2 rounded-md px-2.5 py-2 text-sm text-foreground hover:bg-accent"
      >
        <svg
          class="size-4 shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="3" />
          <path
            d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"
          />
        </svg>
        Settings
      </a>
      <button
        type="button"
        onclick={() => {
          close();
          logoutModalOpen = true;
        }}
        class="flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm text-foreground hover:bg-accent"
      >
        <svg
          class="size-4 shrink-0"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
          <polyline points="16 17 21 12 16 7" />
          <line x1="21" y1="12" x2="9" y2="12" />
        </svg>
        Logout
      </button>
    {/snippet}
  </Popover.Content>
</Popover.Root>

<Modal.Root bind:open={logoutModalOpen}>
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

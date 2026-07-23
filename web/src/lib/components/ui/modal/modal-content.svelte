<script lang="ts">
  import type { Snippet } from "svelte";
  import { getModalContext } from "./context";

  interface Props {
    children: Snippet<[close: () => void]>;
  }

  let { children }: Props = $props();
  const ctx = getModalContext();
</script>

{#if ctx.open}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <button
      type="button"
      class="fixed inset-0 bg-background/80 backdrop-blur-sm"
      onclick={() => ctx.setOpen(false)}
      aria-label="Close"
    ></button>
    <div class="relative w-full max-w-md rounded-lg border border-border bg-popover p-6 shadow-lg">
      {@render children(() => ctx.setOpen(false))}
    </div>
  </div>
{/if}

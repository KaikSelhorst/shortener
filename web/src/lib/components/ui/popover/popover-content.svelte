<script lang="ts">
  import type { Snippet } from "svelte";
  import { getPopoverContext } from "./context";

  interface Props {
    children: Snippet<[close: () => void]>;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
  const ctx = getPopoverContext();
</script>

{#if ctx.open}
  <button
    type="button"
    class="fixed inset-0 z-10 cursor-default"
    onclick={() => ctx.setOpen(false)}
    aria-label="Close popover"
  ></button>
  <div class="absolute left-0 z-20 mt-2 rounded-lg border border-border bg-popover p-1 shadow-lg {className}">
    {@render children(() => ctx.setOpen(false))}
  </div>
{/if}

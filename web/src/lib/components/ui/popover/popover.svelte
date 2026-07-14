<script lang="ts">
  import type { Snippet } from "svelte";
  import { setPopoverContext } from "./context";

  interface Props {
    open?: boolean;
    class?: string;
    children: Snippet;
  }

  let { open = $bindable(false), class: className = "", children }: Props = $props();

  setPopoverContext({
    get open() {
      return open;
    },
    setOpen(value: boolean) {
      open = value;
    },
  });
</script>

<div class="relative {className}">
  {@render children()}
</div>

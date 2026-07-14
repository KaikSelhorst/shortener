<script lang="ts">
  import type { Snippet } from "svelte";
  import { getPopoverContext } from "./context";

  interface Props {
    children: Snippet<[toggle: () => void]>;
  }

  let { children }: Props = $props();
  const ctx = getPopoverContext();
</script>

{@render children(() => ctx.setOpen(!ctx.open))}

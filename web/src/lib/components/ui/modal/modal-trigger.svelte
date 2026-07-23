<script lang="ts">
  import type { Snippet } from "svelte";
  import { getModalContext } from "./context";

  interface Props {
    children: Snippet<[open: () => void]>;
  }

  let { children }: Props = $props();
  const ctx = getModalContext();
</script>

{@render children(() => ctx.setOpen(true))}

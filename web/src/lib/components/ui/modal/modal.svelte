<script lang="ts">
  import type { Snippet } from "svelte";
  import { setModalContext } from "./context";

  interface Props {
    open?: boolean;
    children: Snippet;
  }

  let { open = $bindable(false), children }: Props = $props();

  setModalContext({
    get open() {
      return open;
    },
    setOpen(value: boolean) {
      open = value;
    },
  });
</script>

{@render children()}

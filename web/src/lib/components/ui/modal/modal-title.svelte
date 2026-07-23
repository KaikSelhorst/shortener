<script lang="ts">
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

<h3 class="text-base font-medium text-popover-foreground">{@render children()}</h3>

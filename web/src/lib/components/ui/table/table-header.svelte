<script lang="ts">
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

<!-- box-shadow, not border-b: a real border on a sticky thead/tr stops painting once it's pinned (Chromium clips it in the sticky layer) -->
<thead
  class="sticky top-0 z-10 rounded-t-lg bg-card shadow-[0_1px_0_0_var(--color-border)] [&>tr]:border-b-0"
>
  {@render children()}
</thead>

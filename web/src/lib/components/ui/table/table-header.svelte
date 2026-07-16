<script lang="ts">
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

<thead class="sticky top-0 z-10 rounded-t-lg bg-card">
  {@render children()}
</thead>

<script lang="ts">
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

<tbody class="[&>tr:last-child]:border-b-0 [&_td:hover]:bg-accent/30">
  {@render children()}
</tbody>

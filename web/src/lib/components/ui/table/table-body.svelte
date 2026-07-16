<script lang="ts">
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

<tbody class="[&>tr:last-child]:border-b-0 [&>tr:hover]:bg-accent">
  {@render children()}
</tbody>

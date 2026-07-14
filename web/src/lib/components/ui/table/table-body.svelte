<script lang="ts">
  import type { Snippet } from "svelte";

  let { children }: { children: Snippet } = $props();
</script>

<tbody class="[&>tr:last-child]:border-b-0">
  {@render children()}
</tbody>

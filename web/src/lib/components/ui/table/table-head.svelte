<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<th class="px-2 py-3 text-left font-medium text-muted-foreground {className}">
  {@render children()}
</th>

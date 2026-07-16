<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<th class="border-r border-border px-2 py-2 text-left font-medium text-muted-foreground {className}">
  {@render children()}
</th>

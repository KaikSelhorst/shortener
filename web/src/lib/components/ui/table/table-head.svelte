<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<!-- box-shadow, not border-r: matches table-cell.svelte's divider mechanism (see its comment) -->
<th
  class="px-2 py-2 text-left font-medium text-muted-foreground shadow-[inset_-1px_0_0_0_var(--color-border)] {className}"
>
  {@render children()}
</th>

<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLTdAttributes } from "svelte/elements";

  interface Props extends HTMLTdAttributes {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "", ...rest }: Props = $props();
</script>

<!-- box-shadow, not border-r: table-collapse aligns real borders to the shared column edge, but
     table-head.svelte's divider is a box-shadow (which ignores collapse) — matching the mechanism
     here keeps the two dividers pixel-aligned instead of offset from each other -->
<td {...rest} class="px-2 py-2 shadow-[inset_-1px_0_0_0_var(--color-border)] {className}">
  {@render children()}
</td>

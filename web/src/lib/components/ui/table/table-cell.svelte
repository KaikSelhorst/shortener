<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLTdAttributes } from "svelte/elements";

  interface Props extends HTMLTdAttributes {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "", ...rest }: Props = $props();
</script>

<td {...rest} class="border-r border-border px-2 py-2 {className}">
  {@render children()}
</td>

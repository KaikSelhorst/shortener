<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<td class="border-r border-border px-2 py-2 {className}">
  {@render children()}
</td>

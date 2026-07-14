<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<tr class="border-b border-border hover:bg-accent {className}">
  {@render children()}
</tr>

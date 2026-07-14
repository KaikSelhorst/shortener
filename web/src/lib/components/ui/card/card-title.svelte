<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<h3 class="text-sm font-medium text-foreground {className}">
  {@render children()}
</h3>

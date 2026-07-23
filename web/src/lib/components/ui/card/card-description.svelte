<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<p class="text-sm text-muted-foreground {className}">
  {@render children()}
</p>

<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<div class="rounded-lg border border-border bg-card p-4 {className}">
  {@render children()}
</div>

<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    children: Snippet;
    class?: string;
  }

  let { children, class: className = "" }: Props = $props();
</script>

<div class="overflow-hidden rounded-lg border border-border bg-card {className}">
  <table class="w-full text-sm">
    {@render children()}
  </table>
</div>

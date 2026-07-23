<script lang="ts">
  import type { Snippet } from "svelte";

  interface Props {
    label?: string;
    error?: string;
    description?: string;
    class?: string;
    children: Snippet;
  }

  let { label, error, description, class: className = "", children }: Props = $props();
</script>

<label class="flex flex-col gap-1.5 {className}">
  {#if label}
    <span class="text-xs font-medium text-muted-foreground">{label}</span>
  {/if}
  {@render children()}
  {#if error}
    <span class="text-xs text-destructive">{error}</span>
  {:else if description}
    <span class="text-xs text-muted-foreground">{description}</span>
  {/if}
</label>

<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends HTMLInputAttributes {
    error?: boolean;
  }

  let { error = false, class: className = "", ...rest }: Props = $props();
</script>

<input
  {...rest}
  class="w-full rounded-md border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 disabled:cursor-not-allowed disabled:opacity-50 {error
    ? 'border-destructive focus-visible:ring-destructive'
    : 'border-input focus-visible:ring-ring'} {className}"
/>

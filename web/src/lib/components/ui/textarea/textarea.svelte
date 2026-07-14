<script lang="ts">
  import type { HTMLTextareaAttributes } from "svelte/elements";

  interface Props extends HTMLTextareaAttributes {
    error?: boolean;
  }

  let { error = false, class: className = "", rows = 3, ...rest }: Props = $props();
</script>

<textarea
  {rows}
  {...rest}
  class="resize-none rounded-md border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 disabled:cursor-not-allowed disabled:opacity-50 {error
    ? 'border-destructive focus-visible:ring-destructive'
    : 'border-input focus-visible:ring-ring'} {className}"
></textarea>

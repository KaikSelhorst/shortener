<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLSelectAttributes } from "svelte/elements";

  interface Props extends HTMLSelectAttributes {
    children: Snippet;
  }

  let { class: className = "", children, ...rest }: Props = $props();
</script>

<select
  {...rest}
  class="rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 {className}"
>
  {@render children()}
</select>

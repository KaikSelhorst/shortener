<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAnchorAttributes } from "svelte/elements";

  interface Props extends HTMLAnchorAttributes {
    active?: boolean;
    children: Snippet;
  }

  let { active = false, class: className = "", children, ...rest }: Props = $props();
</script>

<a
  {...rest}
  class="rounded-md px-3 py-2 text-sm font-medium transition-colors {active
    ? 'bg-secondary text-foreground'
    : 'text-muted-foreground hover:bg-accent hover:text-foreground'} {className}"
>
  {@render children()}
</a>

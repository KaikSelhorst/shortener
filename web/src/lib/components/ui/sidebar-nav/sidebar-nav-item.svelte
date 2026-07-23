<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAnchorAttributes } from "svelte/elements";

  interface Props extends HTMLAnchorAttributes {
    active?: boolean;
    icon?: Snippet;
    children: Snippet;
  }

  let { active = false, class: className = "", icon, children, ...rest }: Props = $props();
</script>

<a
  {...rest}
  class="flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors {active
    ? 'bg-secondary text-foreground'
    : 'text-muted-foreground hover:bg-accent hover:text-foreground'} {className}"
>
  {#if icon}
    {@render icon()}
  {/if}
  {@render children()}
</a>

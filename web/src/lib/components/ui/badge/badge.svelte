<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLAttributes } from "svelte/elements";
  import type { BadgeToken } from "./badge-variants";

  interface Props extends HTMLAttributes<HTMLSpanElement> {
    token?: BadgeToken;
    children: Snippet;
  }

  let { token = "neutral", class: className = "", children, ...rest }: Props = $props();
</script>

<span
  {...rest}
  class="inline-flex items-center rounded-md px-2.5 py-1 text-xs font-medium {className}"
  style="
    background: color-mix(in oklab, var(--{token}) 18%, var(--card));
    color: var(--{token});
  "
>
  {@render children()}
</span>

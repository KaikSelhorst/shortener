<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";
  import { buttonVariants, buttonSizes, type ButtonVariant, type ButtonSize } from "./button-variants";

  interface Props extends HTMLButtonAttributes {
    variant?: ButtonVariant;
    size?: ButtonSize;
    children: Snippet;
  }

  let { variant = "primary", size = "md", type = "button", class: className = "", children, ...rest }: Props = $props();
</script>

<button
  {type}
  {...rest}
  class="rounded-md font-medium transition-opacity disabled:cursor-not-allowed disabled:opacity-40 {buttonVariants[
    variant
  ]} {buttonSizes[size]} {className}"
>
  {@render children()}
</button>

<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";
  import { buttonVariants, buttonSizes, type ButtonVariant, type ButtonSize } from "./button-variants";

  interface Props extends Omit<HTMLButtonAttributes, "onclick"> {
    variant?: ButtonVariant;
    size?: ButtonSize;
    href?: string;
    onclick?: (event: MouseEvent) => void;
    children: Snippet;
  }

  let {
    variant = "primary",
    size = "md",
    type = "button",
    href,
    disabled,
    class: className = "",
    children,
    onclick,
    ...rest
  }: Props = $props();

  let classes = $derived(
    `rounded-md font-medium transition-opacity disabled:cursor-not-allowed disabled:opacity-40 ${buttonVariants[variant]} ${buttonSizes[size]} ${className}`,
  );
</script>

{#if href}
  <a
    href={disabled ? undefined : href}
    aria-disabled={disabled}
    {onclick}
    class="{classes} {disabled ? 'pointer-events-none' : ''}"
  >
    {@render children()}
  </a>
{:else}
  <button {type} {disabled} {onclick} {...rest} class={classes}>
    {@render children()}
  </button>
{/if}

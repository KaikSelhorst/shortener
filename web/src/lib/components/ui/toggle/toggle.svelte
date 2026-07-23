<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends Omit<HTMLInputAttributes, "type"> {}

  let { class: className = "", ...rest }: Props = $props();
</script>

<span class="relative inline-flex h-5 w-9 shrink-0 items-center {className}">
  <input type="checkbox" {...rest} class="peer sr-only" />
  <span
    class="absolute inset-0 rounded-md bg-secondary transition-colors peer-checked:bg-lime peer-focus-visible:ring-2 peer-focus-visible:ring-ring peer-focus-visible:ring-offset-2 peer-focus-visible:ring-offset-card"
  ></span>
  <span
    class="absolute left-0.5 size-4 rounded-md bg-foreground transition-transform peer-checked:translate-x-4 peer-checked:bg-background"
  ></span>
</span>

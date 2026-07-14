<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends Omit<HTMLInputAttributes, "type"> {}

  let { class: className = "", ...rest }: Props = $props();
</script>

<span class="relative inline-flex size-4 shrink-0 {className}">
  <input
    type="radio"
    {...rest}
    class="peer size-4 shrink-0 appearance-none rounded-full border border-input bg-background checked:border-lime focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-card disabled:cursor-not-allowed"
  />
  <span class="pointer-events-none absolute inset-0 hidden items-center justify-center peer-checked:flex">
    <span class="size-1.5 rounded-full bg-lime"></span>
  </span>
</span>

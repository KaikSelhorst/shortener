<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends Omit<HTMLInputAttributes, "type" | "value"> {
    value?: number | string;
  }

  let { value = $bindable(), class: className = "", ...rest }: Props = $props();
</script>

<input type="range" bind:value {...rest} class={className} />

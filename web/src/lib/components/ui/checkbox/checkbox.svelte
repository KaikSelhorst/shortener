<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  interface Props extends Omit<HTMLInputAttributes, "type"> {
    indeterminate?: boolean;
  }

  let { indeterminate = false, class: className = "", ...rest }: Props = $props();
</script>

{#if indeterminate}
  <span class="flex size-4 shrink-0 items-center justify-center rounded border border-secondary bg-secondary {className}">
    <span class="h-0.5 w-2 rounded-full bg-lime"></span>
  </span>
{:else}
  <span class="relative inline-flex size-4 shrink-0 {className}">
    <input
      type="checkbox"
      {...rest}
      class="peer size-4 shrink-0 appearance-none rounded border border-input bg-background checked:border-lime checked:bg-lime focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed"
    />
    <svg
      class="pointer-events-none absolute inset-0 hidden size-4 p-0.5 text-background peer-checked:block"
      viewBox="0 0 16 16"
      fill="none"
      stroke="currentColor"
      stroke-width="2.5"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d="M3 8l3.5 3.5L13 5" />
    </svg>
  </span>
{/if}

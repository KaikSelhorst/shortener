<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";
  import { Input } from "$lib/components/ui/input";

  interface Props extends Omit<HTMLInputAttributes, "type"> {
    error?: boolean;
  }

  let { error = false, class: className = "", ...rest }: Props = $props();

  let visible = $state(false);
</script>

<div class="relative">
  <Input
    type={visible ? "text" : "password"}
    {error}
    class={className}
    style="padding-right: 2.25rem;"
    {...rest}
  />
  <button
    type="button"
    onclick={() => (visible = !visible)}
    class="flex items-center px-2.5 text-muted-foreground hover:text-foreground"
    style="position: absolute; top: 0; bottom: 0; right: 0;"
    aria-label={visible ? "Hide password" : "Show password"}
  >
    <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      {#if visible}
        <path d="M9.9 4.24A10.94 10.94 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
        <path d="M6.61 6.61A18.5 18.5 0 0 0 1 12s4 8 11 8a10.94 10.94 0 0 0 5.39-1.61" />
        <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
        <path d="M1 1l22 22" />
      {:else}
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
        <circle cx="12" cy="12" r="3" />
      {/if}
    </svg>
  </button>
</div>

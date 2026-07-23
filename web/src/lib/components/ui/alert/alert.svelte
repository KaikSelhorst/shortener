<script lang="ts">
  import type { Snippet } from "svelte";

  export type AlertToken = "info" | "success" | "warning" | "destructive";

  interface Props {
    token?: AlertToken;
    title: string;
    children: Snippet;
    onDismiss?: () => void;
    class?: string;
  }

  let { token = "info", title, children, onDismiss, class: className = "" }: Props = $props();
</script>

<div
  class="flex gap-3 rounded-lg border p-4 {className}"
  style="
    background: color-mix(in oklab, var(--{token}) 10%, var(--card));
    border-color: color-mix(in oklab, var(--{token}) 35%, transparent);
  "
>
  <svg
    class="size-5 shrink-0"
    style="color: var(--{token})"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  >
    {#if token === "info"}
      <circle cx="12" cy="12" r="9" />
      <path d="M12 8h.01" />
      <path d="M11 11.5h1v4.5" />
    {:else if token === "success"}
      <circle cx="12" cy="12" r="9" />
      <path d="M8 12.5l2.5 2.5L16 9.5" />
    {:else if token === "warning"}
      <path d="M10.29 3.86L1.82 18a1 1 0 0 0 .86 1.5h18.64a1 1 0 0 0 .86-1.5L13.71 3.86a1 1 0 0 0-1.72 0z" />
      <path d="M12 9v4" />
      <path d="M12 17h.01" />
    {:else}
      <circle cx="12" cy="12" r="9" />
      <path d="M9 9l6 6M15 9l-6 6" />
    {/if}
  </svg>
  <div class="flex-1">
    <p class="text-sm font-medium" style="color: var(--{token})">{title}</p>
    <p class="mt-0.5 text-sm" style="color: color-mix(in oklab, var(--{token}) 55%, var(--muted-foreground))">
      {@render children()}
    </p>
  </div>
  {#if onDismiss}
    <button
      type="button"
      onclick={onDismiss}
      aria-label="Dismiss"
      class="shrink-0 rounded-md p-1 text-muted-foreground hover:text-foreground"
    >
      <svg class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M6 6l12 12M18 6l-12 12" />
      </svg>
    </button>
  {/if}
</div>

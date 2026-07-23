<script lang="ts">
  import type { Snippet } from "svelte";
  import { tooltipPlacements, type TooltipPlacement } from "./tooltip-variants";

  interface Props {
    label: string;
    placement?: TooltipPlacement;
    children: Snippet;
    class?: string;
  }

  let { label, placement = "top", children, class: className = "" }: Props = $props();
</script>

<span class="group relative inline-flex {className}">
  {@render children()}
  <span
    role="tooltip"
    class="pointer-events-none absolute z-30 whitespace-nowrap rounded-md border border-border bg-popover px-2 py-1 text-xs text-popover-foreground opacity-0 shadow-lg transition-opacity group-hover:opacity-100 group-focus-within:opacity-100 {tooltipPlacements[
      placement
    ]}"
  >
    {label}
  </span>
</span>

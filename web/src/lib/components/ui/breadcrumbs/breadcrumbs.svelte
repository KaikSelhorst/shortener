<script lang="ts">
  interface BreadcrumbItem {
    label: string;
    href?: string;
  }

  interface Props {
    items: BreadcrumbItem[];
    class?: string;
  }

  let { items, class: className = "" }: Props = $props();
</script>

<nav class="flex items-center gap-1.5 text-xs font-medium {className}">
  {#each items as item, i}
    {#if i > 0}
      <svg
        class="size-3 shrink-0 text-muted-foreground"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M9 6l6 6-6 6" />
      </svg>
    {/if}
    {#if i === items.length - 1 || !item.href}
      <span class="text-foreground">{item.label}</span>
    {:else}
      <a href={item.href} class="text-muted-foreground hover:text-foreground">{item.label}</a>
    {/if}
  {/each}
</nav>

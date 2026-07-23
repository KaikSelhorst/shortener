<script lang="ts">
  interface Item {
    label: string;
    value: number;
  }

  interface Props {
    items: Item[];
  }

  let { items }: Props = $props();

  const sorted = $derived([...items].filter((item) => item.value > 0).sort((a, b) => b.value - a.value));
  const max = $derived(Math.max(...sorted.map((item) => item.value), 1));
</script>

{#if sorted.length === 0}
  <p class="text-sm text-muted-foreground">No data yet.</p>
{:else}
  <div class="flex flex-col gap-2.5">
    {#each sorted as item (item.label)}
      <div class="flex items-center gap-3">
        <span class="w-20 shrink-0 truncate text-sm capitalize text-foreground">{item.label}</span>
        <div class="h-2 flex-1 overflow-hidden rounded-full bg-secondary">
          <div class="h-full rounded-full bg-lime" style="width: {(item.value / max) * 100}%"></div>
        </div>
        <span class="w-10 shrink-0 text-right text-sm tabular-nums text-muted-foreground">
          {item.value.toLocaleString()}
        </span>
      </div>
    {/each}
  </div>
{/if}

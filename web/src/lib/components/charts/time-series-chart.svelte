<script lang="ts">
  interface Point {
    date: string;
    count: number;
  }

  interface Props {
    data: Point[];
    label?: string;
  }

  let { data, label = "Clicks" }: Props = $props();

  const width = 720;
  const height = 220;
  const padding = { top: 16, right: 12, bottom: 28, left: 36 };
  const plotWidth = width - padding.left - padding.right;
  const plotHeight = height - padding.top - padding.bottom;

  const square = 4;
  const gapX = 1.5;
  const gapY = 2;
  const rowModule = square + gapY;

  function niceTicks(max: number, count: number): number[] {
    if (max <= 0) return [0];
    const step = Math.ceil(max / count / 5) * 5 || 1;
    const ticks: number[] = [];
    for (let v = 0; v <= max + step; v += step) ticks.push(v);
    return ticks;
  }

  const yTicks = $derived(niceTicks(Math.max(...data.map((d) => d.count), 0), 4));
  const scaleMax = $derived(yTicks[yTicks.length - 1] || 1);
  const rows = $derived(Math.max(1, Math.floor(plotHeight / rowModule)));

  const groupWidth = $derived(data.length ? plotWidth / data.length : plotWidth);
  const squaresPerColumn = $derived(Math.max(1, Math.floor((groupWidth - gapX) / (square + gapX))));
  const columnContentWidth = $derived(squaresPerColumn * (square + gapX) - gapX);

  function groupStartX(i: number): number {
    return padding.left + i * groupWidth + (groupWidth - columnContentWidth) / 2;
  }
  function filledRows(count: number): number {
    return Math.round((count / scaleMax) * rows);
  }
  function rowY(r: number): number {
    return padding.top + plotHeight - (r + 1) * rowModule;
  }
  function tickY(value: number): number {
    const r = Math.round((value / scaleMax) * rows);
    return rowY(r) + square / 2;
  }

  function range(n: number): number[] {
    return Array.from({ length: n }, (_, i) => i);
  }

  function formatCompact(value: number): string {
    return value >= 1000 ? `${(value / 1000).toFixed(value % 1000 === 0 ? 0 : 1)}k` : value.toLocaleString();
  }

  const labelStep = $derived(Math.max(1, Math.ceil(data.length / 8)));
  const xLabelIndices = $derived(data.map((_, i) => i).filter((i) => i % labelStep === 0));

  function formatAxisDate(value: string): string {
    return new Date(value).toLocaleDateString(undefined, { month: "short", day: "numeric" });
  }
  function formatTooltipDate(value: string): string {
    return new Date(value).toLocaleDateString(undefined, { weekday: "short", month: "short", day: "numeric" });
  }

  let hovered = $state<number | null>(null);

  function handlePointerMove(event: PointerEvent) {
    if (data.length === 0) return;
    const svg = event.currentTarget as SVGSVGElement;
    const rect = svg.getBoundingClientRect();
    const relX = ((event.clientX - rect.left) / rect.width) * width;
    const index = Math.floor((relX - padding.left) / groupWidth);
    hovered = Math.min(data.length - 1, Math.max(0, index));
  }

  const tooltipAlign = $derived(hovered === 0 ? "start" : hovered === data.length - 1 ? "end" : "center");
  const crosshairX = $derived(hovered === null ? 0 : padding.left + (hovered + 0.5) * groupWidth);
</script>

{#if data.length === 0}
  <p class="text-sm text-muted-foreground">No data yet.</p>
{:else}
  <div class="relative">
    <svg
      viewBox="0 0 {width} {height}"
      class="w-full"
      role="img"
      aria-label="{label} over time"
      onpointermove={handlePointerMove}
      onpointerleave={() => (hovered = null)}
    >
      {#each yTicks as tick (tick)}
        <text x={padding.left - 8} y={tickY(tick)} dy="0.32em" text-anchor="end" class="fill-muted-foreground text-[7px]">
          {formatCompact(tick)}
        </text>
      {/each}

      {#each xLabelIndices as i (i)}
        <text x={groupStartX(i) + columnContentWidth / 2} y={height - 8} text-anchor="middle" class="fill-muted-foreground text-[7px]">
          {formatAxisDate(data[i].date)}
        </text>
      {/each}

      {#if hovered !== null}
        <rect
          x={padding.left + hovered * groupWidth}
          y={padding.top}
          width={groupWidth}
          height={plotHeight}
          class="fill-foreground/5"
        />
        <line x1={crosshairX} x2={crosshairX} y1={padding.top} y2={padding.top + plotHeight} class="stroke-border" stroke-width="1" stroke-dasharray="3,3" />
      {/if}

      {#each data as point, i (point.date)}
        {@const filled = filledRows(point.count)}
        {#each range(rows) as r (r)}
          {#each range(squaresPerColumn) as c (c)}
            <rect
              x={groupStartX(i) + c * (square + gapX)}
              y={rowY(r)}
              width={square}
              height={square}
              rx="1"
              class={r < filled ? "fill-lime" : "fill-border"}
            />
          {/each}
        {/each}
      {/each}
    </svg>

    {#if hovered !== null}
      <div
        class="pointer-events-none absolute top-1 flex items-center gap-2 rounded-md border border-border bg-popover px-2.5 py-1.5 text-xs shadow-md {tooltipAlign ===
        'start'
          ? 'translate-x-0'
          : tooltipAlign === 'end'
            ? '-translate-x-full'
            : '-translate-x-1/2'}"
        style="left: {(crosshairX / width) * 100}%"
      >
        <span class="size-1.5 shrink-0 rounded-full bg-lime"></span>
        <p class="text-muted-foreground">{formatTooltipDate(data[hovered].date)}</p>
        <p class="font-medium text-foreground">{data[hovered].count.toLocaleString()}</p>
      </div>
    {/if}
  </div>
{/if}

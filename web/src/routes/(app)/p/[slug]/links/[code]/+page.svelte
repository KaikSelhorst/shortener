<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import { Button } from "$lib/components/ui/button";
  import StatTile from "$lib/components/stat-tile.svelte";
  import TimeSeriesChart from "$lib/components/charts/time-series-chart.svelte";
  import RankedBarList from "$lib/components/charts/ranked-bar-list.svelte";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  const PERIOD_LABELS: Record<string, string> = { "7d": "7 days", "30d": "30 days", "90d": "90 days" };

  function toItems(record: Record<string, number>) {
    return Object.entries(record).map(([label, value]) => ({ label, value }));
  }
</script>

<div class="flex flex-col gap-4 p-6">
  <div class="flex items-start justify-between gap-4">
    <div class="flex flex-col gap-1">
      <a
        href={data.link.short_url}
        target="_blank"
        rel="noopener noreferrer"
        class="text-lg font-medium text-foreground hover:underline"
      >
        {data.link.short_url}
      </a>
      <p class="max-w-xl truncate text-sm text-muted-foreground">{data.link.original_url}</p>
    </div>

    <div class="flex items-center gap-1">
      {#each data.periods as period (period)}
        <Button href="?period={period}" variant={data.period === period ? "primary" : "ghost"} size="sm">
          {PERIOD_LABELS[period]}
        </Button>
      {/each}
    </div>
  </div>

  <div class="grid grid-cols-2 gap-4">
    <StatTile label="Total clicks" value={data.analytics.total_clicks} />
    <StatTile label="Unique clicks" value={data.analytics.unique_clicks} />
  </div>

  <Card.Root>
    <Card.Title>Clicks over time</Card.Title>
    <Card.Content>
      <TimeSeriesChart data={data.analytics.over_time} />
    </Card.Content>
  </Card.Root>

  <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
    <Card.Root>
      <Card.Title>Devices</Card.Title>
      <Card.Content>
        <RankedBarList items={toItems(data.analytics.devices)} />
      </Card.Content>
    </Card.Root>
    <Card.Root>
      <Card.Title>Browsers</Card.Title>
      <Card.Content>
        <RankedBarList items={toItems(data.analytics.browsers)} />
      </Card.Content>
    </Card.Root>
    <Card.Root>
      <Card.Title>Referrers</Card.Title>
      <Card.Content>
        <RankedBarList items={toItems(data.analytics.referrers)} />
      </Card.Content>
    </Card.Root>
  </div>
</div>

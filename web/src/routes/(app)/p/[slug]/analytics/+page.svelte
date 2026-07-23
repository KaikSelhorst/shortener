<script lang="ts">
  import * as Card from "$lib/components/ui/card";
  import * as Table from "$lib/components/ui/table";
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
  <div class="flex items-center justify-end gap-1">
    {#each data.periods as period (period)}
      <Button href="?period={period}" variant={data.period === period ? "primary" : "ghost"} size="sm">
        {PERIOD_LABELS[period]}
      </Button>
    {/each}
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

  <Card.Root>
    <Card.Title>Top links</Card.Title>
    <Card.Content>
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.Head>Link</Table.Head>
            <Table.Head>Original URL</Table.Head>
            <Table.Head>Clicks</Table.Head>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {#if data.analytics.top_links.length === 0}
            <Table.Row>
              <Table.Cell colspan={3} class="text-center text-muted-foreground">No clicks in this period.</Table.Cell>
            </Table.Row>
          {:else}
            {#each data.analytics.top_links as link (link.short_code)}
              <Table.Row>
                <Table.Cell class="font-medium text-foreground">{link.short_code}</Table.Cell>
                <Table.Cell class="max-w-xs truncate text-muted-foreground">{link.original_url}</Table.Cell>
                <Table.Cell>{link.total_clicks}</Table.Cell>
              </Table.Row>
            {/each}
          {/if}
        </Table.Body>
      </Table.Root>
    </Card.Content>
  </Card.Root>
</div>

<script lang="ts">
  import * as Modal from "$lib/components/ui/modal";
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import { Badge } from "$lib/components/ui/badge";
  import type { BadgeToken } from "$lib/components/ui/badge";
  import { formatDateTime } from "$lib/format";
  import type { PageProps } from "./$types";

  let { data }: PageProps = $props();

  const STATUS_TOKEN: Record<string, BadgeToken> = {
    delivered: "success",
    failed: "destructive",
    pending: "warning",
  };

  const currentPage = $derived(data.deliveries.page);
</script>

<div class="flex flex-col gap-4 p-6">
  <Table.Root>
    <Table.Header>
      <Table.Row>
        <Table.Head>Event</Table.Head>
        <Table.Head>Status</Table.Head>
        <Table.Head>Attempts</Table.Head>
        <Table.Head>Response</Table.Head>
        <Table.Head>Created</Table.Head>
        <Table.Head>Payload</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#if data.deliveries.data.length === 0}
        <Table.Row>
          <Table.Cell colspan={6} class="text-center text-muted-foreground">No deliveries yet.</Table.Cell>
        </Table.Row>
      {:else}
        {#each data.deliveries.data as delivery (delivery.id)}
          <Table.Row>
            <Table.Cell class="font-medium text-foreground">{delivery.event}</Table.Cell>
            <Table.Cell>
              <Badge token={STATUS_TOKEN[delivery.status] ?? "neutral"}>{delivery.status}</Badge>
            </Table.Cell>
            <Table.Cell>{delivery.attempts}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{delivery.response_status ?? "—"}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{formatDateTime(delivery.created_at)}</Table.Cell>
            <Table.Cell>
              <Modal.Root>
                <Modal.Trigger>
                  {#snippet children(open)}
                    <Button variant="outline" size="sm" onclick={open}>View</Button>
                  {/snippet}
                </Modal.Trigger>
                <Modal.Content>
                  {#snippet children(close)}
                    <Modal.Header>
                      <Modal.Title>{delivery.event}</Modal.Title>
                      <Modal.Description>Delivery payload sent to the webhook endpoint.</Modal.Description>
                    </Modal.Header>
                    <pre
                      class="mt-4 max-h-80 overflow-auto rounded-md bg-secondary p-3 text-xs text-foreground">{JSON.stringify(
                        delivery.payload,
                        null,
                        2,
                      )}</pre>
                    <Modal.Footer>
                      <Button type="button" variant="outline" onclick={close}>Close</Button>
                    </Modal.Footer>
                  {/snippet}
                </Modal.Content>
              </Modal.Root>
            </Table.Cell>
          </Table.Row>
        {/each}
      {/if}
    </Table.Body>
  </Table.Root>

  {#if currentPage > 1 || data.deliveries.has_more}
    <div class="flex items-center justify-end gap-2">
      <Button variant="outline" size="sm" href={currentPage > 1 ? `?page=${currentPage - 1}` : undefined} disabled={currentPage <= 1}>
        Previous
      </Button>
      <Button
        variant="outline"
        size="sm"
        href={data.deliveries.has_more ? `?page=${currentPage + 1}` : undefined}
        disabled={!data.deliveries.has_more}
      >
        Next
      </Button>
    </div>
  {/if}
</div>

<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import * as Table from "$lib/components/ui/table";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { Badge } from "$lib/components/ui/badge";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import { Tooltip } from "$lib/components/ui/tooltip";
  import DeleteWebhookButton from "$lib/components/delete-webhook-button.svelte";
  import { createFormSubmit } from "$lib/form-submit.svelte";
  import { formatDate } from "$lib/format";
  import type { PageProps } from "./$types";

  const EVENTS = [
    { value: "link.created", label: "Link created" },
    { value: "link.updated", label: "Link updated" },
    { value: "link.deleted", label: "Link deleted" },
    { value: "link.clicked", label: "Link clicked" },
  ];

  let { data, form }: PageProps = $props();

  let revealedSecret = $state<string | null>(null);
  $effect(() => {
    if (form?.secret) revealedSecret = form.secret;
  });
</script>

<div class="flex flex-col gap-4 p-6">
  {#if revealedSecret}
    <Alert token="success" title="Webhook secret" onDismiss={() => (revealedSecret = null)}>
      Copy this now — it won't be shown again: <code class="font-mono">{revealedSecret}</code>
    </Alert>
  {/if}

  <div class="flex items-center justify-end">
    <Modal.Root>
      <Modal.Trigger>
        {#snippet children(open)}
          <Button onclick={open}>New webhook</Button>
        {/snippet}
      </Modal.Trigger>
      <Modal.Content>
        {#snippet children(close)}
          <Modal.Header>
            <Modal.Title>New webhook</Modal.Title>
            <Modal.Description>Get notified when events happen in this project.</Modal.Description>
          </Modal.Header>

          {#if form?.error}
            <Alert token="destructive" title="Couldn't create webhook" class="mt-4">{form.error}</Alert>
          {/if}

          {@const formSubmit = createFormSubmit(close)}
          <form class="mt-4 flex flex-col gap-4" method="POST" action="?/create" use:enhance={formSubmit.submit}>
            <Field label="URL">
              <Input type="url" name="url" placeholder="https://example.com/webhook" required />
            </Field>
            <Field label="Events">
              <div class="flex flex-col gap-2">
                {#each EVENTS as event (event.value)}
                  <label class="flex items-center gap-2 text-sm text-foreground">
                    <Checkbox name="events" value={event.value} />
                    {event.label}
                  </label>
                {/each}
              </div>
            </Field>
            <Modal.Footer>
              <Button type="button" variant="outline" onclick={close}>Cancel</Button>
              <Button type="submit" disabled={formSubmit.submitting}>
                {formSubmit.submitting ? "Creating…" : "Create"}
              </Button>
            </Modal.Footer>
          </form>
        {/snippet}
      </Modal.Content>
    </Modal.Root>
  </div>

  <Table.Root>
    <Table.Header>
      <Table.Row>
        <Table.Head>URL</Table.Head>
        <Table.Head>Events</Table.Head>
        <Table.Head>Status</Table.Head>
        <Table.Head>Created</Table.Head>
        <Table.Head>Actions</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#if data.webhooks.length === 0}
        <Table.Row>
          <Table.Cell colspan={5} class="text-center text-muted-foreground">No webhooks yet.</Table.Cell>
        </Table.Row>
      {:else}
        {#each data.webhooks as webhook (webhook.id)}
          <Table.Row>
            <Table.Cell class="max-w-xs truncate font-medium text-foreground">{webhook.url}</Table.Cell>
            <Table.Cell>
              <div class="flex flex-wrap gap-1">
                {#each webhook.events as event (event)}
                  <Badge token="neutral">{event}</Badge>
                {/each}
              </div>
            </Table.Cell>
            <Table.Cell>
              <Badge token={webhook.enabled ? "success" : "neutral"}>{webhook.enabled ? "Enabled" : "Disabled"}</Badge>
            </Table.Cell>
            <Table.Cell class="text-muted-foreground">{formatDate(webhook.created_at)}</Table.Cell>
            <Table.Cell>
              <div class="flex items-center gap-1">
                <Tooltip label="View deliveries">
                  <a
                    href="/p/{data.project.slug}/webhooks/{webhook.id}/deliveries"
                    aria-label="View deliveries"
                    class="rounded-md p-1.5 text-muted-foreground hover:bg-accent hover:text-foreground"
                  >
                    <svg
                      class="size-4"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <line x1="18" y1="20" x2="18" y2="10" />
                      <line x1="12" y1="20" x2="12" y2="4" />
                      <line x1="6" y1="20" x2="6" y2="14" />
                    </svg>
                  </a>
                </Tooltip>
                <DeleteWebhookButton id={webhook.id} />
              </div>
            </Table.Cell>
          </Table.Row>
        {/each}
      {/if}
    </Table.Body>
  </Table.Root>
</div>

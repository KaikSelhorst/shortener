<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { Tooltip } from "$lib/components/ui/tooltip";

  interface Link {
    short_code: string;
    original_url: string;
    title: string | null;
    description: string | null;
    og_image: string | null;
    expires_at: string | null;
    max_clicks: number | null;
  }

  interface Props {
    link: Link;
  }

  let { link }: Props = $props();

  let submitting = $state(false);
  let error = $state<string | null>(null);

  // datetime-local inputs carry no timezone — render the stored UTC instant
  // using the browser's own local time so the field shows what the user expects.
  function toDatetimeLocal(value: string | null) {
    if (!value) return "";
    const date = new Date(value);
    const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000);
    return local.toISOString().slice(0, 16);
  }
</script>

<Modal.Root>
  <Modal.Trigger>
    {#snippet children(open)}
      <Tooltip label="Edit link">
        <button
          type="button"
          onclick={() => {
            error = null;
            open();
          }}
          aria-label="Edit link"
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
            <path d="M12 20h9" />
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
          </svg>
        </button>
      </Tooltip>
    {/snippet}
  </Modal.Trigger>
  <Modal.Content>
    {#snippet children(close)}
      <Modal.Header>
        <Modal.Title>Edit link</Modal.Title>
        <Modal.Description>The short code can't be changed here.</Modal.Description>
      </Modal.Header>

      {#if error}
        <Alert token="destructive" title="Couldn't update link" class="mt-4">{error}</Alert>
      {/if}

      <form
        class="mt-4 flex flex-col gap-4"
        method="POST"
        action="?/update"
        use:enhance={({ formData }) => {
          const expiresAt = formData.get("expires_at");
          if (typeof expiresAt === "string" && expiresAt) {
            formData.set("expires_at", new Date(expiresAt).toISOString());
          }
          submitting = true;
          return async ({ update, result }) => {
            await update();
            submitting = false;
            if (result.type === "success") {
              close();
            } else if (result.type === "failure") {
              error = (result.data as { error?: string } | undefined)?.error ?? "Something went wrong.";
            }
          };
        }}
      >
        <input type="hidden" name="code" value={link.short_code} />
        <input type="hidden" name="description" value={link.description ?? ""} />
        <input type="hidden" name="og_image" value={link.og_image ?? ""} />

        <Field label="URL">
          <Input type="url" name="url" value={link.original_url} required />
        </Field>
        <Field label="Title" description="Optional.">
          <Input name="title" value={link.title ?? ""} />
        </Field>
        <Field label="Max clicks" description="Optional. The link stops working after this many clicks.">
          <Input type="number" name="max_clicks" min="1" value={link.max_clicks ?? ""} placeholder="Unlimited" />
        </Field>
        <Field label="Expires at" description="Optional, in your local time.">
          <Input type="datetime-local" name="expires_at" value={toDatetimeLocal(link.expires_at)} />
        </Field>
        <Modal.Footer>
          <Button type="button" variant="outline" onclick={close}>Cancel</Button>
          <Button type="submit" disabled={submitting}>
            {submitting ? "Saving…" : "Save"}
          </Button>
        </Modal.Footer>
      </form>
    {/snippet}
  </Modal.Content>
</Modal.Root>

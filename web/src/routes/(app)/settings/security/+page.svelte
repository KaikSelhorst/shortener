<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Card from "$lib/components/ui/card";
  import * as Modal from "$lib/components/ui/modal";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { Badge } from "$lib/components/ui/badge";
  import { createFormSubmit } from "$lib/form-submit.svelte";
  import type { PageProps } from "./$types";

  let { data, form }: PageProps = $props();

  let pendingSetup = $state<{ uri: string; secret: string } | null>(null);
  $effect(() => {
    if (form?.setup) pendingSetup = { uri: form.uri, secret: form.secret };
  });

  const setupSubmit = createFormSubmit();
  const confirmSubmit = createFormSubmit(() => (pendingSetup = null));

  function formatSecret(secret: string) {
    return secret.match(/.{1,4}/g)?.join(" ") ?? secret;
  }
</script>

<div class="flex flex-col gap-4 p-6">
  <Card.Root>
    <Card.Header>
      <Card.Title>Account</Card.Title>
    </Card.Header>
    <Card.Content>
      <div class="flex flex-col gap-1">
        <span class="text-xs font-medium text-muted-foreground">Email</span>
        <span class="text-sm text-foreground">{data.me.email}</span>
      </div>
    </Card.Content>
  </Card.Root>

  <Card.Root>
    <Card.Header>
      <Card.Title>Two-factor authentication</Card.Title>
      <Card.Description>Require a code from an authenticator app when signing in.</Card.Description>
    </Card.Header>
    <Card.Content>
      {#if data.me.totp_enabled}
        <div class="flex items-center justify-between gap-4">
          <Badge token="success">Enabled</Badge>
          <Modal.Root>
            <Modal.Trigger>
              {#snippet children(open)}
                <Button variant="outline" onclick={open}>Disable</Button>
              {/snippet}
            </Modal.Trigger>
            <Modal.Content>
              {#snippet children(close)}
                <Modal.Header>
                  <Modal.Title>Disable two-factor authentication?</Modal.Title>
                  <Modal.Description>Enter a current code from your authenticator app to confirm.</Modal.Description>
                </Modal.Header>

                {#if form?.error}
                  <Alert token="destructive" title="Couldn't disable 2FA" class="mt-4">{form.error}</Alert>
                {/if}

                {@const disableSubmit = createFormSubmit(close)}
                <form class="mt-4" method="POST" action="?/disable" use:enhance={disableSubmit.submit}>
                  <Field label="Code">
                    <Input
                      type="text"
                      name="code"
                      inputmode="numeric"
                      autocomplete="one-time-code"
                      maxlength={6}
                      placeholder="000000"
                      required
                    />
                  </Field>
                  <Modal.Footer>
                    <Button type="button" variant="outline" onclick={close}>Cancel</Button>
                    <Button type="submit" variant="destructive" disabled={disableSubmit.submitting}>
                      {disableSubmit.submitting ? "Disabling…" : "Disable"}
                    </Button>
                  </Modal.Footer>
                </form>
              {/snippet}
            </Modal.Content>
          </Modal.Root>
        </div>
      {:else if pendingSetup}
        {#if form?.error}
          <Alert token="destructive" title="Couldn't confirm code" class="mb-4">{form.error}</Alert>
        {/if}
        <div class="flex flex-col gap-3">
          <p class="text-sm text-muted-foreground">
            Add this key to your authenticator app, then enter the 6-digit code it generates.
          </p>
          <code class="block rounded-md bg-secondary p-3 text-center font-mono text-sm tracking-widest text-foreground">
            {formatSecret(pendingSetup.secret)}
          </code>
          <form class="flex items-end gap-2" method="POST" action="?/confirm" use:enhance={confirmSubmit.submit}>
            <Field label="Code" class="flex-1">
              <Input
                type="text"
                name="code"
                inputmode="numeric"
                autocomplete="one-time-code"
                maxlength={6}
                placeholder="000000"
                required
              />
            </Field>
            <Button type="submit" disabled={confirmSubmit.submitting}>
              {confirmSubmit.submitting ? "Confirming…" : "Confirm"}
            </Button>
            <Button type="button" variant="outline" onclick={() => (pendingSetup = null)}>Cancel</Button>
          </form>
        </div>
      {:else}
        {#if form?.error}
          <Alert token="destructive" title="Couldn't start 2FA setup" class="mb-4">{form.error}</Alert>
        {/if}
        <div class="flex items-center justify-between gap-4">
          <p class="text-sm text-muted-foreground">Two-factor authentication is not enabled.</p>
          <form method="POST" action="?/setup" use:enhance={setupSubmit.submit}>
            <Button type="submit" disabled={setupSubmit.submitting}>
              {setupSubmit.submitting ? "Starting…" : "Enable 2FA"}
            </Button>
          </form>
        </div>
      {/if}
    </Card.Content>
  </Card.Root>
</div>

<script lang="ts">
  import { enhance } from "$app/forms";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import type { PageProps } from "./$types";

  let { form }: PageProps = $props();

  let submitting = $state(false);
</script>

<svelte:head>
  <title>Two-factor authentication — Shortener</title>
</svelte:head>

<div class="flex flex-col gap-1 text-center">
  <h1 class="text-2xl font-semibold text-foreground">Two-factor authentication</h1>
  <p class="text-sm text-muted-foreground">Enter the 6-digit code from your authenticator app.</p>
</div>

{#if form?.error}
  <Alert token="destructive" title="Couldn't verify code">{form.error}</Alert>
{/if}

<form
  class="flex flex-col gap-4"
  method="POST"
  use:enhance={() => {
    submitting = true;
    return async ({ update }) => {
      await update();
      submitting = false;
    };
  }}
>
  <Field label="Code">
    <Input
      type="text"
      name="code"
      inputmode="numeric"
      autocomplete="one-time-code"
      maxlength={6}
      placeholder="000000"
      class="text-center text-lg tracking-[0.5em]"
      required
    />
  </Field>
  <Button type="submit" class="mt-2 w-full" disabled={submitting}>
    {submitting ? "Verifying…" : "Verify"}
  </Button>
</form>

<p class="text-center text-sm text-muted-foreground">
  <a href="/login" class="font-medium text-foreground hover:underline">Back to sign in</a>
</p>

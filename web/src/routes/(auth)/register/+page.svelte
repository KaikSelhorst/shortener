<script lang="ts">
  import { enhance } from "$app/forms";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { PasswordInput } from "$lib/components/ui/password-input";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import type { PageProps } from "./$types";

  let { form }: PageProps = $props();

  let submitting = $state(false);
</script>

<div class="flex flex-col gap-1 text-center">
  <h1 class="text-2xl font-semibold text-foreground">Create an account</h1>
  <p class="text-sm text-muted-foreground">Enter your email and choose a password to get started.</p>
</div>

{#if form?.error}
  <Alert token="destructive" title="Couldn't create account">{form.error}</Alert>
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
  <Field label="Email">
    <Input type="email" name="email" placeholder="you@example.com" autocomplete="email" required />
  </Field>
  <Field label="Password" description="8–72 characters.">
    <PasswordInput name="password" placeholder="••••••••" autocomplete="new-password" minlength={8} maxlength={72} required />
  </Field>
  <Button type="submit" class="mt-2 w-full" disabled={submitting}>
    {submitting ? "Creating account…" : "Create account"}
  </Button>
</form>

<p class="text-center text-sm text-muted-foreground">
  Already have an account? <a href="/login" class="font-medium text-foreground hover:underline">Sign in</a>
</p>

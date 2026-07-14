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
  <h1 class="text-2xl font-semibold text-foreground">Sign in</h1>
  <p class="text-sm text-muted-foreground">Enter your email and password to continue.</p>
</div>

{#if form?.error}
  <Alert token="destructive" title="Couldn't sign in">{form.error}</Alert>
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
  <Field label="Password">
    <PasswordInput name="password" placeholder="••••••••" autocomplete="current-password" required />
  </Field>
  <Button type="submit" class="mt-2 w-full" disabled={submitting}>
    {submitting ? "Signing in…" : "Sign in"}
  </Button>
</form>

<p class="text-center text-sm text-muted-foreground">
  Don't have an account? <a href="/register" class="font-medium text-foreground hover:underline">Sign up</a>
</p>

<script lang="ts">
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { PasswordInput } from "$lib/components/ui/password-input";
  import { Button } from "$lib/components/ui/button";

  function onsubmit(e: SubmitEvent) {
    e.preventDefault();
  }
</script>

<div class="flex flex-col gap-1 text-center">
  <h1 class="text-2xl font-semibold text-foreground">Create an account</h1>
  <p class="text-sm text-muted-foreground">Enter your email and choose a password to get started.</p>
</div>

<form class="flex flex-col gap-4" {onsubmit}>
  <Field label="Email">
    <Input type="email" name="email" placeholder="you@example.com" autocomplete="email" required />
  </Field>
  <Field label="Password" description="8–72 characters.">
    <PasswordInput name="password" placeholder="••••••••" autocomplete="new-password" minlength={8} maxlength={72} required />
  </Field>
  <Button type="submit" class="mt-2 w-full">Create account</Button>
</form>

<p class="text-center text-sm text-muted-foreground">
  Already have an account? <a href="/login" class="font-medium text-foreground hover:underline">Sign in</a>
</p>

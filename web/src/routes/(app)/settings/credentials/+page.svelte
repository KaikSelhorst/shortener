<script lang="ts">
  import { enhance } from "$app/forms";
  import * as Modal from "$lib/components/ui/modal";
  import * as Table from "$lib/components/ui/table";
  import { Field } from "$lib/components/ui/field";
  import { Input } from "$lib/components/ui/input";
  import { Select } from "$lib/components/ui/select";
  import { Button } from "$lib/components/ui/button";
  import { Alert } from "$lib/components/ui/alert";
  import { Badge } from "$lib/components/ui/badge";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import DeleteCredentialButton from "$lib/components/delete-credential-button.svelte";
  import { createFormSubmit } from "$lib/form-submit.svelte";
  import { formatDate } from "$lib/format";
  import type { PageProps } from "./$types";

  const SCOPE_GROUPS = [
    { label: "Links", scopes: ["links:create", "links:read", "links:update", "links:delete"] },
    { label: "Projects", scopes: ["projects:create", "projects:read", "projects:update", "projects:delete"] },
    { label: "Webhooks", scopes: ["webhooks:read", "webhooks:create", "webhooks:delete"] },
  ];

  let { data, form }: PageProps = $props();

  let revealedToken = $state<string | null>(null);
  $effect(() => {
    if (form?.token) revealedToken = form.token;
  });

  function projectName(projectId: number | null) {
    if (projectId === null) return "All projects";
    return data.projects.find((p) => p.id === projectId)?.name ?? `#${projectId}`;
  }
</script>

<div class="flex flex-col gap-4 p-6">
  {#if revealedToken}
    <Alert token="success" title="Credential" onDismiss={() => (revealedToken = null)}>
      Copy this now — it won't be shown again: <code class="font-mono">{revealedToken}</code>
    </Alert>
  {/if}

  <div class="flex items-center justify-end">
    <Modal.Root>
      <Modal.Trigger>
        {#snippet children(open)}
          <Button onclick={open}>New credential</Button>
        {/snippet}
      </Modal.Trigger>
      <Modal.Content>
        {#snippet children(close)}
          <Modal.Header>
            <Modal.Title>New credential</Modal.Title>
            <Modal.Description>Use this to authenticate requests to the API instead of a user session.</Modal.Description>
          </Modal.Header>

          {#if form?.error}
            <Alert token="destructive" title="Couldn't create credential" class="mt-4">{form.error}</Alert>
          {/if}

          {@const formSubmit = createFormSubmit(close)}
          <form class="mt-4 flex flex-col gap-4" method="POST" action="?/create" use:enhance={formSubmit.submit}>
            <Field label="Name">
              <Input name="name" placeholder="CI deploy key" required />
            </Field>
            <Field label="Project" description="Restrict this key to one project, or leave unrestricted.">
              <Select name="project_id">
                <option value="">All projects</option>
                {#each data.projects as project (project.id)}
                  <option value={project.id}>{project.name}</option>
                {/each}
              </Select>
            </Field>
            <Field label="Scopes">
              <div class="flex flex-col gap-2">
                <label class="flex items-center gap-2 text-sm text-foreground">
                  <Checkbox name="scopes" value="*" />
                  Full access (*)
                </label>
              </div>
            </Field>
            <div class="grid grid-cols-3 gap-4">
              {#each SCOPE_GROUPS as group (group.label)}
                <div class="flex flex-col gap-2">
                  <span class="text-xs font-medium text-muted-foreground">{group.label}</span>
                  {#each group.scopes as scope (scope)}
                    <label class="flex items-center gap-2 text-sm text-foreground">
                      <Checkbox name="scopes" value={scope} />
                      {scope.split(":")[1]}
                    </label>
                  {/each}
                </div>
              {/each}
            </div>
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
        <Table.Head>Name</Table.Head>
        <Table.Head>Key</Table.Head>
        <Table.Head>Scopes</Table.Head>
        <Table.Head>Project</Table.Head>
        <Table.Head>Last used</Table.Head>
        <Table.Head>Created</Table.Head>
        <Table.Head>Actions</Table.Head>
      </Table.Row>
    </Table.Header>
    <Table.Body>
      {#if data.credentials.length === 0}
        <Table.Row>
          <Table.Cell colspan={7} class="text-center text-muted-foreground">No credentials yet.</Table.Cell>
        </Table.Row>
      {:else}
        {#each data.credentials as key (key.id)}
          <Table.Row>
            <Table.Cell class="font-medium text-foreground">{key.name}</Table.Cell>
            <Table.Cell class="font-mono text-muted-foreground">{key.key_prefix}…</Table.Cell>
            <Table.Cell>
              <div class="flex flex-wrap gap-1">
                {#each key.scopes as scope (scope)}
                  <Badge token="neutral">{scope}</Badge>
                {/each}
              </div>
            </Table.Cell>
            <Table.Cell class="text-muted-foreground">{projectName(key.project_id)}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{key.last_used_at ? formatDate(key.last_used_at) : "Never"}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{formatDate(key.created_at)}</Table.Cell>
            <Table.Cell>
              <DeleteCredentialButton id={key.id} />
            </Table.Cell>
          </Table.Row>
        {/each}
      {/if}
    </Table.Body>
  </Table.Root>
</div>

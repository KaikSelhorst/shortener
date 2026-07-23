<script lang="ts">
  import { page } from "$app/state";
  import { Breadcrumbs } from "$lib/components/ui/breadcrumbs";
  import * as SidebarNav from "$lib/components/ui/sidebar-nav";
  import AccountMenu from "$lib/components/account-menu.svelte";
  import type { LayoutProps } from "./$types";

  let { children }: LayoutProps = $props();

  const credentialsHref = "/settings/credentials";
  const securityHref = "/settings/security";

  const SECTION_LABELS: Record<string, string> = { credentials: "Credentials", security: "Security" };
  const section = $derived(page.url.pathname.split("/")[2] ?? "");
  const sectionLabel = $derived(SECTION_LABELS[section] ?? "");

  const crumbs = $derived([
    { label: "Home", href: "/p" },
    ...(sectionLabel ? [{ label: sectionLabel }] : [{ label: "Settings" }]),
  ]);
</script>

<div class="flex h-screen">
  <aside class="flex w-56 shrink-0 flex-col overflow-y-auto border-r border-border bg-card">
    <div class="flex flex-1 flex-col gap-3 p-3">
      <div class="flex flex-col gap-1">
        <a
          href="/p"
          class="flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground"
        >
          <svg
            class="size-4 shrink-0"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M19 12H5" />
            <path d="m12 19-7-7 7-7" />
          </svg>
          Back
        </a>

        <SidebarNav.Root>
          <SidebarNav.Item href={credentialsHref} active={page.url.pathname.startsWith(credentialsHref)}>
            {#snippet icon()}
              <svg
                class="size-4 shrink-0"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            {/snippet}
            Credentials
          </SidebarNav.Item>

          <SidebarNav.Item href={securityHref} active={page.url.pathname.startsWith(securityHref)}>
            {#snippet icon()}
              <svg
                class="size-4 shrink-0"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
              </svg>
            {/snippet}
            Security
          </SidebarNav.Item>
        </SidebarNav.Root>
      </div>

      <div class="mt-auto">
        <AccountMenu />
      </div>
    </div>
  </aside>

  <div class="flex min-h-0 flex-1 flex-col">
    <header class="flex h-14 shrink-0 items-center border-b border-border px-6">
      <Breadcrumbs items={crumbs} />
    </header>
    <main class="min-h-0 flex-1 overflow-y-auto">
      {@render children()}
    </main>
  </div>
</div>

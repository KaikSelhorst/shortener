<script lang="ts">
  import { page } from "$app/state";
  import { Breadcrumbs } from "$lib/components/ui/breadcrumbs";
  import * as SidebarNav from "$lib/components/ui/sidebar-nav";
  import ProjectSwitcher from "$lib/components/project-switcher.svelte";
  import LogoutButton from "$lib/components/logout-button.svelte";
  import type { LayoutProps } from "./$types";

  let { data, children }: LayoutProps = $props();

  const linksHref = $derived(`/p/${data.project.slug}/links`);
  const analyticsHref = $derived(`/p/${data.project.slug}/analytics`);
  const webhooksHref = $derived(`/p/${data.project.slug}/webhooks`);

  const section = $derived(page.url.pathname.split("/")[3] ?? "");
  const sectionLabel = $derived(section ? section.charAt(0).toUpperCase() + section.slice(1) : "");

  const crumbs = $derived([
    { label: "Home", href: "/p" },
    { label: data.project.name, href: `/p/${data.project.slug}/links` },
    ...(sectionLabel ? [{ label: sectionLabel }] : []),
  ]);
</script>

<div class="flex h-screen">
  <aside class="flex w-56 shrink-0 flex-col overflow-y-auto border-r border-border bg-card">
    <ProjectSwitcher projects={data.projects} active={data.project} />

    <div class="flex flex-1 flex-col gap-3 p-3">
      <SidebarNav.Root>
        <SidebarNav.Item href={linksHref} active={page.url.pathname.startsWith(linksHref)}>
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
              <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
              <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
            </svg>
          {/snippet}
          Links
        </SidebarNav.Item>

        <SidebarNav.Item href={analyticsHref} active={page.url.pathname.startsWith(analyticsHref)}>
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
              <line x1="18" y1="20" x2="18" y2="10" />
              <line x1="12" y1="20" x2="12" y2="4" />
              <line x1="6" y1="20" x2="6" y2="14" />
            </svg>
          {/snippet}
          Analytics
        </SidebarNav.Item>

        <SidebarNav.Item href={webhooksHref} active={page.url.pathname.startsWith(webhooksHref)}>
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
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
            </svg>
          {/snippet}
          Webhooks
        </SidebarNav.Item>
      </SidebarNav.Root>

      <div class="mt-auto">
        <LogoutButton />
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

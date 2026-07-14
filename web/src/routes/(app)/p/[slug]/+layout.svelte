<script lang="ts">
  import { page } from "$app/state";
  import { Breadcrumbs } from "$lib/components/ui/breadcrumbs";
  import * as SidebarNav from "$lib/components/ui/sidebar-nav";
  import ProjectSwitcher from "$lib/components/project-switcher.svelte";
  import type { LayoutProps } from "./$types";

  let { data, children }: LayoutProps = $props();

  const navItems = $derived([
    { label: "Links", href: `/p/${data.project.slug}/links` },
    { label: "Analytics", href: `/p/${data.project.slug}/analytics` },
    { label: "Webhooks", href: `/p/${data.project.slug}/webhooks` },
  ]);

  const section = $derived(page.url.pathname.split("/")[3] ?? "");
  const sectionLabel = $derived(section ? section.charAt(0).toUpperCase() + section.slice(1) : "");

  const crumbs = $derived([
    { label: "Home", href: "/p" },
    { label: data.project.name, href: `/p/${data.project.slug}/links` },
    ...(sectionLabel ? [{ label: sectionLabel }] : []),
  ]);
</script>

<div class="flex min-h-screen">
  <aside class="flex w-56 shrink-0 flex-col gap-4 border-r border-border bg-card p-4">
    <ProjectSwitcher projects={data.projects} active={data.project} />

    <SidebarNav.Root>
      {#each navItems as item (item.href)}
        <SidebarNav.Item href={item.href} active={page.url.pathname.startsWith(item.href)}>
          {item.label}
        </SidebarNav.Item>
      {/each}
    </SidebarNav.Root>

    <form method="POST" action="/logout" class="mt-auto">
      <button
        type="submit"
        class="w-full rounded-md px-3 py-2 text-left text-sm font-medium text-muted-foreground hover:bg-accent hover:text-foreground"
      >
        Logout
      </button>
    </form>
  </aside>

  <div class="flex flex-1 flex-col">
    <header class="border-b border-border px-6 py-4">
      <Breadcrumbs items={crumbs} />
    </header>
    <main class="flex-1">
      {@render children()}
    </main>
  </div>
</div>

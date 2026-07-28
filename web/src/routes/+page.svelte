<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import * as Table from "$lib/components/ui/table";
  import AsciiPlasma from "$lib/components/ascii-plasma.svelte";
  import TimeSeriesChart from "$lib/components/charts/time-series-chart.svelte";

  const SNIPPET = `curl https://docut.xyz/projects/marketing/links \\
  -H "Authorization: Bearer sk_live_..." \\
  -H "content-type: application/json" \\
  -d '{"url":"https://example.com/launch","custom_code":"launch"}'`;

  let copied = $state(false);

  async function copySnippet() {
    await navigator.clipboard.writeText(SNIPPET);
    copied = true;
    setTimeout(() => (copied = false), 1500);
  }

  const MOCK_SERIES = [
    { date: "2026-07-01", count: 420 },
    { date: "2026-07-02", count: 512 },
    { date: "2026-07-03", count: 380 },
    { date: "2026-07-04", count: 610 },
    { date: "2026-07-05", count: 705 },
    { date: "2026-07-06", count: 540 },
    { date: "2026-07-07", count: 890 },
    { date: "2026-07-08", count: 760 },
    { date: "2026-07-09", count: 920 },
    { date: "2026-07-10", count: 1040 },
    { date: "2026-07-11", count: 880 },
    { date: "2026-07-12", count: 1150 },
    { date: "2026-07-13", count: 1320 },
    { date: "2026-07-14", count: 1480 },
  ];

  const MOCK_LINKS = [
    { short: "docut.xyz/launch", clicks: "12,480" },
    { short: "docut.xyz/promo-q3", clicks: "8,210" },
    { short: "docut.xyz/docs", clicks: "3,190" },
  ];
</script>

<svelte:head>
  <title>Docut</title>
</svelte:head>

<div class="flex min-h-screen flex-col bg-background">
  <header class="border-b border-border">
    <div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
      <a href="/" class="text-sm font-semibold tracking-tight text-foreground">
        Docut
      </a>
      <nav class="flex items-center gap-2">
        <Button href="/login" variant="ghost" size="sm">Sign in</Button>
        <Button href="/register" size="sm">Get started</Button>
      </nav>
    </div>
  </header>

  <section class="relative overflow-hidden border-b border-border">
    <div class="absolute inset-0 -z-10">
      <AsciiPlasma />
    </div>
    <div class="mx-auto max-w-6xl px-6 py-20 lg:py-28">
      <div class="mx-auto flex max-w-2xl flex-col items-center gap-6 text-center">
        <h1 class="text-4xl leading-tight font-semibold text-foreground md:text-5xl">Short links. Full picture.</h1>
        <p class="text-base text-muted-foreground md:text-lg">
          Shorten, organize, and track links per project — with click analytics, webhooks, and a scoped API for
          everything you automate.
        </p>
        <div class="flex flex-wrap items-center justify-center gap-3">
          <Button href="/register" size="lg">Get started free</Button>
          <Button href="/login" variant="outline" size="lg">Sign in</Button>
        </div>
      </div>

      <div class="mx-auto mt-14 max-w-4xl overflow-hidden rounded-lg border border-border bg-card shadow-2xl">
        <div class="flex items-center gap-1.5 border-b border-border px-4 py-3">
          <span class="size-2.5 rounded-full bg-border"></span>
          <span class="size-2.5 rounded-full bg-border"></span>
          <span class="size-2.5 rounded-full bg-border"></span>
          <span class="ml-3 truncate rounded-md bg-secondary px-2 py-1 text-xs text-muted-foreground">
            docut.xyz/p/marketing/analytics
          </span>
        </div>
        <div class="flex flex-col gap-5 p-5">
          <div class="grid grid-cols-3 gap-4">
            <div>
              <p class="text-xs text-muted-foreground">Total clicks</p>
              <p class="text-lg font-semibold text-foreground">48,213</p>
            </div>
            <div>
              <p class="text-xs text-muted-foreground">Unique visitors</p>
              <p class="text-lg font-semibold text-foreground">31,904</p>
            </div>
            <div>
              <p class="text-xs text-muted-foreground">Active links</p>
              <p class="text-lg font-semibold text-foreground">142</p>
            </div>
          </div>

          <TimeSeriesChart data={MOCK_SERIES} />

          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.Head>Short link</Table.Head>
                <Table.Head>Clicks</Table.Head>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {#each MOCK_LINKS as link (link.short)}
                <Table.Row>
                  <Table.Cell class="font-medium text-foreground">{link.short}</Table.Cell>
                  <Table.Cell class="text-muted-foreground">{link.clicks}</Table.Cell>
                </Table.Row>
              {/each}
            </Table.Body>
          </Table.Root>
        </div>
      </div>
    </div>
  </section>

  <section class="border-b border-border">
    <div class="mx-auto max-w-6xl px-6 py-16">
      <div class="grid gap-10 md:grid-cols-3">
        <div>
          <span class="text-xs font-medium text-muted-foreground">01</span>
          <h3 class="mt-2 text-base font-medium text-foreground">Create a project</h3>
          <p class="mt-1 text-sm text-muted-foreground">
            Group links by product, campaign, or team — each with its own credentials and webhooks.
          </p>
        </div>
        <div>
          <span class="text-xs font-medium text-muted-foreground">02</span>
          <h3 class="mt-2 text-base font-medium text-foreground">Shorten & share</h3>
          <p class="mt-1 text-sm text-muted-foreground">
            Custom codes, expiration dates, and max click limits — set once at creation.
          </p>
        </div>
        <div>
          <span class="text-xs font-medium text-muted-foreground">03</span>
          <h3 class="mt-2 text-base font-medium text-foreground">Track everything</h3>
          <p class="mt-1 text-sm text-muted-foreground">
            Device, browser, referrer, and unique-visitor breakdowns for every click.
          </p>
        </div>
      </div>
    </div>
  </section>

  <section class="border-b border-border">
    <div class="mx-auto max-w-6xl px-6 py-16">
      <h2 class="text-xl font-medium text-foreground">Everything you need</h2>
      <div class="mt-8 grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div class="flex gap-3">
          <svg class="mt-0.5 size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" x2="12" y1="20" y2="10" />
            <line x1="18" x2="18" y1="20" y2="4" />
            <line x1="6" x2="6" y1="20" y2="16" />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-foreground">Analytics</h3>
            <p class="mt-1 text-sm text-muted-foreground">Device, browser, referrer, and unique-visitor breakdowns per link.</p>
          </div>
        </div>
        <div class="flex gap-3">
          <svg class="mt-0.5 size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path
              d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z"
            />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-foreground">Webhooks</h3>
            <p class="mt-1 text-sm text-muted-foreground">Real-time delivery for link creation, updates, deletes, and clicks.</p>
          </div>
        </div>
        <div class="flex gap-3">
          <svg class="mt-0.5 size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-foreground">Credentials</h3>
            <p class="mt-1 text-sm text-muted-foreground">Scoped API keys per project, for anything you want to automate.</p>
          </div>
        </div>
        <div class="flex gap-3">
          <svg class="mt-0.5 size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-foreground">Link controls</h3>
            <p class="mt-1 text-sm text-muted-foreground">Expiration dates and max click limits, enforced at redirect time.</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <section class="border-b border-border">
    <div class="mx-auto max-w-6xl px-6 py-16">
      <div class="grid gap-10 lg:grid-cols-2 lg:items-center">
        <div class="flex flex-col gap-4">
          <h2 class="text-xl font-medium text-foreground">Built for developers</h2>
          <p class="text-sm text-muted-foreground">
            Every project ships with scoped API keys and real-time webhooks — automate link creation, react to
            clicks, and build on top of Docut without touching the dashboard.
          </p>
          <ul class="mt-2 flex flex-col gap-2 text-sm text-muted-foreground">
            <li class="flex items-center gap-2">
              <svg class="size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12" />
              </svg>
              Scoped credentials per project
            </li>
            <li class="flex items-center gap-2">
              <svg class="size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12" />
              </svg>
              Webhooks for created, updated, deleted, and clicked
            </li>
            <li class="flex items-center gap-2">
              <svg class="size-4 shrink-0 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12" />
              </svg>
              Full OpenAPI spec at <code class="text-foreground">/openapi.json</code>
            </li>
          </ul>
        </div>

        <div class="overflow-hidden rounded-lg border border-border bg-card shadow-2xl">
          <div class="flex items-center gap-1.5 border-b border-border px-4 py-3">
            <span class="size-2.5 rounded-full bg-border"></span>
            <span class="size-2.5 rounded-full bg-border"></span>
            <span class="size-2.5 rounded-full bg-border"></span>
            <span class="ml-3 truncate rounded-md bg-secondary px-2 py-1 text-xs text-muted-foreground">zsh</span>
            <button
              type="button"
              onclick={copySnippet}
              class="ml-auto flex items-center gap-1.5 rounded-md px-2 py-1 text-xs text-muted-foreground hover:bg-accent hover:text-foreground"
            >
              {#if copied}
                <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12" />
                </svg>
                Copied
              {:else}
                <svg class="size-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="14" height="14" x="8" y="8" rx="2" ry="2" />
                  <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2" />
                </svg>
                Copy
              {/if}
            </button>
          </div>

          <pre class="overflow-x-auto p-6 font-mono text-xs leading-loose text-foreground"><code
            ><span class="text-muted-foreground">$</span> curl https://docut.xyz/projects/marketing/links \
    -H "Authorization: Bearer sk_live_..." \
    -H "content-type: application/json" \
    -d '<span class="text-lime">{"{"}</span>"url":"https://example.com/launch","custom_code":"launch"<span class="text-lime">{"}"}</span>'

<span class="text-muted-foreground"># 201 Created</span>
<span class="text-lime">{"{"}</span>"short_url":"https://docut.xyz/launch","total_clicks":0<span class="text-lime">{"}"}</span></code></pre>
        </div>
      </div>
    </div>
  </section>

  <section>
    <div class="mx-auto max-w-6xl px-6 py-16">
      <div class="flex flex-col items-center gap-4 rounded-lg border border-border bg-card p-10 text-center">
        <h2 class="text-xl font-medium text-foreground">Ready to get started?</h2>
        <p class="max-w-md text-sm text-muted-foreground">Create an account and shorten your first link in under a minute.</p>
        <Button href="/register" size="lg">Get started free</Button>
      </div>
    </div>
  </section>

  <footer class="mt-auto border-t border-border">
    <div class="mx-auto max-w-6xl px-6 py-3 text-center text-xs text-muted-foreground">
      © {new Date().getFullYear()} Docut
    </div>
  </footer>
</div>

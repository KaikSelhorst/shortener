<script lang="ts">
  import "../app.css";
  import { dev } from "$app/environment";
  import DevPerfMonitor from "$lib/components/dev-perf-monitor.svelte";

  let { children } = $props();
</script>

<svelte:head>
  <title>Shortener</title>
  <meta name="description" content="Fast, self-hosted link shortening with click analytics and webhooks." />
</svelte:head>

{@render children()}

{#if dev}
  <DevPerfMonitor />
{/if}

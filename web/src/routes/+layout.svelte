<script lang="ts">
  import "../app.css";

  let { children } = $props();
</script>

<svelte:head>
  <title>Shortener</title>
  <meta name="description" content="Fast, self-hosted link shortening with click analytics and webhooks." />
</svelte:head>

{@render children()}

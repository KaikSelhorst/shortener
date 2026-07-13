<script lang="ts">
	import '$lib/styles/theme.css';
	import favicon from '$lib/assets/favicon.svg';
	import FpsMeter from '$lib/components/dev/FpsMeter.svelte';
	import { dev } from '$app/environment';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
{@render children()}
{#if dev}<FpsMeter />{/if}

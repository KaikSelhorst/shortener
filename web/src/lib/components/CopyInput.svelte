<script lang="ts">
	import Button from './Button.svelte'
	import { onDestroy } from 'svelte'

	interface Props {
		value: string
	}

	let { value }: Props = $props()

	let copied = $state(false)
	let timeoutId: ReturnType<typeof setTimeout> | undefined

	onDestroy(() => clearTimeout(timeoutId))

	async function copy() {
		await navigator.clipboard.writeText(value)
		copied = true
		clearTimeout(timeoutId)
		timeoutId = setTimeout(() => (copied = false), 2000)
	}
</script>

<div class="flex items-center gap-2 border border-border bg-card px-3 py-2">
	<code class="flex-1 break-all font-mono text-[10px] text-foreground">{value}</code>
	<Button variant="outline" size="sm" onclick={copy} aria-label="Copy to clipboard">
		{copied ? 'copied!' : 'copy'}
	</Button>
</div>

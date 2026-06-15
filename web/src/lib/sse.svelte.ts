import { onMount } from 'svelte'

export function useSSE(getUrl: () => string, onMessage: (e: MessageEvent) => void) {
	let connected = $state(false)

	onMount(() => {
		const es = new EventSource(getUrl())
		es.onopen = () => { connected = true }
		es.onmessage = onMessage
		es.onerror = () => { connected = false }
		return () => es.close()
	})

	return {
		get connected() { return connected },
	}
}

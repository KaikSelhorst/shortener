import { onMount } from 'svelte'

const MAX_RETRIES = 5

export function useSSE(getUrl: () => string, onMessage: (e: MessageEvent) => void) {
	let connected = $state(false)

	onMount(() => {
		let es: EventSource
		let retries = 0
		let closed = false

		function connect() {
			if (closed) return
			es = new EventSource(getUrl())
			es.onopen = () => { connected = true; retries = 0 }
			es.onmessage = onMessage
			es.onerror = () => {
				connected = false
				es.close()
				if (closed || ++retries > MAX_RETRIES) return
				setTimeout(connect, Math.min(1000 * 2 ** retries, 30_000))
			}
		}

		connect()
		return () => { closed = true; es?.close() }
	})

	return {
		get connected() { return connected },
	}
}

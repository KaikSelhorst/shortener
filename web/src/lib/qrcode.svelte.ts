export function useQrcode() {
	let open = $state(false)
	let url = $state<string | null>(null)
	let code = $state<string | null>(null)
	let container = $state<HTMLDivElement | undefined>(undefined)

	$effect(() => {
		if (!open) {
			url = null
			code = null
		}
	})

	function show(shortUrl: string, shortCode: string) {
		url = shortUrl
		code = shortCode
		open = true
	}

	function downloadSVG() {
		const svg = container?.querySelector('svg')
		if (!svg) return
		const str = new XMLSerializer().serializeToString(svg)
		const blob = new Blob([str], { type: 'image/svg+xml' })
		const href = URL.createObjectURL(blob)
		const a = document.createElement('a')
		a.href = href
		a.download = `qr-${code}.svg`
		a.click()
		URL.revokeObjectURL(href)
	}

	function downloadPNG() {
		const svg = container?.querySelector('svg')
		if (!svg) return
		const str = new XMLSerializer().serializeToString(svg)
		const blob = new Blob([str], { type: 'image/svg+xml' })
		const svgUrl = URL.createObjectURL(blob)
		const img = new Image()
		img.onerror = () => URL.revokeObjectURL(svgUrl)
		img.onload = () => {
			const scale = 4
			const canvas = document.createElement('canvas')
			canvas.width = img.width * scale
			canvas.height = img.height * scale
			const ctx = canvas.getContext('2d')
			if (!ctx) { URL.revokeObjectURL(svgUrl); return }
			ctx.scale(scale, scale)
			ctx.drawImage(img, 0, 0)
			URL.revokeObjectURL(svgUrl)
			const a = document.createElement('a')
			a.href = canvas.toDataURL('image/png')
			a.download = `qr-${code}.png`
			a.click()
		}
		img.src = svgUrl
	}

	return {
		get open() { return open },
		set open(v: boolean) { open = v },
		get url() { return url },
		get container() { return container },
		set container(v: HTMLDivElement | undefined) { container = v },
		show,
		downloadSVG,
		downloadPNG,
	}
}

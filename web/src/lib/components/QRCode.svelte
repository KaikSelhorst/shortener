<script lang="ts">
	import qrcodegen from '$lib/qrcodegen'

	interface Props {
		/** The text/URI to encode — e.g. an otpauth:// URI. */
		data: string
		/** Pixel size of each module (default 4). */
		moduleSize?: number
		/** Quiet-zone border in modules (minimum 4 per spec). */
		border?: number
	}

	let { data, moduleSize = 4, border = 4 }: Props = $props()

	// Derive the SVG path string from the QR matrix.
	// Each dark module becomes a unit square appended to a single <path d>.
	let svgContent = $derived.by(() => {
		const qr = qrcodegen.QrCode.encodeText(data, qrcodegen.QrCode.Ecc.MEDIUM)
		const s = qr.size
		const total = (s + border * 2) * moduleSize

		let d = ''
		for (let y = 0; y < s; y++) {
			for (let x = 0; x < s; x++) {
				if (qr.getModule(x, y)) {
					const px = (x + border) * moduleSize
					const py = (y + border) * moduleSize
					d += `M${px},${py}h${moduleSize}v${moduleSize}h-${moduleSize}z`
				}
			}
		}

		return { d, size: total }
	})
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	width={svgContent.size}
	height={svgContent.size}
	viewBox="0 0 {svgContent.size} {svgContent.size}"
	role="img"
	aria-label="QR code"
>
	<!-- White background including quiet zone -->
	<rect width="100%" height="100%" fill="white" />
	<!-- All dark modules as a single path -->
	<path d={svgContent.d} fill="black" />
</svg>

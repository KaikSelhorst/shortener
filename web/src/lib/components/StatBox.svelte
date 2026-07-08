<script lang="ts">
	type Color = 'green' | 'cyan' | 'default'

	interface Props {
		label: string
		value: string | number
		color?: Color
	}

	let { label, value, color = 'default' }: Props = $props()

	const colorCls: Record<Color, string> = {
		green:   'text-success',
		cyan:    'text-info',
		default: 'text-foreground',
	}

	const display = $derived(typeof value === 'number' ? value.toLocaleString() : value)
</script>

<div class="px-5 py-4">
	<div class="tui-label text-muted-foreground">{label}</div>
	<div class="mt-1.5 text-2xl font-semibold tracking-tight {colorCls[color]}">{display}</div>
</div>

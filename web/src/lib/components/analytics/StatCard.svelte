<script lang="ts">
	interface Props {
		label: string
		value: string | number
		color?: 'success' | 'default'
		bordered?: boolean
	}

	let { label, value, color = 'default', bordered = true }: Props = $props()
</script>

<div class="px-6 py-4 {bordered ? 'border-r border-border' : ''}">
	<p class="text-[11px] font-semibold uppercase tracking-wider text-muted-foreground">{label}</p>
	<p class="mt-1 text-2xl font-semibold {color === 'success' ? 'text-success' : ''}">
		{typeof value === 'number' ? value.toLocaleString() : value}
	</p>
</div>

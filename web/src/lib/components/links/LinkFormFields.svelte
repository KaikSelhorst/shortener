<script lang="ts">
	import { Input } from '$lib/components/ui'

	interface UtmValues {
		source?: string
		medium?: string
		campaign?: string
		term?: string
		content?: string
	}

	interface Props {
		mode: 'create' | 'edit'
		utmExpanded: boolean
		url?: string
		title?: string
		expiresAt?: string
		maxClicks?: number | string
		utm?: UtmValues
	}

	let {
		mode,
		utmExpanded = $bindable(false),
		url = '',
		title = '',
		expiresAt = '',
		maxClicks = '',
		utm = {},
	}: Props = $props()
</script>

<Input name="url" type="url" label="URL" placeholder="https://example.com" required value={url} />
<Input name="title" type="text" label="Title" placeholder="optional title" value={title} />
{#if mode === 'create'}
	<Input
		name="custom_code"
		type="text"
		label="Short code"
		placeholder="my-link (leave blank to auto-generate)"
		pattern="[a-zA-Z0-9_\-]+"
		minlength={3}
		maxlength={50}
	/>
{/if}
<Input name="expires_at" type="datetime-local" label="Expires at" value={expiresAt} />
<Input name="max_clicks" type="number" label="Max clicks" placeholder="unlimited" min={1} value={maxClicks} />
<button
	type="button"
	onclick={() => (utmExpanded = !utmExpanded)}
	class="text-sm text-muted-foreground hover:text-foreground transition-colors text-left"
>
	{utmExpanded ? '▾' : '▸'} UTM parameters
</button>
{#if utmExpanded}
	<div class="flex flex-col gap-4 border-l border-border pl-3">
		<Input name="utm_source" type="text" label="Source" placeholder="newsletter, twitter, google" value={utm.source ?? ''} />
		<Input name="utm_medium" type="text" label="Medium" placeholder="email, social, cpc" value={utm.medium ?? ''} />
		<Input name="utm_campaign" type="text" label="Campaign" placeholder="spring_sale" value={utm.campaign ?? ''} />
		<Input name="utm_term" type="text" label="Term" placeholder="keywords (optional)" value={utm.term ?? ''} />
		<Input name="utm_content" type="text" label="Content" placeholder="variant_a (optional)" value={utm.content ?? ''} />
	</div>
{/if}

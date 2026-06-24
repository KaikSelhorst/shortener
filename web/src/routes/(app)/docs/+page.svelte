<script lang="ts">
	import type { PageData } from './$types'

	let { data }: { data: PageData } = $props()

	type SchemaField = { name: string; type: string; required: boolean; description: string; depth: number }

	type Endpoint = {
		id: string
		method: string
		path: string
		op: any
		tag: string
	}

	const HTTP_METHODS = ['get', 'post', 'put', 'delete', 'patch']

	const { spec, grouped, baseURL } = $derived.by(() => {
		const s = data.spec as any

		const eps: Endpoint[] = []
		for (const [path, pathItem] of Object.entries<any>(s.paths ?? {})) {
			for (const [method, op] of Object.entries<any>(pathItem)) {
				if (HTTP_METHODS.includes(method)) {
					eps.push({ id: `${method}::${path}`, method, path, op, tag: op.tags?.[0] ?? 'Other' })
				}
			}
		}

		const tagOrder: string[] = (s.tags ?? []).map((t: any) => t.name)
		const tagMeta: Record<string, string> = Object.fromEntries(
			(s.tags ?? []).map((t: any) => [t.name, t.description ?? ''])
		)

		return {
			spec: s,
			grouped: tagOrder.map((tag) => ({
				tag,
				description: tagMeta[tag] ?? '',
				endpoints: eps.filter((e) => e.tag === tag)
			})),
			baseURL: s.servers?.[0]?.url ?? ''
		}
	})

	let selectedId = $state<string | null>(null)

	const selected = $derived(
		selectedId
			? grouped.flatMap((g) => g.endpoints).find((e) => e.id === selectedId) ?? null
			: null
	)

	// ── schema resolution ────────────────────────────────────────────────────

	function resolveRef(schema: any): any {
		if (!schema?.$ref) return schema
		const name = (schema.$ref as string).split('/').pop()!
		return spec?.components?.schemas?.[name] ?? null
	}

	// Check $ref BEFORE resolving so the name is preserved
	function schemaTypeName(schema: any): string {
		if (!schema) return 'any'
		if (schema.$ref) return (schema.$ref as string).split('/').pop()!
		if (schema.type === 'array') return `${schemaTypeName(schema.items)}[]`
		if (schema.allOf) return schema.allOf.map(schemaTypeName).join(' & ')
		if (schema.format) return `${schema.type}(${schema.format})`
		return schema.type ?? 'object'
	}

	// Flat list for request body (top-level fields only)
	function flattenSchema(schema: any): SchemaField[] {
		if (!schema) return []
		const s = resolveRef(schema) ?? schema
		if (!s) return []
		if (s.allOf) return s.allOf.flatMap((p: any) => flattenSchema(resolveRef(p) ?? p))
		const props = s.properties ?? {}
		const required: string[] = s.required ?? []
		return Object.entries<any>(props).map(([name, def]) => ({
			name,
			type: schemaTypeName(def),
			required: required.includes(name),
			description: def.description ?? '',
			depth: 0
		}))
	}

	// Recursive expansion for responses — expands array items inline
	function expandSchema(schema: any, depth = 0, visited = new Set<string>()): SchemaField[] {
		if (!schema || depth > 2) return []

		// Top-level array: unwrap items and expand them
		if (schema.type === 'array' && schema.items) {
			const itemsResolved = resolveRef(schema.items) ?? schema.items
			if (itemsResolved?.properties || itemsResolved?.allOf) {
				return expandSchema(schema.items, depth, visited)
			}
			return []
		}

		const s = resolveRef(schema) ?? schema
		if (!s) return []

		if (schema.$ref) {
			const refName = (schema.$ref as string).split('/').pop()!
			if (visited.has(refName)) return []
			visited = new Set(visited)
			visited.add(refName)
		}

		if (s.allOf) {
			return s.allOf.flatMap((p: any) => expandSchema(resolveRef(p) ?? p, depth, visited))
		}

		const props = s.properties ?? {}
		const required: string[] = s.required ?? []
		const result: SchemaField[] = []

		for (const [name, def] of Object.entries<any>(props)) {
			result.push({
				name,
				type: schemaTypeName(def),
				required: required.includes(name),
				description: def.description ?? '',
				depth
			})

			// Expand array-of-objects items as indented rows
			if (def.type === 'array' && def.items) {
				const itemsResolved = resolveRef(def.items) ?? def.items
				if (itemsResolved?.properties || itemsResolved?.allOf) {
					result.push(...expandSchema(def.items, depth + 1, visited))
				}
			}
		}

		return result
	}

	function requestBodySchema(op: any): any {
		return op.requestBody?.content?.['application/json']?.schema ?? null
	}

	function responseSchema(res: any): any {
		return res.content?.['application/json']?.schema ?? null
	}

	// ── example generation + copy ────────────────────────────────────────────

	function schemaToExample(schema: any, visited = new Set<string>()): any {
		if (!schema) return null
		if (schema.$ref) {
			const name = (schema.$ref as string).split('/').pop()!
			if (visited.has(name)) return {}
			const next = new Set(visited)
			next.add(name)
			return schemaToExample(resolveRef(schema), next)
		}
		if (schema.allOf) {
			return schema.allOf.reduce((acc: any, part: any) => ({
				...acc,
				...schemaToExample(resolveRef(part) ?? part, visited)
			}), {})
		}
		if (schema.type === 'array') {
			const item = schemaToExample(schema.items, visited)
			return item !== null ? [item] : []
		}
		if (schema.type === 'object' || schema.properties) {
			const obj: Record<string, any> = {}
			for (const [key, def] of Object.entries<any>(schema.properties ?? {})) {
				obj[key] = schemaToExample(def as any, visited)
			}
			return obj
		}
		if (schema.enum) return schema.enum[0]
		if (schema.format === 'date-time') return '2024-01-01T00:00:00Z'
		if (schema.format === 'email') return 'user@example.com'
		if (schema.format === 'uri') return 'https://example.com'
		if (schema.format === 'uuid') return '00000000-0000-0000-0000-000000000000'
		if (schema.nullable) return null
		if (schema.type === 'integer') return 0
		if (schema.type === 'number') return 0
		if (schema.type === 'boolean') return false
		if (schema.type === 'string') return ''
		return null
	}

	let copiedKey = $state<string | null>(null)

	function copy(key: string, text: string) {
		navigator.clipboard.writeText(text).then(() => {
			copiedKey = key
			setTimeout(() => { copiedKey = null }, 1500)
		}).catch(() => {})
	}

	function copySchema(key: string, schema: any) {
		copy(key, JSON.stringify(schemaToExample(schema), null, 2))
	}

	// ── styling ──────────────────────────────────────────────────────────────

	const METHOD_COLOR: Record<string, string> = {
		get:    'text-tui-cyan',
		post:   'text-primary',
		put:    'text-tui-yellow',
		delete: 'text-tui-red',
		patch:  'text-tui-green'
	}

	const STATUS_COLOR: Record<string, string> = {
		'2': 'text-tui-green',
		'3': 'text-tui-cyan',
		'4': 'text-tui-yellow',
		'5': 'text-tui-red'
	}

	function statusColor(code: string) {
		return STATUS_COLOR[code[0]] ?? 'text-muted-foreground'
	}

	function hasAuth(op: any) {
		return op.security && op.security.length > 0
	}

	function authLabel(op: any): string {
		if (!hasAuth(op)) return 'none'
		const keys = op.security.flatMap((s: any) => Object.keys(s))
		return keys.includes('apiKeyAuth') ? 'jwt / api key' : 'jwt'
	}
</script>

<svelte:head>
	<title>API Reference — Shortener</title>
</svelte:head>

<div class="flex items-start gap-4">

		<!-- sidebar: sticky so it stays visible while scrolling route detail -->
		<nav class="tui-panel sticky top-4 w-56 shrink-0">
			<div class="tui-panel-header">▌ routes</div>
			{#each grouped as group}
				{#if group.endpoints.length > 0}
					<div class="border-b border-border">
						<div class="bg-card px-3 py-1.5">
							<span class="tui-label text-[9px]">{group.tag}</span>
						</div>
						{#each group.endpoints as ep}
							<button
								onclick={() => (selectedId = ep.id)}
								class="flex w-full items-center gap-2 border-l-2 px-3 py-1.5 text-left transition-colors hover:bg-secondary/60
									{selectedId === ep.id
										? 'border-primary bg-secondary/60'
										: 'border-transparent'}"
							>
								<span
									class="w-10 shrink-0 text-right font-mono text-[9px] uppercase
										{METHOD_COLOR[ep.method] ?? 'text-muted-foreground'}"
								>
									{ep.method}
								</span>
								<code
									class="min-w-0 flex-1 truncate font-mono text-[10px]
										{selectedId === ep.id ? 'text-foreground' : 'text-muted-foreground'}"
								>
									{ep.path}
								</code>
							</button>
						{/each}
					</div>
				{/if}
			{/each}
		</nav>

		<!-- detail -->
		<main class="min-w-0 flex-1">
			{#if selected}
				{@const reqSchema = requestBodySchema(selected.op)}
				{@const reqProps = flattenSchema(reqSchema)}
				<div class="space-y-4">

					<!-- headline -->
					<div class="tui-panel">
						<div class="tui-panel-header gap-2">
							<span class="{METHOD_COLOR[selected.method] ?? 'text-accent'}">
								{selected.method.toUpperCase()}
							</span>
							<code class="text-foreground">{selected.path}</code>
						</div>
						<div class="px-4 py-3">
							<p class="font-mono text-xs text-muted-foreground">{selected.op.summary ?? ''}</p>
							{#if selected.op.description}
								<p class="mt-1 font-mono text-xs text-muted-foreground/60">
									{selected.op.description}
								</p>
							{/if}
						</div>
					</div>

					<!-- security + parameters -->
					<div class="tui-panel">
						<div class="tui-panel-header">▌ security</div>
						<div class="px-4 py-3">
							<span
								class="font-mono text-xs
									{hasAuth(selected.op) ? 'text-tui-green' : 'text-muted-foreground'}"
							>
								{authLabel(selected.op)}
							</span>
						</div>

						{#if selected.op.parameters?.length}
							<div class="border-t border-border">
								<div class="tui-panel-header">▌ parameters</div>
								<div class="divide-y divide-border">
									{#each selected.op.parameters as param}
										{@const pKey = `param-${selected.id}-${param.name}`}
										<div class="flex items-start gap-4 px-4 py-2.5 font-mono text-xs">
											<span
												class="mt-0.5 w-9 shrink-0 text-center text-[9px] uppercase tracking-wide text-muted-foreground/40"
											>
												{param.in}
											</span>
											<div class="flex-1">
												<span
													class="{param.required
														? 'text-foreground'
														: 'text-muted-foreground'}"
												>
													{param.name}{param.required ? '*' : ''}
												</span>
												{#if param.description}
													<p class="mt-0.5 text-[11px] text-muted-foreground/60">
														{param.description}
													</p>
												{/if}
											</div>
											<span class="shrink-0 text-muted-foreground/50">
												{schemaTypeName(param.schema)}
											</span>
											<button
												onclick={() => copy(pKey, param.name)}
												class="shrink-0 cursor-pointer rounded-[2px] px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider transition-colors hover:bg-secondary/60
													{copiedKey === pKey ? 'text-tui-green' : 'text-muted-foreground/30 hover:text-muted-foreground'}"
											>
												{copiedKey === pKey ? 'copied' : 'copy'}
											</button>
										</div>
									{/each}
								</div>
							</div>
						{/if}
					</div>

					<!-- request body -->
					{#if reqSchema}
						{@const rbKey = `rb-${selected.id}`}
						<div class="tui-panel">
							<div class="tui-panel-header justify-between">
								<span>▌ request body</span>
								<button
									onclick={() => copySchema(rbKey, reqSchema)}
									class="cursor-pointer rounded-[2px] px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider transition-colors hover:bg-secondary/60
										{copiedKey === rbKey ? 'text-tui-green' : 'text-muted-foreground/40 hover:text-muted-foreground'}"
								>
									{copiedKey === rbKey ? 'copied' : 'copy json'}
								</button>
							</div>
							{#if reqProps.length > 0}
								<div class="divide-y divide-border">
									{#each reqProps as prop}
										<div class="flex items-start gap-4 px-4 py-2.5 font-mono text-xs">
											<div class="flex-1">
												<span
													class="{prop.required
														? 'text-foreground'
														: 'text-muted-foreground'}"
												>
													{prop.name}{prop.required ? '*' : ''}
												</span>
												{#if prop.description}
													<p class="mt-0.5 text-[11px] text-muted-foreground/60">
														{prop.description}
													</p>
												{/if}
											</div>
											<span class="shrink-0 text-muted-foreground/50">{prop.type}</span>
										</div>
									{/each}
								</div>
							{:else}
								<div class="px-4 py-3">
									<span class="font-mono text-xs text-muted-foreground/60">
										{schemaTypeName(reqSchema)}
									</span>
								</div>
							{/if}
						</div>
					{/if}

					<!-- responses -->
					<div class="tui-panel">
						<div class="tui-panel-header">▌ responses</div>
						<div class="divide-y divide-border">
							{#each Object.entries<any>(selected.op.responses ?? {}) as [code, res]}
								{@const rSchema = responseSchema(res)}
								{@const rFields = expandSchema(rSchema)}
								{@const rKey = `res-${selected.id}-${code}`}
								<div>
									<div class="tui-panel-header justify-between">
										<div class="flex items-baseline gap-2">
											<span class="font-bold {statusColor(code)}">{code}</span>
											<span class="text-muted-foreground">{res.description}</span>
										</div>
										{#if rSchema}
											<button
												onclick={() => copySchema(rKey, rSchema)}
												class="cursor-pointer rounded-[2px] px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider transition-colors hover:bg-secondary/60
													{copiedKey === rKey ? 'text-tui-green' : 'text-muted-foreground/40 hover:text-muted-foreground'}"
											>
												{copiedKey === rKey ? 'copied' : 'copy json'}
											</button>
										{/if}
									</div>
									{#if rFields.length > 0}
										<div class="divide-y divide-border">
											{#if rSchema?.type === 'array'}
												<div class="px-4 py-1.5 font-mono text-[10px] text-muted-foreground/40 italic">
													{schemaTypeName(rSchema)}
												</div>
											{/if}
											{#each rFields as field}
												<div
													class="flex items-center gap-4 px-4 py-1.5 font-mono text-[11px]"
													style={field.depth > 0 ? `padding-left: ${1 + field.depth * 0.875}rem` : ''}
												>
													{#if field.depth > 0}
														<span class="shrink-0 text-muted-foreground/25">└</span>
													{/if}
													<span class="flex-1 {field.depth > 0 ? 'text-muted-foreground/50' : 'text-muted-foreground/70'}">
														{field.name}
													</span>
													<span class="shrink-0 text-muted-foreground/40">{field.type}</span>
												</div>
											{/each}
										</div>
									{:else if rSchema}
										<div class="px-4 py-2">
											<p class="font-mono text-[11px] text-muted-foreground/50">
												{schemaTypeName(rSchema)}
											</p>
										</div>
									{/if}
								</div>
							{/each}
						</div>
					</div>

				</div>
			{:else}
				<div class="tui-panel">
					<div class="tui-panel-header">▌ route detail</div>
					<div class="px-4 py-14 text-center">
						<p class="font-mono text-xs text-muted-foreground">-- select a route --</p>
						<p class="mt-1 font-mono text-[10px] text-muted-foreground/40">
							{spec.info?.description}
						</p>
					</div>
				</div>
			{/if}
		</main>

</div>

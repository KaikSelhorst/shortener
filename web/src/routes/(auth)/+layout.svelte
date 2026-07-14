<script lang="ts">
  import type { Snippet } from "svelte";
  import AsciiPlasma from "$lib/components/ascii-plasma.svelte";

  let { children }: { children: Snippet } = $props();
</script>

<div class="flex min-h-screen">
  <div class="flex w-full flex-col items-center justify-center bg-card p-4 md:w-1/2">
    <div class="flex w-full max-w-sm flex-col gap-6">
      {@render children()}
    </div>
  </div>
  <div class="relative hidden items-center justify-center overflow-hidden bg-card md:flex md:w-1/2">
    <AsciiPlasma />
  </div>
</div>

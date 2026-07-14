<script lang="ts">
  import { page } from "$app/state";
  import { buttonVariants, buttonSizes } from "$lib/components/ui/button";
  import { renderAsciiCode } from "$lib/ascii-digits";

  const isNotFound = $derived(page.status === 404);
  const ascii = $derived(renderAsciiCode(String(page.status)));
</script>

<main class="fixed inset-0 flex flex-col items-center justify-center gap-3 p-8 text-center">
  <pre class="overflow-x-auto text-left text-[9px] leading-normal text-muted-foreground/40 sm:text-xs md:text-sm">{ascii}</pre>

  <h1 class="text-2xl font-semibold text-foreground">
    {isNotFound ? "Page not found" : "Something went wrong"}
  </h1>

  <p class="max-w-sm text-sm text-muted-foreground">
    {isNotFound
      ? "The page you're looking for doesn't exist or was moved."
      : (page.error?.message ?? "An unexpected error occurred.")}
  </p>

  <a
    href="/"
    class="mt-3 rounded-md px-3.5 py-2 text-sm font-medium transition-opacity {buttonVariants.primary} {buttonSizes.md}"
  >
    Back to home
  </a>
</main>

<script lang="ts">
  import { onMount } from "svelte";
  import { navigating } from "$app/state";

  let collapsed = $state(false);

  let fps = $state(0);
  let frameTime = $state(0);
  let heapUsed = $state<number | null>(null);
  let heapLimit = $state<number | null>(null);
  let lcp = $state<number | null>(null);
  let cls = $state(0);
  let inp = $state<number | null>(null);
  let ttfb = $state<number | null>(null);
  let domContentLoaded = $state<number | null>(null);
  let loadComplete = $state<number | null>(null);
  let resourceCount = $state(0);
  let transferSize = $state(0);
  let domNodes = $state(0);
  let longTaskCount = $state(0);
  let lastNavMs = $state<number | null>(null);

  let navStartedAt: number | null = null;
  $effect(() => {
    if (navigating.to) {
      navStartedAt = performance.now();
    } else if (navStartedAt !== null) {
      lastNavMs = performance.now() - navStartedAt;
      navStartedAt = null;
    }
  });

  function formatBytes(bytes: number) {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / 1024 / 1024).toFixed(2)} MB`;
  }

  function formatMs(ms: number | null) {
    if (ms === null) return "—";
    return `${ms.toFixed(0)}ms`;
  }

  function tone(value: number | null, good: number, poor: number) {
    if (value === null) return "text-muted-foreground";
    if (value <= good) return "text-success";
    if (value <= poor) return "text-warning";
    return "text-destructive";
  }

  const fpsTone = $derived(fps >= 50 ? "text-success" : fps >= 30 ? "text-warning" : "text-destructive");
  const lcpTone = $derived(tone(lcp, 2500, 4000));
  const clsTone = $derived(tone(cls, 0.1, 0.25));
  const inpTone = $derived(tone(inp, 200, 500));

  onMount(() => {
    let raf = 0;
    let frames = 0;
    let lastFpsSample = performance.now();
    let lastFrame = performance.now();

    function loop(now: number) {
      frameTime = now - lastFrame;
      lastFrame = now;
      frames++;
      if (now - lastFpsSample >= 500) {
        fps = Math.round((frames * 1000) / (now - lastFpsSample));
        frames = 0;
        lastFpsSample = now;
      }
      raf = requestAnimationFrame(loop);
    }
    raf = requestAnimationFrame(loop);

    const sampleInterval = setInterval(() => {
      const mem = (performance as unknown as { memory?: { usedJSHeapSize: number; jsHeapSizeLimit: number } })
        .memory;
      if (mem) {
        heapUsed = mem.usedJSHeapSize;
        heapLimit = mem.jsHeapSizeLimit;
      }
      domNodes = document.getElementsByTagName("*").length;
      const resources = performance.getEntriesByType("resource") as PerformanceResourceTiming[];
      resourceCount = resources.length;
      transferSize = resources.reduce((sum, r) => sum + (r.transferSize || 0), 0);
    }, 1000);

    const nav = performance.getEntriesByType("navigation")[0] as PerformanceNavigationTiming | undefined;
    if (nav) {
      ttfb = nav.responseStart - nav.requestStart;
      domContentLoaded = nav.domContentLoadedEventEnd - nav.startTime;
      loadComplete = nav.loadEventEnd - nav.startTime;
    }

    const observers: PerformanceObserver[] = [];

    try {
      const lcpObserver = new PerformanceObserver((list) => {
        const entries = list.getEntries() as PerformanceEntry[];
        const last = entries[entries.length - 1] as PerformanceEntry & { renderTime?: number; loadTime?: number };
        if (last) lcp = last.renderTime || last.loadTime || last.startTime;
      });
      lcpObserver.observe({ type: "largest-contentful-paint", buffered: true });
      observers.push(lcpObserver);
    } catch {
      // not supported in this browser
    }

    try {
      const clsObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries() as (PerformanceEntry & { hadRecentInput: boolean; value: number })[]) {
          if (!entry.hadRecentInput) cls += entry.value;
        }
      });
      clsObserver.observe({ type: "layout-shift", buffered: true });
      observers.push(clsObserver);
    } catch {
      // not supported in this browser
    }

    try {
      const inpObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries() as (PerformanceEntry & { interactionId?: number; duration: number })[]) {
          if (entry.interactionId && (inp === null || entry.duration > inp)) inp = entry.duration;
        }
      });
      inpObserver.observe({ type: "event", buffered: true, durationThreshold: 16 } as PerformanceObserverInit);
      observers.push(inpObserver);
    } catch {
      // not supported in this browser
    }

    try {
      const longTaskObserver = new PerformanceObserver((list) => {
        longTaskCount += list.getEntries().length;
      });
      longTaskObserver.observe({ type: "longtask", buffered: true });
      observers.push(longTaskObserver);
    } catch {
      // not supported in this browser
    }

    return () => {
      cancelAnimationFrame(raf);
      clearInterval(sampleInterval);
      for (const observer of observers) observer.disconnect();
    };
  });
</script>

<div class="fixed bottom-2 right-2 z-[9999] select-none font-mono text-[10px] leading-tight">
  {#if collapsed}
    <button
      type="button"
      onclick={() => (collapsed = false)}
      class="flex items-center gap-1.5 rounded-md border border-border bg-card/95 px-2 py-1 text-foreground shadow-lg backdrop-blur"
    >
      <span class="size-1.5 rounded-full bg-lime"></span>
      <span class={fpsTone}>{fps} fps</span>
    </button>
  {:else}
    <div class="w-56 rounded-md border border-border bg-card/95 p-2 text-foreground shadow-lg backdrop-blur">
      <div class="mb-1.5 flex items-center justify-between border-b border-border pb-1.5">
        <span class="flex items-center gap-1.5 font-semibold text-muted-foreground">
          <span class="size-1.5 rounded-full bg-lime"></span>
          PERF
        </span>
        <button type="button" onclick={() => (collapsed = true)} class="text-muted-foreground hover:text-foreground">
          ✕
        </button>
      </div>

      <div class="grid grid-cols-2 gap-x-2 gap-y-0.5">
        <span class="text-muted-foreground">FPS</span>
        <span class="text-right {fpsTone}">{fps} ({frameTime.toFixed(1)}ms)</span>

        <span class="text-muted-foreground">Heap</span>
        <span class="text-right"
          >{heapUsed !== null ? formatBytes(heapUsed) : "—"}{heapLimit !== null ? ` / ${formatBytes(heapLimit)}` : ""}</span
        >

        <span class="text-muted-foreground">DOM nodes</span>
        <span class="text-right">{domNodes}</span>

        <span class="text-muted-foreground">Long tasks</span>
        <span class="text-right {longTaskCount > 0 ? 'text-warning' : ''}">{longTaskCount}</span>

        <span class="col-span-2 mt-1 border-t border-border pt-1 text-muted-foreground">Web Vitals</span>

        <span class="text-muted-foreground">LCP</span>
        <span class="text-right {lcpTone}">{formatMs(lcp)}</span>

        <span class="text-muted-foreground">CLS</span>
        <span class="text-right {clsTone}">{cls.toFixed(3)}</span>

        <span class="text-muted-foreground">INP</span>
        <span class="text-right {inpTone}">{formatMs(inp)}</span>

        <span class="col-span-2 mt-1 border-t border-border pt-1 text-muted-foreground">Navigation</span>

        <span class="text-muted-foreground">TTFB</span>
        <span class="text-right">{formatMs(ttfb)}</span>

        <span class="text-muted-foreground">DOM ready</span>
        <span class="text-right">{formatMs(domContentLoaded)}</span>

        <span class="text-muted-foreground">Load</span>
        <span class="text-right">{formatMs(loadComplete)}</span>

        <span class="text-muted-foreground">Last route nav</span>
        <span class="text-right">{formatMs(lastNavMs)}</span>

        <span class="col-span-2 mt-1 border-t border-border pt-1 text-muted-foreground">Network</span>

        <span class="text-muted-foreground">Resources</span>
        <span class="text-right">{resourceCount}</span>

        <span class="text-muted-foreground">Transferred</span>
        <span class="text-right">{formatBytes(transferSize)}</span>
      </div>
    </div>
  {/if}
</div>

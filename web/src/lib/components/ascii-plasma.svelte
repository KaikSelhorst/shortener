<script lang="ts">
  const FONT_SIZE = 12; // px
  const LINE_HEIGHT = 1.15;
  const PROBE_TEXT = "MMMMMMMMMM";

  // Kept to the sparse end on purpose — dense glyphs (@, #, %) read as bold shapes
  // and pull focus away from the form, which is the opposite of what a background should do.
  const RAMP = "  ..::--";

  let probe: HTMLSpanElement;

  let width = $state(0);
  let height = $state(0);
  // Fallback estimate used only until the probe measures the real glyph metrics.
  let charWidth = $state(FONT_SIZE * 0.6);
  let charHeight = $state(FONT_SIZE * LINE_HEIGHT);

  // Overshoot by one cell in each axis so the grid always covers the container;
  // the wrapper clips the excess instead of leaving a gap from rounding/estimation error.
  const cols = $derived(Math.max(10, Math.ceil(width / charWidth) + 1));
  const rows = $derived(Math.max(10, Math.ceil(height / charHeight) + 1));

  let t = $state(0);

  $effect(() => {
    if (!probe) return;
    const rect = probe.getBoundingClientRect();
    charWidth = rect.width / PROBE_TEXT.length;
    charHeight = rect.height;
  });

  function render(time: number, c: number, r: number): string {
    if (c === 0 || r === 0) return "";

    const cx = c / 2 + (c / 4) * Math.sin(time * 0.3);
    const cy = r / 2 + (r / 4) * Math.cos(time * 0.4);

    const lines: string[] = [];
    for (let y = 0; y < r; y++) {
      let line = "";
      for (let x = 0; x < c; x++) {
        const v1 = Math.sin(x * 0.18 + time);
        const v2 = Math.sin(y * 0.14 - time * 0.8);
        const v3 = Math.sin((x + y) * 0.1 + time * 1.3);
        const dx = x - cx;
        const dy = (y - cy) * 2;
        const v4 = Math.sin(Math.sqrt(dx * dx + dy * dy) * 0.2 - time);

        const v = (v1 + v2 + v3 + v4) / 4; // -1..1
        const index = Math.floor(((v + 1) / 2) * (RAMP.length - 1));
        line += RAMP[index];
      }
      lines.push(line);
    }

    return lines.join("\n");
  }

  const text = $derived(render(t, cols, rows));

  $effect(() => {
    const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    if (reduceMotion) return;

    const id = setInterval(() => {
      t += 0.06;
    }, 90);
    return () => clearInterval(id);
  });
</script>

<div class="absolute inset-0 overflow-hidden" bind:clientWidth={width} bind:clientHeight={height}>
  <span
    bind:this={probe}
    aria-hidden="true"
    class="invisible absolute whitespace-pre font-mono leading-[1.15]"
    style="font-size: {FONT_SIZE}px;">{PROBE_TEXT}</span>
  <pre
    class="pointer-events-none m-0 select-none whitespace-pre font-mono leading-[1.15] text-foreground/[12%]"
    style="font-size: {FONT_SIZE}px;">{text}</pre>
</div>

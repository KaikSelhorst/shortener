<script lang="ts">
  import { renderProjectIcon } from "$lib/project-icon";

  interface Props {
    seed: string;
    size?: number;
  }

  let { seed, size = 28 }: Props = $props();
  let canvas: HTMLCanvasElement;

  $effect(() => {
    if (canvas) renderProjectIcon(canvas, seed, size);
  });
</script>

<canvas
  bind:this={canvas}
  width={size}
  height={size}
  style="width: {size}px; height: {size}px;"
  class="shrink-0 rounded-md"
  aria-hidden="true"
></canvas>

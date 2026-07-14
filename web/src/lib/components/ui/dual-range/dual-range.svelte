<script lang="ts">
  interface Props {
    min?: number;
    max?: number;
    minValue?: number;
    maxValue?: number;
    class?: string;
  }

  let { min = 0, max = 100, minValue = $bindable(0), maxValue = $bindable(100), class: className = "" }: Props = $props();

  const percent = (value: number) => ((value - min) / (max - min)) * 100;
</script>

<div class="relative h-4 w-full {className}">
  <div class="absolute inset-x-0 top-1/2 h-1 -translate-y-1/2 rounded-full bg-secondary"></div>
  <div
    class="absolute top-1/2 h-1 -translate-y-1/2 rounded-full bg-lime"
    style="left: {percent(minValue)}%; right: {100 - percent(maxValue)}%"
  ></div>
  <input
    type="range"
    {min}
    {max}
    value={minValue}
    oninput={(e) => (minValue = Math.min(Number(e.currentTarget.value), maxValue - 1))}
    class="range-dual"
  />
  <input
    type="range"
    {min}
    {max}
    value={maxValue}
    oninput={(e) => (maxValue = Math.max(Number(e.currentTarget.value), minValue + 1))}
    class="range-dual"
  />
</div>

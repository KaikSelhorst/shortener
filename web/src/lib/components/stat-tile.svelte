<script lang="ts">
  import * as Card from "$lib/components/ui/card";

  interface Props {
    label: string;
    value: number;
  }

  let { label, value }: Props = $props();
</script>

<Card.Root>
  <p class="text-sm text-muted-foreground">{label}</p>
  <p class="mt-1 text-3xl font-semibold text-foreground">{value.toLocaleString()}</p>
</Card.Root>

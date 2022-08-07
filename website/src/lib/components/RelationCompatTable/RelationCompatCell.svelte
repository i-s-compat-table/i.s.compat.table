<script lang="ts">
  import type { SupportRange } from "$lib/types";

  export let support: SupportRange | null = null;
</script>

<td
  class="support-cell {support?.isCurrent
    ? 'current'
    : support?.range
    ? 'deprecated'
    : ''}"
>
  {support?.range || ""}
</td>

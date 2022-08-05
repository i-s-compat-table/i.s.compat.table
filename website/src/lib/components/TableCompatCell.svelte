<script lang="ts">
  import type { RangeRef } from "$lib/types";
  export let data: RangeRef | null;
</script>

{#if data}
  <td class={data.isCurrent ? "currently-supported" : "deprecated"}>
    {#if data.url}
      <a href={data.url} target="_blank">{data.range}</a>
    {:else}
      <span>{data.range}</span>
    {/if}
  </td>
{:else}
  <td class="unsupported">X</td>
{/if}

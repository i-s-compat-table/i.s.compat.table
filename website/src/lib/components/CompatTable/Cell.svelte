<script lang="ts">
  import type { Munged } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  type _ = Munged[string][string][AllDbs];
  export let support: Munged[string][string][AllDbs] | undefined | null;
  const supported = support !== null && support !== undefined;
</script>

<td class="support-cell {supported ? 'supported' : 'unsupported'}">
  {#if supported}
    {JSON.stringify(support)}
    <!-- {#each Object.entries(support) as [kind, range]}
      <code>{kind}</code>: {range.range}
    {/each} -->
  {:else}
    nope
  {/if}
</td>

<style>
  :global(.support-cell) {
    text-align: center;
    background-color: lightsalmon;
  }
  :global(.support-cell.supported) {
    background-color: aquamarine;
  }
</style>

<script lang="ts">
  import type { Munged } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  type Kinds = Munged[string][string][AllDbs];
  export let kinds: Kinds | undefined | null;
  const kindRanges = Object.entries(kinds ?? {});
  const supported = kindRanges.some(([, ranges]) =>
    ranges.some((range) => range.isCurrent),
  );
</script>

<td class="support-cell {supported ? 'supported' : 'unsupported'}">
  {#if kindRanges.length === 0}
    nope
  {:else}
    {#each kindRanges as [kind, ranges]}
      {#if ranges.length > 1}
        Hey hey hey
      {/if}

      <code>{kind}</code>: {#each ranges as range}
        <div>
          {#if range.url}
            <a target="_blank" href={range.url}>{range.range}</a>
          {:else}
            {range.range}
          {/if}
        </div>
      {/each}
    {/each}
    <!-- {JSON.stringify(kinds)} -->
    <!-- {#each Object.entries(support) as [kind, range]}
      <code>{kind}</code>: {range.range}
    {/each} -->
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

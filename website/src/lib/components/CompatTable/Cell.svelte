<script lang="ts">
  import type { Munged } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  type Kinds = Munged[string][string][AllDbs];
  export let kinds: Kinds | undefined | null;
  const kindRanges = Object.entries(kinds ?? {});
  const anySupport = kindRanges.length > 0;
  const current = kindRanges.some(([, ranges]) =>
    ranges.some((range) => range.isCurrent),
  );
</script>

<td
  class="support-cell {current ? 'current' : anySupport ? 'deprecated' : 'unsupported'}"
>
  {#if !anySupport}
    unsupported
  {:else}
    {#each kindRanges as [kind, ranges]}
      {#if kind}
        <code>{kind}</code>:
      {/if}
      {#each ranges as range}
        <div>
          {#if range.url}
            <a target="_blank" href={range.url}>{range.range}</a>
          {:else}
            {range.range}
          {/if}
        </div>
      {/each}
    {/each}
  {/if}
</td>

<style>
  :global(.support-cell) {
    text-align: center;
    background-color: lightsalmon;
  }
  :global(.support-cell.deprecated) {
    background-color: aliceblue;
  }
  :global(.support-cell.current) {
    background-color: aquamarine;
  }
</style>

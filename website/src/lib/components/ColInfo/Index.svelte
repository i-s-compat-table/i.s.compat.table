<script lang="ts" context="module">
  const getMaxVersion = (r: UniqueVersionRange) =>
    Array.from(r.versions).sort(sortByVersion).reverse()[0];
</script>

<script lang="ts">
  import { sortByVersion, type Munged, type UniqueVersionRange } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  type Info = Munged[string][string][AllDbs];
  export let info: Info;

  const kinds = Object.entries(info ?? {})
    .reduce((a, [kind, ranges]) => {
      a.push([
        kind,
        ranges
          .sort((a, b) => sortByVersion(getMaxVersion(a), getMaxVersion(b)))
          .reverse(),
      ]);
      return a;
    }, [] as [string, UniqueVersionRange[]][])
    .sort(([, a], [, b]) => sortByVersion(getMaxVersion(a[0]), getMaxVersion(b[0])))
    .reverse();
</script>

{#each kinds as [kind, noteRanges]}
  <div class="type-history">
    <div>
      <span>type:</span> <code class="column-type">{kind || "unknown"}</code>
    </div>
    <div>
      {#each noteRanges as noteRange}
        {#if noteRange.note}
          <div>
            <div>Note:</div>
            <div>{noteRange.note}</div>
          </div>
          {#if noteRange.license}
            <div>
              {#if noteRange.license_url}
                <a href={noteRange.license_url} target="_blank"
                  >{noteRange.attribution || noteRange.license}</a
                >
              {:else}
                <span>{noteRange.attribution || noteRange.license}</span>
              {/if}
            </div>
          {/if}
        {/if}
        <div>
          {#if noteRange.url}
            <a href={noteRange.url} target="_blank">{noteRange.range}</a>
          {:else}
            <span>{noteRange.range}</span>
          {/if}
        </div>
      {/each}
    </div>
  </div>
{/each}

<style>
  .type-history {
    margin-top: 1ch;
    margin-bottom: 1ch;
  }
</style>

<script lang="ts" context="module">
  import { base } from "$app/paths";
  import { munge } from "$lib/munge";
  import type { Load } from "@sveltejs/kit";
  export const prerender = true;
  export const load: Load = async ({ fetch }) => {
    const target = `${base}/columns.tsv`;
    const [munged, tableSupport] = await fetch(target)
      .then((r) => {
        if (r.ok) return r.text();
        else {
          throw new Error(`error ${r.status}: ${r.statusText}`);
        }
      })
      .then((tsv) => munge(tsv)); // TODO: accept gzip/bzip
    return { props: { support: tableSupport } };
  };
</script>

<script lang="ts">
  import CommonalitySelector from "$lib/components/CommonalitySelector.svelte";
  import RelationCompatRow from "$lib/components/RelationCompatTable/RelationCompatRow.svelte";
  import type { TableSupportRow } from "$lib/components/RelationCompatTable/types";
  import type { TableSupport } from "$lib/munge";
  import { byKey } from "$lib/sort";
  import dbStore from "$lib/stores/dbs";
  export let support: TableSupport;

  let _rows: TableSupportRow[] = [];
  $: dbs = $dbStore;
  $: {
    _rows = Object.entries(support)
      .map(([table, dbSupport]) => ({
        name: table,
        specific: Object.keys(dbSupport).length === 1,
        universal: dbs.every((db) => db in dbSupport),
        support: dbSupport,
      }))
      .sort(byKey("name"));
  }
</script>

<h1 class="centered"><code>information_schema.*</code></h1>

<div class="centered">
  <CommonalitySelector />
</div>

<div class="centered">
  <table class="sticky-header">
    <thead>
      <tr>
        <th>relation</th>
        {#each dbs as db}
          <th>{db}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each _rows as row}
        <RelationCompatRow {row} />
      {/each}
    </tbody>
  </table>
</div>

<style>
  .centered {
    display: flex;
    justify-content: center;
  }
  .sticky-header > thead > tr > * {
    position: sticky;
    position: -webkit-sticky;
    position: -moz-sticky;
    position: -ms-sticky;
    position: -o-sticky;
    top: 0px;
    background-color: var(--bg-color, #fff);
    z-index: 1;
    border: none;
  }
  /* .sticky-header > tbody > tr:first-child > * {
    border: none;
  } */
  .sticky-header > thead > tr > *::after {
    content: "";
    position: absolute;
    left: 0;
    bottom: 0;
    width: 100%;
    border-bottom: 2px solid #ddd;
  }
  :global(.hidden) {
    display: none;
  }
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

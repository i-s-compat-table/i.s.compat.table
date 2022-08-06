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
  import type { TableSupport } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  export let support: TableSupport;
  type SupportRange = { range: string; isCurrent: boolean };
  type Desired = {
    name: string;
    specific: boolean;
    universal: boolean;
    support: Partial<Record<AllDbs, SupportRange>>;
  };
  let _rows: Desired[] = Object.entries(support)
    .map(([table, dbSupport]) => ({
      name: table,
      specific: Object.keys(dbSupport).length === 1,
      universal: dbs.every((db) => db in dbSupport),
      support: dbSupport,
    }))
    .sort((a, b) => (a.name > b.name ? 1 : b.name > a.name ? -1 : 0));
  const dbs: AllDbs[] = ["mysql", "mariadb", "tidb", "postgres", "cockroachdb", "mssql"];
  let showSpecific = true;
</script>

<!-- TODO: factor the specificity checkbox into its own store & component. -->
<label
  ><input type="checkbox" bind:checked={showSpecific} />show relations specific to one
  database</label
>

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
      <tr class={row.specific ? (showSpecific ? "" : "hidden") : ""}>
        <th id="relation-{row.name}">{row.name}</th>
        <!-- TODO: display self-link on-hover? -->
        {#each dbs as db}
          <!-- TODO: factor table support cells into their own component-->
          <td
            class="support-cell {row.support[db]?.isCurrent
              ? 'current'
              : row.support[db]?.range
              ? 'deprecated'
              : ''}"
          >
            {row.support[db]?.range || ""}
          </td>
        {/each}
      </tr>
    {/each}
  </tbody>
</table>

<style>
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

<script lang="ts">
  import type { Munged } from "$lib/munge";
  import { objEntries } from "$lib/sort";
  import dbs from "$lib/stores/dbs";
  import InfoTable from "./Table.svelte";
  export let columnSupport: Munged = {};
  $: rows = objEntries(columnSupport);
</script>

<table class="sticky-header">
  <thead>
    <tr>
      <th>relation</th>
      <th>column</th>
      {#each $dbs as db}
        <th>{db}</th>
      {/each}
      <th />
    </tr>
  </thead>
  <tbody style="will-change: contents;">
    {#each rows as [tableName, tableData]}
      <InfoTable name={tableName} columns={tableData} />
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
</style>

<script lang="ts">
  import type { Munged, TableSupport } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  import InfoTable from "./Table.svelte";
  export let dbs: AllDbs[] = [];
  export let columnSupport: Munged = {};
  export let tableSupport: TableSupport = {};
</script>

<table class="sticky-header">
  <thead>
    <tr>
      <th>relation</th>
      <th>column</th>
      {#each dbs as db}
        <th>{db}</th>
      {/each}
      <th />
    </tr>
  </thead>
  <tbody>
    {#each Object.entries(columnSupport) as [tableName, tableData]}
      <InfoTable
        name={tableName}
        support={tableSupport[tableName] ?? {}}
        columns={tableData}
        {dbs}
      />
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

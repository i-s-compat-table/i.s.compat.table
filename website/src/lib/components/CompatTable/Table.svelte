<script lang="ts">
  import type { Munged, TableSupport } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  import Cell from "./Cell.svelte";
  import Column from "./Column.svelte";
  type TableInfo = Munged[string];
  export let dbs: AllDbs[];
  export let name: string;
  export let support: TableSupport[string] = {};
  export let columns: TableInfo;
</script>

<tr>
  <th class="table-support-header" rowspan={Object.keys(columns).length + 1}
    ><code>{name}</code></th
  >
  <td>
    <!-- "columns" column should be empty for a table-row -->
  </td>
  {#each dbs as db}
    <Cell />
    <!-- TODO: any compatibility? -->
  {/each}
</tr>
{#each Object.entries(columns) as [colName, colInfo]}
  <Column name={colName} support={colInfo} {dbs} />
{/each}

<style>
  :global(.table-support-header) {
    vertical-align: top;
  }
</style>

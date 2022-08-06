<script lang="ts">
  import type { Munged } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  import Column from "./Column.svelte";
  type TableInfo = Munged[string];
  export let dbs: AllDbs[];
  export let name: string;
  // export let support: TableSupport[string] = {};
  export let columns: TableInfo;
</script>

<tr>
  <th
    id="table-{name}"
    class="table-support-header"
    rowspan={Object.keys(columns).length + 1}><code>{name}</code></th
  >
  <td colspan={dbs.length + 1}>
    <hr />
    <!-- "columns" column should be empty for a table-row -->
  </td>
  <!-- {#each dbs as db}
    <Cell />
    TODO: any compatibility?
  {/each} -->
</tr>
{#each Object.entries(columns) as [colName, colInfo]}
  <Column name={colName} support={colInfo} {dbs} />
{/each}

<style>
  :global(.table-support-header) {
    vertical-align: top;
  }
</style>

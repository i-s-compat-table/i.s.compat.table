<script lang="ts">
  import type { Munged } from "$lib/munge";
  import dbs from "$lib/stores/dbs";
  import Column from "./Column.svelte";
  export let name: string;
  type TableInfo = Munged[string];
  export let columns: TableInfo;
</script>

<tr>
  <th
    id="table-{name}"
    class="table-support-header"
    rowspan={Object.keys(columns).length + 1}><code>{name}</code></th
  >
  <td colspan={$dbs.length + 1} class="table-support-header" />
</tr>
{#each Object.entries(columns) as [colName, colInfo]}
  <Column name={colName} support={colInfo} />
{/each}

<style>
  :global(.table-support-header) {
    vertical-align: top;
    position: sticky;
    top: 1em;
    background-color: var(--bg-color, #fff);
    border-top: 1px solid black;
  }
</style>

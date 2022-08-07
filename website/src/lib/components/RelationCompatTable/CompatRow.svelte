<script lang="ts">
  import type { AllDbs, CompatibilityRow } from "$lib/types";
  import TableCompatCell from "./TableCompatCell.svelte";
  export let data: CompatibilityRow;
  export let dbNames: AllDbs[] = ["postgres", "mysql", "mariadb", "mssql"];
  const { tableName, columns, ...dbs } = data;
</script>

<tr>
  <th colspan={columns.length}>{tableName}</th>
  {#each dbNames as dbName}
    <TableCompatCell data={dbs[dbName] ?? null} />
  {/each}
</tr>

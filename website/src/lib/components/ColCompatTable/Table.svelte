<script lang="ts">
  import type { Munged } from "$lib/munge";
  import { objEntries } from "$lib/sort";
  import commonality from "$lib/stores/commonality";
  import dbs from "$lib/stores/dbs";
  import Column from "./Column.svelte";
  export let name: string;
  type TableInfo = Munged[string];
  export let columns: TableInfo;
  $: rows = objEntries(columns).filter(([colName, col]) => {
    switch ($commonality) {
      case "any":
        return true;
      case "db-specific":
        return $dbs.filter((db) => col[db]).length === 1;
      case "shared":
        if (name === "administrable_role_authorizations" && colName === "grantee") {
          console.log($dbs.filter((db) => col[db]));
          console.log(colName, col["mysql"]);
          console.log($dbs.filter((db) => col[db]).length);
        }
        return $dbs.filter((db) => col[db]).length > 1;
      case "universal":
        if (name === "columns" && colName === "character_set_name") {
          console.log(
            colName,
            $dbs.filter((db) => col[db]),
            col["mariadb"],
          );
        }
        return $dbs.filter((db) => col[db]).length === $dbs.length;
    }
  });
  $: hidden = rows.length === 0;
</script>

<tr class:hidden>
  <th id="table-{name}" class="table-support-header" rowspan={rows.length + 1}>
    <code>{name}</code>
  </th>
  <td colspan={$dbs.length + 1} class="table-support-header" />
</tr>
{#each rows as [name, support] (name)}
  <Column {name} {support} />
{/each}

<style>
  :global(.table-support-header) {
    vertical-align: top;
    position: sticky;
    top: 1em;
    background-color: var(--bg-color, #fff);
    border-top: 1px solid black;
  }
  :global(.hidden) {
    display: none;
  }
</style>

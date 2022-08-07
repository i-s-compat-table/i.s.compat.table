<script lang="ts">
  import commonality from "$lib/stores/commonality";
  import dbs from "$lib/stores/dbs";
  import RelationCompatCell from "./RelationCompatCell.svelte";
  import type { TableSupportRow } from "./types";
  export let row: TableSupportRow;
  let hidden: boolean = false;
  $: {
    switch ($commonality) {
      case "any":
        hidden = false;
        break;
      case "db-specific":
        hidden = !row.specific;
        break;
      case "shared":
        hidden = row.specific;
        break;
      case "universal":
        hidden = !row.universal;
        break;
    }
  }
</script>

<tr class:hidden>
  <th id="relation-{row.name}"><code>{row.name}</code></th>
  <!-- TODO: display self-link on-hover? -->
  {#each $dbs as db}
    <RelationCompatCell support={row.support[db] ?? null} />
  {/each}
</tr>

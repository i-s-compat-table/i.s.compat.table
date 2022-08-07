<script lang="ts">
  import { base } from "$app/paths";
  import type { Munged } from "$lib/munge";
  import { getColumnPath } from "$lib/nav";
  import dbs from "$lib/stores/dbs";
  import Cell from "./Cell.svelte";
  export let table: string;
  export let name: string;
  export let support: Munged[string][string];
  export let hidden: boolean = false;
  // $: cells = $dbs.map((db) => support[db] ?? null);
</script>

<tr class:hidden>
  <th><a href={getColumnPath(base, table, name)}><code>{name}</code></a></th>
  {#each $dbs as db (db)}
    <Cell kinds={support[db] ?? null} />
  {/each}
</tr>

<script lang="ts">
  import type { Munged } from "$lib/munge";
  import dbs from "$lib/stores/dbs";
  import Cell from "./Cell.svelte";
  export let name: string;
  export let support: Munged[string][string];
  // $: cells = $dbs.map((db) => support[db] ?? null);
</script>

<tr>
  <th><code>{name}</code></th>
  {#each $dbs as db (db)}
    <Cell kinds={support[db] ?? null} />
  {/each}
</tr>

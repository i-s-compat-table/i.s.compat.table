<script lang="ts">
  import type { Munged } from "$lib/munge";
  import type { AllDbs } from "$lib/types";
  import Cell from "./Cell.svelte";
  export let name: string;
  export let dbs: AllDbs[] = [];
  export let support: Munged[string][string];
</script>

<tr>
  <th><code>{name}</code></th>
  {#each dbs as db}
    <Cell support={support[db] ?? null} />
  {/each}
</tr>

<script lang="ts">
  import type { Munged } from "$lib/munge";
  import { objEntries } from "$lib/sort";
  import commonality from "$lib/stores/commonality";
  import dbs from "$lib/stores/dbs";
  import Column from "./Column.svelte";
  export let name: string;
  type TableInfo = Munged[string];
  export let columns: TableInfo;
  $: rows = objEntries(columns);
  $: ubiquity = rows.map(([colName, col]) => $dbs.filter((db) => col[db]).length);
  const isSpecific = (u: number) => u === 1;
  const isShared = (u: number) => u > 1;
  const isAny = (_: any) => true;
  $: isUniversal = (u: number) => u === $dbs.length;
  $: shouldShow = (() => {
    switch ($commonality) {
      case "any":
        return isAny;
      case "db-specific":
        return isSpecific;
      case "shared":
        return isShared;
      case "universal":
        return isUniversal;
    }
  })();
  $: show = ubiquity.map(shouldShow);
  $: nShown = show.filter(Boolean).length;
  $: hidden = nShown === 0;
</script>

<tr class:hidden>
  <th id="table-{name}" class="table-support-header" rowspan={nShown + 1}>
    <code>{name}</code>
  </th>
  <td colspan={$dbs.length + 1} class="table-support-header" />
</tr>
{#each rows as [name, support], i (name)}
  <Column {name} {support} hidden={!show[i]} />
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

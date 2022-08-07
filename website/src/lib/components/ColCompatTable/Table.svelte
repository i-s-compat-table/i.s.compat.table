<script lang="ts">
  import { base } from "$app/paths";

  import type { Munged } from "$lib/munge";
  import { getRelationPath } from "$lib/nav";
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
    <a style="witdh: 100%" href={getRelationPath(base, name)}><code>{name}</code></a>
  </th>
  <td class="table-support-header" colspan={$dbs.length + 1}>
    <span>&nbsp;</span><!-- <hr /> -->
  </td>
</tr>
{#each rows as [colName, support], i (colName)}
  <Column table={name} name={colName} {support} hidden={!show[i]} />
{/each}

<style>
  :global(.table-support-header) {
    vertical-align: top;
    position: sticky;
    top: 1em;
    background-color: var(--bg-color, #fff);
  }
  :global(td.table-support-header) {
    border-bottom: 1px solid black;
  }
  :global(.table-support-header > a) {
    text-decoration: none;
    display: inline-block;
    width: 100%;
    border-bottom: 1px solid;
  }
  :global(.hidden) {
    display: none;
  }
</style>

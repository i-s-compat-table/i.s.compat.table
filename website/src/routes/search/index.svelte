<script lang="ts" context="module">
  import { base } from "$app/paths";
  import type { RawData } from "$lib/munge";
  import type { Load } from "@sveltejs/kit";
  import { tsvParse } from "d3-dsv";

  export const prerender = true;
  export const load: Load = async ({ fetch }) => {
    const target = `${base}/columns.tsv`;
    const data = await fetch(target)
      .then((r) => {
        if (r.ok) return r.text();
        else {
          throw new Error(`error ${r.status}: ${r.statusText}`);
        }
      })
      .then((tsv) =>
        Array.from(
          (tsvParse(tsv) as RawData[]).reduce((a, r) => {
            const ref = r.table_name + "." + r.column_name;
            a.add(ref);
            return a;
          }, new Set<string>()),
        )
          .sort()
          .map((ref) => {
            const [table, column] = ref.split(".");
            return { table, column };
          }),
      )
      .then((r) => r);
    return { props: { data } };
  };
</script>

<script lang="ts">
  import SearchResult from "$lib/components/Search/SearchResult.svelte";
  import query from "$lib/stores/search";
  import Fuse from "fuse.js";
  export let data: { table: string; column: string }[];
  // TODO: allow selecting keys
  const index = new Fuse(data, {
    keys: ["table", "column"],
    includeMatches: true,
    minMatchCharLength: 3,
    threshold: 0.1,
  });
  $: results = index.search($query);
  // TODO: debounce search?
</script>

<div>
  <input
    type="text"
    bind:value={$query}
    placeholder="search by relation and/or column name"
    style="width: 50%"
  />
  <ol>
    {#each results as result (result.item.table + "." + result.item.column)}
      <SearchResult {result} />
    {/each}
  </ol>
</div>

<style>
  ol {
    list-style: none;
  }
</style>

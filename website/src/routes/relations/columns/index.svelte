<script lang="ts" context="module">
  import { base } from "$app/paths";
  import { munge, type Munged } from "$lib/munge";
  import type { Load } from "@sveltejs/kit";
  export const prerender = true;
  export const load: Load = async ({ fetch }) => {
    const target = `${base}/columns.tsv`;
    const [munged, tableSupport] = await fetch(target)
      .then((r) => {
        if (r.ok) return r.text();
        else {
          throw new Error(`error ${r.status}: ${r.statusText}`);
        }
      })
      .then((tsv) => munge(tsv)); // TODO: accept gzip/bzip?
    return { props: { columnSupport: munged, tableSupport } };
  };
</script>

<script lang="ts">
  import ColCompatTable from "$lib/components/ColCompatTable/Index.svelte";
  import CommonalitySelector from "$lib/components/CommonalitySelector.svelte";
  import type { TableSupport } from "$lib/munge";

  export let columnSupport: Munged;
  export let tableSupport: TableSupport;
</script>

<h1><code>information_schema</code> compatibility table</h1>
<CommonalitySelector />
<!-- TODO: docs -->
<ColCompatTable {columnSupport} />

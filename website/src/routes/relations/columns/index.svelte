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
    return { props: { columnSupport: munged } };
  };
</script>

<script lang="ts">
  import ColCompatTable from "$lib/components/ColCompatTable/Index.svelte";
  import CommonalitySelector from "$lib/components/CommonalitySelector.svelte";
  export let columnSupport: Munged;
</script>

<h1 class="centered"><code>information_schema.*.*</code></h1>
<div class="centered">
  <CommonalitySelector />
</div>
<!-- TODO: docs -->
<div class="centered">
  <ColCompatTable {columnSupport} />
</div>

<style>
  .centered {
    display: flex;
    justify-content: center;
  }
</style>

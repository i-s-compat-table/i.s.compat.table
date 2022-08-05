<script lang="ts" context="module">
  import { base } from "$app/paths";
  import { munge, type Munged } from "$lib/munge";
  import type { Load } from "@sveltejs/kit";
  export const prerender = true;
  export const load: Load = async ({ fetch }) => {
    const target = `${base}/columns.tsv`;
    const munged = await fetch(target)
      .then((r) => {
        if (r.ok) return r.text();
        else {
          throw new Error(`error ${r.status}: ${r.statusText}`);
        }
      })
      .then((tsv) => munge(tsv)); // TODO: accept gzip/bzip
    return { props: { data: munged } };
  };
</script>

<script lang="ts">
  import CompatTable from "$lib/components/CompatTable/Index.svelte";
  import type { AllDbs } from "$lib/types";
  const allDbs: AllDbs[] = [
    "mysql",
    "mariadb",
    "tidb",
    "postgres",
    "cockroachdb",
    "mssql",
  ];
  export let data: Munged;
</script>

<h1><code>information_schema</code> compatibility table</h1>
<!-- TODO: docs -->
<CompatTable dbs="{allDbs}columnSupport{data}" />

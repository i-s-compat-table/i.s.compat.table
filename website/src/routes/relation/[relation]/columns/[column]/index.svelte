<script lang="ts" context="module">
  import { base } from "$app/paths";
  import { munge, type Munged } from "$lib/munge";
  import type { Load } from "@sveltejs/kit";
  export const prerender = true;
  export const load: Load = async ({ fetch, params }) => {
    const target = `${base}/columns.tsv`;
    const [munged] = await fetch(target)
      .then((r) => {
        if (r.ok) return r.text();
        else {
          throw new Error(`error ${r.status}: ${r.statusText}`);
        }
      })
      .then((tsv) =>
        munge(
          tsv,
          (r) => r.table_name === params.relation && r.column_name === params.column,
        ),
      );
    return {
      props: {
        tableName: params.relation,
        columnName: params.column,
        columnSupport: munged[params.relation][params.column],
      },
    };
  };
</script>

<script lang="ts">
  import Db from "$lib/components/ColInfo/Db.svelte";
  import dbs from "$lib/stores/dbs";
  export let tableName: string;
  export let columnName: string;
  export let columnSupport: Munged[string][string];
</script>

<h1><code>information_schema.{tableName}.{columnName}</code></h1>

{#each $dbs as db}
  <section id="db-${db}">
    <h2><a href="#db-{db}">{db}</a></h2>
    <Db info={columnSupport[db]} />
  </section>
{/each}

<details open>
  <pre><code>{JSON.stringify(columnSupport, null, 2)}</code></pre>
</details>

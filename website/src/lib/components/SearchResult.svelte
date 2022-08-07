<script lang="ts">
  /// reference path="fuse.js"
  import { base } from "$app/paths";
  import { flatten } from "$lib/array";

  import { getColumnPath, getRelationPath } from "$lib/nav";
  import SearchHighlight from "./SearchHighlight.svelte";
  export let result: Fuse.FuseResult<{ table: string; column: string }>;
</script>

<li>
  <span class="relation">
    <a href={getRelationPath(base, result.item.table)}>
      <SearchHighlight
        result={result.item.table}
        matches={result.matches
          .filter((m) => m.key === "table")
          .map((m) => m.indices)
          .reduce(flatten, [])}
      />
    </a>
  </span>
  <span class="column">
    <a href={getColumnPath(base, result.item.table, result.item.column)}>
      <SearchHighlight
        result={result.item.column}
        matches={result.matches
          .filter((m) => m.key === "column")
          .map((m) => m.indices)
          .reduce(flatten, [])}
      />
    </a>
  </span>
</li>

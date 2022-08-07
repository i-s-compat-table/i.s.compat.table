<script lang="ts">
  import { base } from "$app/paths";
  import { page } from "$app/stores";
  import { issues } from "$lib/nav";
</script>

<ul>
  {#if $page.routeId}
    <li>
      <a href={base + "/"}>Home</a>
    </li>
  {/if}
  <li>
    <a href="{base}/relations">View all relations in <code>information_schema</code></a>
  </li>
  <li>
    <a href="{base}/columns">View all columns in <code>information_schema</code></a>
  </li>
  <li><a href="{base}/search">Search for a relation or column</a></li>
  <li><a href={issues} target="_blank">Report an issue or inaccuracy</a></li>
</ul>

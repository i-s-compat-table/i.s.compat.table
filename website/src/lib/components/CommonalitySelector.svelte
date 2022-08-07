<script lang="ts">
  import commonality from "$lib/stores/commonality";
</script>

<label>
  <input type="radio" bind:group={$commonality} value="any" />
  show all items
</label>

<label>
  <input type="radio" bind:group={$commonality} value="shared" />
  show items implemented by at least two databases
</label>

<label>
  <input type="radio" bind:group={$commonality} value="db-specific" />
  show items implemented by only one database
</label>

<label>
  <input type="radio" bind:group={$commonality} value="universal" />
  show items implemented by all databases
</label>

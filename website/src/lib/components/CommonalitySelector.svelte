<script lang="ts">
  import commonality from "$lib/stores/commonality";
</script>

<div style="display: flex">
  <!-- TODO: layout? -->
  <div style="margin: 1ch">show items implemented by</div>
  <label>
    <input type="radio" bind:group={$commonality} value="any" />
    any
  </label>

  <label>
    <input type="radio" bind:group={$commonality} value="shared" />
    at least two
  </label>

  <label>
    <input type="radio" bind:group={$commonality} value="db-specific" />
    only one database
  </label>

  <label>
    <input type="radio" bind:group={$commonality} value="universal" />
    all databases
  </label>
</div>

<style>
  label {
    margin: 1ch;
  }
</style>
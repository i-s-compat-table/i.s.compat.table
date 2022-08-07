<script lang="ts">
  import NavLinks from "$lib/components/Navigation/NavLinks.svelte";
</script>

<main>
  <h1 class="centered"><code>information_schema</code> compatibility tables</h1>
  <p>
    This site provides compatibility tables comparing different SQL databases'
    implementaions of the <code>information_schema</code> views.
  </p>
  <NavLinks />
</main>

<style>
  main {
    max-width: 70ch;
    padding: 2ch;
    margin: auto;
  }
</style>

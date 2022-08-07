<script lang="ts">
  import { browser } from "$app/env";
  import { afterNavigate } from "$app/navigation";
  import { Hamburger } from "svelte-hamburgers";
  import NavLinks from "./NavLinks.svelte";

  let open = false;

  const toggleMenu = () => (open = false);
  if (browser) {
    document.addEventListener("keyup", (e) => {
      if (e.key === "Escape") toggleMenu();
    });
    afterNavigate(toggleMenu);
  }
</script>

<div class:open>
  <Hamburger bind:open --color="var(--font-color, black)" />
  <nav style="display: {open ? 'flex' : 'none'}">
    <NavLinks />
  </nav>
</div>

<style>
  div,
  nav {
    position: fixed;
    background-color: var(--bg-color, #fff);
    justify-content: center;
    align-items: center;
    z-index: 999;
    right: 1ch;
    top: 1ch;
  }
  div.open {
    width: 100%;
    height: 100%;
  }
  nav {
    width: 100%;
    flex-direction: column;
    max-width: 50vw;
    height: 70vh;
    max-height: 1040px;
    left: 25vw;
    right: 25vw;
  }
  :global(button.hamburger.hamburger--spin) {
    background-color: var(--bg-color, #fff) !important;
    position: fixed;
    top: 1ch;
    right: 1ch;
  }
</style>

<script lang="ts">
  export let result: string;
  export let matches: [number, number][];
  // let chunks: string[] = [];
  // const chunkMeaning: boolean[] = []; // true = highlight, false = normal
  // let index = 0;
  $: chunks = matches.reduce(
    (() => {
      let index = 0;
      return (acc: [string, boolean][], [start, end], i, arr) => {
        if (start > index) {
          acc.push([result.slice(index, start), false]);
        }
        acc.push([result.slice(start, end + 1), true]);
        index = end + 1;
        if (i === arr.length - 1 && index < result.length) {
          acc.push([result.slice(index), false]);
        }
        return acc;
      };
    })(),
    [] as [string, boolean][],
  );
</script>

<code>
  {#if chunks.length}
    {#each chunks as [chunk, highlighted]}
      {#if highlighted}
        <b>{chunk}</b>
      {:else}
        {chunk}
      {/if}
    {/each}
  {:else}
    {result}
  {/if}
</code>

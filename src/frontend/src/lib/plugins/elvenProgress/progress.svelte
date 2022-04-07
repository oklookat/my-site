<script lang="ts">
  import { onMount } from "svelte";
  import ElvenProgress from ".";

  let line: HTMLDivElement;

  onMount(() => {
    new ElvenProgress(line);
  });
</script>

<div class="progress">
  <div class="progress__line" bind:this={line} />
</div>

<style lang="scss">
  .progress {
    cursor: default;
    position: absolute;
    width: 100%;
    height: 2px;
    &__line {
      height: 100%;
      width: 0;
      background-color: #a097dc;
    }
  }
</style>

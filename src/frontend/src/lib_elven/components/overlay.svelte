<script lang="ts">
  import { onDestroy, onMount } from "svelte";

  export let onClose: () => void;

  onMount(() => {
    document.body.classList.add("no-scroll");
  });

  onDestroy(() => {
    document.body.classList.remove("no-scroll");
  });
</script>

<div class="overlay base__overlay" on:click|self={onClose}>
  <div class="overlay__main">
    <div class="overlay__content">
      <slot />
    </div>
  </div>
</div>

<style lang="scss">
  .overlay {
    &__main {
      border-radius: var(--border-radius);
      z-index: 9999;
      overflow: auto;
      max-width: 365px;
      max-height: 365px;
      width: 250px;
      height: 250px;
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      @media (prefers-color-scheme: light) {
        color: #000;
        background-color: #fff;
      }
      @media (prefers-color-scheme: dark) {
        color: #fff;
        background-color: #202020;
      }
    }
    &__content {
      box-sizing: border-box;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 24px;
      font-size: 1.2rem;
      height: 100%;
      width: 100%;
    }
  }
</style>

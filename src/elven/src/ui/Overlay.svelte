<script lang="ts">
  import { createEventDispatcher, onDestroy } from "svelte";

  export let active: boolean;

  const dispatch = createEventDispatcher<{ deactivated: boolean }>();

  onDestroy(() => {
    document.body.classList.remove("no-scroll");
  });

  $: watchActive(active);
  function watchActive(value: boolean) {
    if (value) {
      document.body.classList.add("no-scroll");
      return
    }
    document.body.classList.remove("no-scroll");
  }

  function deactivate() {
    document.body.classList.remove("no-scroll");
    dispatch("deactivated");
  }
</script>

{#if active}
  <div class="overlay" on:click|self={deactivate}>
    <div class="overlay__main">
      <div class="overlay__content">
        <slot />
      </div>
    </div>
  </div>
{/if}

<style lang="scss">
  .overlay {
    background-color: rgba(0, 0, 0, 0.4);
    z-index: 9998;
    max-width: 100vw;
    width: 100%;
    height: 100%;
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
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

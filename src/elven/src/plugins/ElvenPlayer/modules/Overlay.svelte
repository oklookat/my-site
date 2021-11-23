<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from "svelte";

  const dispatch = createEventDispatcher();

  const setScroll = (v: boolean) => {
    if (!v) {
      dispatch("deactivated");
    }
    document.body.style.overflow = v ? "hidden" : null;
  };

  let container: HTMLDivElement;

  onMount(() => {
    document.body.appendChild(container);
    setScroll(true);
  });

  onDestroy(() => {
    setScroll(false);
  });

  function deactivate(e: MouseEvent | TouchEvent) {
    const notLMB = e instanceof MouseEvent && e.button !== 0;
    if (notLMB) {
      // prevent decativate when not LMB clicked
      return;
    }
    setScroll(false);
  }
</script>

<div
  class="overlay"
  bind:this={container}
  on:mousedown|self|stopPropagation={deactivate}
>
  <div class="overlay__main">
    <div class="overlay__content">
      <slot />
    </div>
  </div>
</div>

<style lang="scss">
  .overlay {
    background-color: rgba(0, 0, 0, 0.7);
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
    align-items: center;
    justify-content: center;
    &__main {
      border-radius: 12px;
      z-index: 9999;
      max-width: 365px;
      max-height: 524px;
      width: 294px;
      height: 224px;
      box-sizing: border-box;
      font-size: 1rem;
    }
    &__content {
      background-color: white;
      overflow: auto;
      border-radius: inherit;
      width: 100%;
      height: 100%;
      @media (prefers-color-scheme: light) {
        background-color: #ececec;
        color: black;
      }
      @media (prefers-color-scheme: dark) {
        background-color: #202020;
        color: white;
      }
    }
  }
</style>

<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from "svelte";
  import SliderLogic from "./logic";

  const dispatch = createEventDispatcher();

  // dispatch event only after mouse up
  export let afterUp: boolean = false;
  // initial percents
  export let percents = 0;

  // used only for slider style
  let elementPercents = "0%";

  // slider container - used by logic class. Need for set events etc.
  let container;

  // main logic
  let sl = new SliderLogic();

  // init
  let ready = false;

  // set style on slider element
  function calcStyles(percents: number) {
    elementPercents = `${percents}%`;
  }

  let isMouseDown = false;
  let slidePercents = 0;
  // when user mouse down on slider
  const s1 = sl.slider.isMouseDown.subscribe((v) => {
    isMouseDown = v;
    dispatcher();
  });
  // when user change percents by dragging a slider
  const s2 = sl.slider.percents.subscribe((v) => {
    slidePercents = v;
    dispatcher();
  });

  // dispatch events
  function dispatcher() {
    const isMouseUp = !isMouseDown;
    if (afterUp) {
      // dispatch event only when mouse up after dragging
      if (isMouseUp) {
        dispatch("slide", slidePercents);
      }
    } else {
      // dispatch event after dragging
      dispatch("slide", slidePercents);
    }
    // set percents on element anyway
    calcStyles(slidePercents);
  }

  // on slider percents changed by prop
  $: watchPercents(percents);
  function watchPercents(v) {
    if (isMouseDown) {
      // block change percents because now user drag slider
      return;
    }
    calcStyles(percents);
  }

  // init
  onMount(() => {
    if (percents > 0) {
      calcStyles(percents);
    }
    if (!container) {
      return;
    }
    sl.init(container);
    ready = true;
  });

  // cleanup
  onDestroy(() => {
    s1();
    s2();
    sl.destroy();
  });
</script>

<div class="slider" bind:this={container}>
  {#if ready}
    <div class="slider__line" style="width: {elementPercents}" />
    <div class="slider__bubble" style="left: calc({elementPercents} - 3px)" />
  {/if}
</div>

<style lang="scss">
  .slider {
    width: calc(100% - 3px);
    height: 100%;
    box-sizing: border-box;
    position: relative;
    display: flex;
    align-items: center;
    cursor: pointer;
    &__line,
    &__bubble {
      position: absolute;
      background-color: #a097dc;
    }
    &__line {
      border-top-left-radius: 2px;
      border-bottom-left-radius: 2px;
      top: auto;
      bottom: 0;
      width: 0;
      height: 100%;
    }
    &__bubble {
      transition: width, height 100ms;
      border-radius: 50%;
      left: 0;
      margin-left: -6px;
      width: 20px;
      height: 20px;
    }
  }
</style>

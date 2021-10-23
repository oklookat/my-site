<script lang="ts">
  import SliderLogic from "@/common/plugins/ElvenPlayer/components/SliderLogic";
  import { createEventDispatcher, onDestroy, onMount } from "svelte";

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

<div class="slider__container" bind:this={container}>
  {#if ready}
    <div class="slider__sliders">
      <div class="slider__sl" style="width: {elementPercents}" />
      <div class="slider__bubble" style="left: calc({elementPercents} - 3px)" />
    </div>
  {/if}
</div>

<style scoped>
  .slider__container {
    width: 100%;
    height: 100%;
  }

  .slider__sliders {
    width: calc(100% - 3px);
    height: 100%;
    position: relative;
    display: flex;
    align-items: center;
  }

  .slider__container,
  .slider__sliders,
  .slider__sl,
  .slider__bubble {
    cursor: pointer;
  }

  .slider__sl,
  .slider__bubble {
    position: absolute;
    background-color: #9d6a89;
  }

  .slider__sl {
    border-top-left-radius: 2px;
    border-bottom-left-radius: 2px;
    width: 0;
    height: 100%;
    top: auto;
    bottom: 0;
  }

  .slider__bubble {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    box-sizing: border-box;
    left: 0;
    margin-left: -6px;
    transition: width, height 100ms;
  }
</style>

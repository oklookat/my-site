<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from "svelte";
  import Core from "./core";

  const dispatch = createEventDispatcher<{
    /** on slider percents changed */
    slide: number;
    /** on slider mouse up (true) or down (false) */
    mouse: boolean;
  }>();

  /** dispatch percents event only after mouse up */
  export let afterUp: boolean = false;

  /** set percents */
  export let percents: number = 0;

  let container: HTMLDivElement;
  let core: Core;
  let mouseDown = false;

  /** percents when user dragging slider */
  let slidePercents = percents;
  /** this will be displayed */
  let finalPercents = percents;
  /** set style on slider element */
  const setFinalPercents = (perc: number) => {
    finalPercents = perc;
  };

  /** when we subscribing to store, callback immediately called. We don't need this. */
  let ready = false;
  let unsub1;
  let unsub2;

  function subscribe() {
    // when user mouse down on slider
    unsub1 = core.state.store.isMouseDown.subscribe((v) => {
      if (!ready) {
        return;
      }
      mouseDown = v;
      dispatcher("mouse");
    });
    // when user change percents by dragging a slider
    unsub2 = core.state.store.percents.subscribe((v) => {
      if (!ready) {
        return;
      }
      slidePercents = v;
      dispatcher("perc");
    });
  }

  function unsubscribe() {
    unsub1();
    unsub2();
  }

  onMount(() => {
    core = new Core(container);
    subscribe()
    ready = true;
  });

  onDestroy(() => {
    unsubscribe()
    core.destroy();
  });

  /** watch when percents prop changed */
  $: watchPercents(percents);
  function watchPercents(perc: number) {
    if (mouseDown) {
      // if user dragging slider, prevent change style by prop
      return;
    }
    setFinalPercents(perc);
  }

  /** dispatch slide events */
  function dispatcher(by: "mouse" | "perc") {
    if (by === "mouse") {
      dispatch("mouse", mouseDown);
    }
    if (afterUp && !mouseDown) {
      // dispatch event only when mouse up after dragging
      dispatch("slide", slidePercents);
    } else if (!afterUp && mouseDown) {
      // dispatch event after dragging
      dispatch("slide", slidePercents);
    }
    setFinalPercents(slidePercents);
  }
</script>

<div class="slider" bind:this={container}>
  <div class="slider__line" style="width: {finalPercents}%" />
  <div class="slider__bubble" style="left: calc({finalPercents}% - 3px)" />
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
      background-color: #918CE6;
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

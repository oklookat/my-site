<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from "svelte";
  import Slider from "./Slider.svelte";
  import Progress from "./Progress.svelte";
  import type PlayerLogic from "../logic/PlayerLogic";
  import PlaybackControls from "./PlaybackControls.svelte";

  export let pl: PlayerLogic;

  const dispatch = createEventDispatcher();

  let noScroll = false;
  $: watchNoScroll(noScroll);
  function watchNoScroll(value) {
    switchScroll();
  }

  let audioOverlayContainer;
  onMount(() => {
    document.body.appendChild(audioOverlayContainer);
    noScroll = true;
  });

  onDestroy(() => {
    document.body.style.overflow = null;
  });

  function switchScroll() {
    if (noScroll) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = null;
    }
  }

  function deactivate(event) {
    if (event.which !== 1 && event.type !== "touchstart") {
      return;
    }
    noScroll = false;
    dispatch("deactivated");
  }

  function onProgressSliderTriggered(percents: number) {
    pl.setTimeByPercents(percents);
  }

  function onVolumeSliderTriggered(percents: number) {
    pl.setVolumeByPercents(percents);
  }

  // state
  let buffered = 0;
  const s1 = pl.player.currentPlaying.percentsBuffered.subscribe((v) => {
    buffered = v;
  });
  let reached = 0;
  const s2 = pl.player.currentPlaying.percentsReached.subscribe((v) => {
    reached = v;
  });
  let currentTime = "00:00";
  const s3 = pl.player.currentPlaying.currentTime.subscribe((v) => {
    currentTime = v;
  });
  let duration = "00:00";
  const s4 = pl.player.currentPlaying.duration.subscribe((v) => {
    duration = v;
  });
  let percentsVolume = 100;
  const s5 = pl.player.percentsVolume.subscribe((v) => {
    percentsVolume = v;
  });
  onDestroy(() => {
    s1();
    s2();
    s3();
    s4();
    s5();
  });
</script>

<div
  class="overlay__container"
  bind:this={audioOverlayContainer}
  on:mousedown|self|stopPropagation={deactivate}
>
  <div class="content">
    <div class="playback__time">
      <div class="sliders__time">
        <div class="progress__buffered">
          <Progress bind:percents={buffered} />
        </div>
        <div class="slider__time">
          <Slider
            afterUp={true}
            on:slide={(e) => onProgressSliderTriggered(e.detail)}
            bind:percents={reached}
          />
        </div>
      </div>
      <div class="durations">
        <div class="duration__current">
          {currentTime}
        </div>
        <div class="duration__total">
          {duration}
        </div>
      </div>
    </div>

    <div class="playback__volume">
      <div class="sliders__volume">
        <Slider
          percents={percentsVolume}
          afterUp={false}
          on:slide={(e) => onVolumeSliderTriggered(e.detail)}
        />
      </div>
    </div>

    <div class="playback__controls">
      <PlaybackControls {pl} />
    </div>
  </div>
</div>

<style>
  .overlay__container {
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
  }

  .content {
    background-color: white;
    border-radius: 12px;
    z-index: 9999;
    overflow: auto;
    max-width: 365px;
    max-height: 524px;
    width: 294px;
    height: 224px;
    box-sizing: border-box;
    font-size: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 24px;
  }

  .playback__time,
  .playback__volume {
    height: max-content;
    display: flex;
    flex-direction: column;
  }

  .playback__time {
    margin-bottom: 24px;
    width: 80%;
    gap: 12px;
  }

  .sliders__time,
  .sliders__volume {
    position: relative;
    border-radius: 4px;
    background-color: rgba(0, 0, 0, 0.5);
    height: 14px;
    width: 100%;
  }

  .slider__time,
  .progress__buffered {
    border-radius: inherit;
    position: absolute;
    width: 100%;
    height: 100%;
  }

  .durations {
    display: flex;
    flex-direction: row;
  }

  .duration__total {
    margin-left: auto;
  }

  .playback__volume {
    height: 24px;
    width: 50%;
  }

  .playback__controls {
    display: flex;
    width: 128px;
    height: 30px;
    fill: black;
    justify-content: center;
  }

  @media (prefers-color-scheme: light) {
    .content {
      background-color: #ececec;
      color: black;
    }
  }

  @media (prefers-color-scheme: dark) {
    .content {
      background-color: #202020;
      color: white;
    }
    .playback__controls {
      fill: white;
    }
  }
</style>

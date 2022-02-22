<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import type { ElvenPlayer, ComponentState } from "./types";
  import PlaybackControls from "./parts/playbackControls.svelte";
  import OverlayMenu from "./parts/overlayMenu.svelte";
  import type { Unsubscriber } from "svelte/store";

  /** is player active */
  let active: boolean = false;
  export let core: ElvenPlayer;

  /** when player open/close */
  export let onActiveChanged: (active?: boolean) => void | undefined =
    undefined;
  $: onActiveChange(active);
  function onActiveChange(val: boolean) {
    if (!onActiveChanged) {
      return;
    }
    onActiveChanged(val);
  }

  /** is player controls overlay active */
  let isOverlay = false;

  /** player state */
  let state: ComponentState = {
    playing: false,
    volume: {
      percents: 100,
    },
    current: {
      buffered: {
        percents: 0,
      },
      time: {
        draggingNow: false,
        percents: 0,
        pretty: "00:00",
      },
      duration: {
        pretty: "00:00",
      },
    },
  };

  let unsub1: Unsubscriber,
    unsub2: Unsubscriber,
    unsub3: Unsubscriber,
    unsub4: Unsubscriber,
    unsub5: Unsubscriber,
    unsub6: Unsubscriber;
  let unsubs: Unsubscriber[] = [];

  onMount(() => {
    init();
  });

  onDestroy(() => {
    destroy();
  });

  function init() {
    unsub1 = core.store.state.playing.onChange((v) => {
      if (v) {
        active = true;
      }
      state.playing = v;
    });

    unsub2 = core.store.state.current.buffered.percents.onChange((v) => {
      state.current.buffered.percents = v;
    });

    unsub3 = core.store.state.current.time.percents.onChange((v) => {
      state.current.time.percents = v;
    });

    unsub4 = core.store.state.current.time.pretty.onChange((v) => {
      // not setting pretty if user dragging time slider now (time preview)
      if (state.current.time.draggingNow) {
        return;
      }
      state.current.time.pretty = v;
    });

    unsub5 = core.store.state.current.duration.pretty.onChange((v) => {
      state.current.duration.pretty = v;
    });

    unsub6 = core.store.state.volume.percents.onChange((v) => {
      state.volume.percents = v;
    });

    unsubs.push(unsub1, unsub2, unsub3, unsub4, unsub5, unsub6);
  }

  function destroy() {
    for (const unsub of unsubs) {
      unsub();
    }
    unsubs = [];
    active = false;
    core.stop();
  }

  /** events */

  function onPlay() {
    active = true;
    core.play();
  }

  function onPause() {
    core.pause();
  }

  function onNext() {
    core.next();
  }

  function onPrev() {
    core.prev();
  }

  function onClose() {
    active = false;
    core.stop();
  }

  function setVolumePercents(perc: number) {
    core.volumePercents = perc;
  }

  function setCurrentTimePercents(perc: number) {
    core.currentTimePercents = perc;
  }

  function setCurrentTimePreview(perc: number) {
    state.current.time.pretty = core.convertPercentsToCurrentTimePretty(perc);
  }
</script>

{#if active}
  <div class="player">
    <div class="player__show">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 191 168"
        on:click={() => (isOverlay = !isOverlay)}
      >
        <rect x="1" y="138" width="190" height="30" rx="15" />
        <rect y="69" width="190" height="30" rx="15" />
        <rect x="1" width="190" height="30" rx="15" />
      </svg>
    </div>

    <div class="player__controls" on:click|stopPropagation>
      <PlaybackControls
        bind:isPlaying={state.playing}
        on:play={() => onPlay()}
        on:pause={() => onPause()}
        on:next={() => onNext()}
        on:prev={() => onPrev()}
      />
    </div>

    <div class="player__close">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 255.07 295.91"
        on:click={() => onClose()}
      >
        <path
          d="M135,390.18h0a33.69,33.69,0,0,0,48-4.29L369.93,159.34a35.1,35.1,0,0,0-4.19-48.86h0a33.69,33.69,0,0,0-48,4.29L130.83,341.33A35.09,35.09,0,0,0,135,390.18Z"
          transform="translate(-122.85 -102.37)"
        />
        <path
          d="M365.74,390.18h0a33.68,33.68,0,0,1-48-4.29L130.83,159.34A35.1,35.1,0,0,1,135,110.48h0a33.7,33.7,0,0,1,48,4.29L369.93,341.33A35.09,35.09,0,0,1,365.74,390.18Z"
          transform="translate(-122.85 -102.37)"
        />
      </svg>
    </div>
  </div>

  {#if isOverlay}
    <OverlayMenu
      bind:state
      on:deactivated={() => (isOverlay = false)}
      on:volumeChanged={(e) => setVolumePercents(e.detail)}
      on:currentTimeChanged={(e) => setCurrentTimePercents(e.detail)}
      on:currentTimePreviewChanged={(e) => setCurrentTimePreview(e.detail)}
    >
      <PlaybackControls
        slot="playbackControls"
        bind:isPlaying={state.playing}
        on:play={() => onPlay()}
        on:pause={() => onPause()}
        on:next={() => onNext()}
        on:prev={() => onPrev()}
      />
    </OverlayMenu>
  {/if}
{/if}

<style lang="scss">
  .player {
    border-top: var(--color-border) 1px solid;
    background-color: var(--color-level-1);
    width: 100%;
    height: 64px;
    display: grid;
    grid-template-columns: 52px 1fr 52px;
    svg {
      height: 20px;
      width: 100%;
    }
    > div {
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
      > svg {
        cursor: pointer;
      }
    }
    &__controls,
    svg {
      @media (prefers-color-scheme: dark) {
        fill: white;
      }
      @media (prefers-color-scheme: light) {
        fill: black;
      }
    }
  }
</style>

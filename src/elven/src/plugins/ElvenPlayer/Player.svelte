<script lang="ts">
  import { onDestroy } from "svelte";
  import type { IElvenPlayer, TPlaylist, TSource } from "./types";
  import Core from "./core";
  import PlaybackControls from "./components/parts/PlaybackControls.svelte";
  import OverlayMenu from "./components/parts/OverlayMenu.svelte";

  export let isActive = false;

  let isOverlay = false;

  let core = new Core();

  onDestroy(() => {
    core.destroy();
  });

  class Plugin implements IElvenPlayer {
    public async addToPlaylist(url: TSource) {
      core.addToPlaylist(url);
    }

    public setPlaylist(playlist: TPlaylist) {
      core.playlist = playlist;
    }

    public async play(url: string) {
      this.setPlaylist({ position: 0, sources: [url] });
      await core.play();
      isActive = true;
    }
  }
  window.$elvenPlayer = new Plugin();
</script>

{#if isActive}
  <div class="player">
    <div class="player__show">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 191 168"
        on:click={() => (isOverlay = !isOverlay)}
      >
        <rect x="1" y="138" width="190" height="30" rx="15" style="fill:#fff" />
        <rect y="69" width="190" height="30" rx="15" style="fill:#fff" />
        <rect x="1" width="190" height="30" rx="15" style="fill:#fff" />
      </svg>
    </div>

    <div class="player__controls" on:click|stopPropagation>
      <PlaybackControls {core} />
    </div>

    <div class="player__close">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 255.07 295.91"
        on:click={() => (isActive = false)}
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

  <OverlayMenu bind:core bind:active={isOverlay} />
{/if}

<style lang="scss">
  .player {
    background-color: rgba(0, 0, 0, 0.5);
    border-top-right-radius: 6px;
    border-top-left-radius: 6px;
    height: 64px;
    width: 214px;
    position: fixed;
    left: 50%;
    bottom: 0;
    transform: translate(-50%, 0);
    display: grid;
    grid-template-rows: 1fr;
    grid-template-columns: 1fr 2fr 1fr;
    justify-items: center;
    svg {
      fill: white;
      height: 20px;
      width: 20px;
    }
    > div {
      height: 100%;
      width: 100%;
      display: flex;
      justify-items: center;
      justify-content: center;
      > svg {
        cursor: pointer;
        align-self: center;
      }
    }
    &__controls {
      fill: white;
    }
  }
</style>

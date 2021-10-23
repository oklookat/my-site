<script lang="ts">
  import type PlayerLogic from "@/common/plugins/ElvenPlayer/components/PlayerLogic";
  import PlaybackControls from "./PlaybackControls.svelte";
  import Overlay from "@/common/plugins/ElvenPlayer/components/Overlay.svelte";

  // TODO: fix component crash after set audio time to end in overlay slider
  // plugin controls
  export let pl: PlayerLogic;
  export let isActive = false;
  $: watchActive(isActive);
  function watchActive(value) {
    if (value) {
      pl.init();
    } else {
      pl.destroy();
    }
  }
  let isAudioOverlayActive = false;
</script>

{#if isActive}
  <div class="player__container">
    <div class="player__playback-1">
      <div class="player__show">
        <svg
          class="player-item player__show-butt"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 191 168"
          on:click={() => (isAudioOverlayActive = !isAudioOverlayActive)}
        >
          <rect
            x="1"
            y="138"
            width="190"
            height="30"
            rx="15"
            style="fill:#fff"
          />
          <rect y="69" width="190" height="30" rx="15" style="fill:#fff" />
          <rect x="1" width="190" height="30" rx="15" style="fill:#fff" />
        </svg>
      </div>
    </div>
    <div class="player__playback-2" on:click|stopPropagation>
      <div class="playback__controls"><PlaybackControls {pl} /></div>
    </div>
    <div class="player__playback-3">
      <div class="player-item player__close">
        <svg
          class="player__close-butt"
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

    {#if isAudioOverlayActive}
      <Overlay on:deactivated={() => (isAudioOverlayActive = false)} bind:pl />
    {/if}
  </div>
{/if}

<style scoped>
  .player__container {
    border-top-right-radius: 6px;
    border-top-left-radius: 6px;
    height: 64px;
    width: 214px;
    background-color: rgba(0, 0, 0, 0.5);
    position: fixed;
    left: 50%;
    bottom: 0;
    transform: translate(-50%, 0);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .player__container svg {
    max-height: 20px;
    max-width: 20px;
    height: 20px;
    fill: white;
  }

  .player-item {
    cursor: pointer;
    align-self: center;
  }

  .player__playback-1,
  .player__playback-2 {
    width: 36px;
  }

  .player__playback-1,
  .player__playback-2,
  .player__playback-3 {
    height: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
  }

  .player__playback-2 {
    width: 100%;
  }

  .playback__controls {
    fill: white;
  }

  .player__show,
  .player__close {
    display: flex;
    flex-direction: row;
    align-items: center;
  }
  .player__show {
    margin-left: 16px;
  }

  .player__close {
    margin-right: 16px;
  }
</style>

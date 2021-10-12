<template>
  <div class="audio-player-container"
       v-if="player.audioPlayer.active">

    <div class="playback-tools">
      <div class="show-player">
        <svg class="audio-player-item show-player-butt"
             v-on:click="isAudioOverlayActive = !isAudioOverlayActive"
             xmlns="http://www.w3.org/2000/svg"
             viewBox="0 0 191 168">
          <rect x="1" y="138" width="190" height="30" rx="15" style="fill:#fff"/>
          <rect y="69" width="190" height="30" rx="15" style="fill:#fff"/>
          <rect x="1" width="190" height="30" rx="15" style="fill:#fff"/>
        </svg>
      </div>
    </div>
    <div class="playback-tools-2" v-on:click.stop>
      <PlaybackControls class="playback-controls" :player="player"></PlaybackControls>
    </div>
    <div class="playback-tools-3">
      <div class="close-player">
        <svg v-on:click="isPlayerActive = false"
            data-name="close-butt" xmlns="http://www.w3.org/2000/svg"
             viewBox="0 0 255.07 295.91">
          <path d="M135,390.18h0a33.69,33.69,0,0,0,48-4.29L369.93,159.34a35.1,35.1,0,0,0-4.19-48.86h0a33.69,33.69,0,0,0-48,4.29L130.83,341.33A35.09,35.09,0,0,0,135,390.18Z" transform="translate(-122.85 -102.37)"/>
          <path d="M365.74,390.18h0a33.68,33.68,0,0,1-48-4.29L130.83,159.34A35.1,35.1,0,0,1,135,110.48h0a33.7,33.7,0,0,1,48,4.29L369.93,341.33A35.09,35.09,0,0,1,365.74,390.18Z" transform="translate(-122.85 -102.37)"/>
        </svg>
      </div>
    </div>

    <Overlay
        v-if="isAudioOverlayActive"
        v-on:deactivated="isAudioOverlayActive = false"
        v-bind:player="player">
    </Overlay>

  </div>
</template>


<script setup lang="ts">
import {reactive, ref, toRef, watch} from "vue";
import Composition from "@/common/plugins/ElvenPlayer/components/Player/Composition"
import PlaybackControls from "./PlaybackControls.vue";
import Overlay from "./../Overlay/Component.vue"
const SERVICE = ref('ELVEN_PLAYER_C')
const player = reactive(new Composition())
const isAudioOverlayActive = ref(false)

const isPlayerActive = toRef(player.audioPlayer, "active")
watch(isPlayerActive, (active) => {
  if (active) {
    player.init()
  } else {
    player.destroy()
  }
})
</script>

<style scoped>
.audio-player-container {
  border-top-right-radius: 6px;
  border-top-left-radius: 6px;
  height: 64px;
  width: 214px;
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  left: 50%;
  bottom: 0;
  transform: translate(-50%, 0);
  margin: 0 auto;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.audio-player-container svg {
  max-height: 20px;
  max-width: 20px;
  height: 20px;
  fill: white;
}

.audio-player-item {
  cursor: pointer;
  align-self: center;
}

.playback-tools,
.playback-tools-3{
  width: 36px;
}

.playback-tools,
.playback-tools-2,
.playback-tools-3{
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.playback-tools-2 {
  width: 100%;
}


.playback-controls{
  fill: white;
}

.show-player,
.close-player{
  display: flex;
  flex-direction: row;
  align-items: center;
}
.show-player{
  margin-left: 16px;
}
.close-player{
  margin-right: 16px;
}

.tools svg {
  max-width: 100%;
  max-height: 100%;
  width: 80%;
}

.show-player-butt {
  cursor: pointer;
}
</style>
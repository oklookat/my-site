<template>
  <div style="width: 200px; height: 200px">{{ player.audioPlayer }}</div>
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
      <div class="playback-controls">
        <svg v-on:click="player.prev()" class="audio-player-item prev-butt" xmlns="http://www.w3.org/2000/svg"
             viewBox="0 0 406.76 169.37">
          <path
              d="M180.76,320.07V179.49A14.38,14.38,0,0,0,159.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,180.76,320.07Z"
              transform="translate(-28 -165.09)"/>
          <path
              d="M434.76,320.07V179.49A14.38,14.38,0,0,0,413.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,434.76,320.07Z"
              transform="translate(-28 -165.09)"/>
        </svg>
        <div class="audio-player-item play-pause">
          <svg v-if="!player.audioPlayer.isPlaying"
               v-on:click="player.play()"
               class="play-butt"
               xmlns="http://www.w3.org/2000/svg"
               viewBox="0 0 152.76 169.37">
            <path
                d="M173.62,320.29V179.71a14.38,14.38,0,0,1,21.47-12.51l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.38,14.38,0,0,1,173.62,320.29Z"
                transform="translate(-173.62 -165.31)"/>
          </svg>
          <svg v-if="player.audioPlayer.isPlaying"
               v-on:click="player.pause()"
               class="pause-butt"
               xmlns="http://www.w3.org/2000/svg"
               viewBox="0 0 207 310">
            <rect width="70" height="310"/>
            <rect x="137" width="70" height="310"/>
          </svg>
        </div>
        <svg v-on:click="player.next()" class="audio-player-item audio-player-next" xmlns="http://www.w3.org/2000/svg"
             viewBox="0 0 406.76 169.37">
          <path
              d="M308,320.07V179.49A14.37,14.37,0,0,1,329.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,308,320.07Z"
              transform="translate(-54 -165.09)"/>
          <path
              d="M54,320.07V179.49A14.37,14.37,0,0,1,75.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,54,320.07Z"
              transform="translate(-54 -165.09)"/>
        </svg>
      </div>
    </div>

    <Overlay
        v-if="isAudioOverlayActive"
        v-on:deactivated="isAudioOverlayActive = false"
        v-on:slide-progress="onSlideProgress"
        v-bind:percents-time="player.audioPlayer.currentPlaying.percentsReached"
        v-bind:percents-buffered="player.audioPlayer.currentPlaying.percentsBuffered">
    </Overlay>

  </div>
</template>


<script setup lang="ts">
import {reactive, ref, toRef, watch} from "vue";
import Overlay from './Overlay.vue'
import ElvenPlayerComposition from "@/common/plugins/ElvenPlayer/core/ElvenPlayerComposition"

const SERVICE = ref('ELVEN_PLAYER_C')
const player = reactive(new ElvenPlayerComposition())
const isAudioOverlayActive = ref(false)

const isPlayerActive = toRef(player.audioPlayer, "active")
watch(isPlayerActive, (active) => {
  if (active) {
    player.init()
  } else {
    player.destroy()
  }
})

function onSlideProgress(percents){
  player.setTimeByPercents(percents)
}
</script>

<style scoped>
/*// при установке значений слайдера учитывать марджин*/
/*// также, как и с style width*/
.audio-player-container {
  border-top-right-radius: 6px;
  border-top-left-radius: 6px;
  height: 64px;
  width: 184px;
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  left: 50%;
  bottom: 0;
  transform: translate(-50%, 0);
  margin: 0 auto;
  display: flex;
  flex-direction: row;
  align-items: center;


}

.audio-player-container svg {
  max-height: 20px;
  max-width: 20px;
  height: 20px;
  width: 20px;
  fill: white;
}

.audio-player-item {
  cursor: pointer;
  align-self: center;
}

.playback-tools {
  width: 36px;
}

.show-player {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: row;
  justify-content: center;
}

.playback-tools,
.playback-tools-2 {
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.playback-tools-2 {
  width: 100%;
}

.show-player {
  margin-left: 16px;
}

.playback-controls {
  height: 100%;
  z-index: 5;
  display: flex;
  flex-direction: row;
  gap: 14px;
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
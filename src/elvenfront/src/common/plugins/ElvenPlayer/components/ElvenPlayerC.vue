<template>
  <div id="audio-player-container" v-on:click="isAudioOverlayActive = !isAudioOverlayActive">

    <Overlay :active="isAudioOverlayActive" v-on:deactivated="isAudioOverlayActive = false"></Overlay>

    <div class="audio-player-main" v-on:click.stop>

      <div class="playback-controls">
        <div class="audio-player-item audio-player-prev">
          <svg id="audio-player-prev-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 406.76 169.37">
            <path
                d="M180.76,320.07V179.49A14.38,14.38,0,0,0,159.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,180.76,320.07Z"
                transform="translate(-28 -165.09)"/>
            <path
                d="M434.76,320.07V179.49A14.38,14.38,0,0,0,413.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,434.76,320.07Z"
                transform="translate(-28 -165.09)"/>
          </svg>
        </div>
        <div class="audio-player-item audio-player-play-pause">
          <svg id="audio-player-pause-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 207 310">
            <rect width="70" height="310"/>
            <rect x="137" width="70" height="310"/>
          </svg>
          <svg id="audio-player-play-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 152.76 169.37">
            <path
                d="M173.62,320.29V179.71a14.38,14.38,0,0,1,21.47-12.51l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.38,14.38,0,0,1,173.62,320.29Z"
                transform="translate(-173.62 -165.31)"/>
          </svg>
        </div>

        <div class="audio-player-item audio-player-next">
          <svg id="audio-player-next-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 406.76 169.37">
            <path
                d="M308,320.07V179.49A14.37,14.37,0,0,1,329.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,308,320.07Z"
                transform="translate(-54 -165.09)"/>
            <path
                d="M54,320.07V179.49A14.37,14.37,0,0,1,75.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,54,320.07Z"
                transform="translate(-54 -165.09)"/>
          </svg>
        </div>

      </div>
    </div>

  </div>
</template>

<script>
import Overlay from "./Overlay.vue"

export default {
  name: "ElvenPlayerC",
  components: {Overlay},
  data() {
    return {
      SERVICE: 'ELVEN_PLAYER_C',
      isAudioOverlayActive: false,
    }
  },
}
</script>

<style scoped>
// при установке значений слайдера учитывать марджин
// также, как и с style width
#audio-player-container {
  border-top-right-radius: 6px;
  border-top-left-radius: 6px;
  height: 64px;
  width: 150px;
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  left: 50%;
  bottom: 0;
  transform: translate(-50%, 0);
  margin: 0 auto;
  display: none;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
}

.audio-player-main {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.playback-controls{
  height: 100%;
  z-index: 5;
  display: flex;
  flex-direction: row;
  gap: 14px;
}

.audio-player-prev{

}
#audio-player-play-butt{

}
#audio-player-pause-butt{
}
.audio-player-next{

}

.audio-player-item {
  cursor: pointer;
  align-self: center;
}

.audio-player-item svg {
  max-height: 20px;
  max-width: 20px;
  height: 20px;
  width: 20px;
  fill: white;
}
</style>
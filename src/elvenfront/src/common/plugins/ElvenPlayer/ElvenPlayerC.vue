<template>
  <div class="audio-player-container" v-show="show">
    <div class="audio-player-item audio-player-progress" v-on:click="setDurationOnClickWrap">
      <div class="audio-player-progressbar"></div>
      <div class="audio-player-buffered"></div>
    </div>
    <div class="audio-player-main">
      <div class="audio-player-item audio-player-prev" v-on:click="prevTrackWrap">
        <svg id="prev_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 406.76 169.37"><path d="M180.76,320.07V179.49A14.38,14.38,0,0,0,159.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,180.76,320.07Z" transform="translate(-28 -165.09)"/><path d="M434.76,320.07V179.49A14.38,14.38,0,0,0,413.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,434.76,320.07Z" transform="translate(-28 -165.09)"/></svg>
      </div>
      <div class="audio-player-item audio-player-play-pause" v-on:click="playPauseWrap">
        <svg v-if="isPlaying" id="pause_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 207 310"><rect width="70" height="310"/><rect x="137" width="70" height="310"/></svg>
        <svg v-if="!isPlaying" id="play_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 152.76 169.37"><path d="M173.62,320.29V179.71a14.38,14.38,0,0,1,21.47-12.51l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.38,14.38,0,0,1,173.62,320.29Z" transform="translate(-173.62 -165.31)"/></svg>
      </div>
      <div class="audio-player-item audio-player-next" v-on:click="nextTrackWrap">
        <svg id="next_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 406.76 169.37"><path d="M308,320.07V179.49A14.37,14.37,0,0,1,329.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,308,320.07Z" transform="translate(-54 -165.09)"/><path d="M54,320.07V179.49A14.37,14.37,0,0,1,75.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,54,320.07Z" transform="translate(-54 -165.09)"/></svg>
      </div>
      <div class="audio-player-item audio-player-close" v-on:click="closePlayer">
        <svg id="close_butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 245 270"><path d="M127.21,365h0c-12.2-15.67-11.28-40,2-54.33L308.81,117.62c13.33-14.34,34-13.26,46.24,2.41h0c12.2,15.67,11.28,40-2.05,54.33L173.45,367.44C160.11,381.78,139.41,380.7,127.21,365Z" transform="translate(-118.63 -107.53)"/><path d="M351.24,368.93h0A35.14,35.14,0,0,0,355,319L181.4,119.8A36,36,0,0,0,131,116.13h0a35.15,35.15,0,0,0-3.7,49.9L300.86,365.26A36,36,0,0,0,351.24,368.93Z" transform="translate(-118.63 -107.53)"/></svg>
      </div>
    </div>
  </div>
</template>

<script>
import {methods, state} from "@/common/plugins/ElvenPlayer/core/Shared"

export default {
  name: "ElvenPlayerC",
  data() {
    return {
      SERVICE: 'ELVEN_PLAYER_C',
      show: false,
      sharedState: state,
      isPlaying: false,
    }
  },
  mounted() {
    methods.mounted()
    state.audioPlayerControls.addEventListener('playPause', (event) =>{
      this.isPlaying = event
    })
  },
  methods: {
    showFunc(isShow) {
      this.show = isShow
    },
    playPauseWrap() {
      this.isPlaying = this.sharedState.audioPlayerControls.playPause()
    },
    nextTrackWrap() {
      this.sharedState.audioPlayerControls.next()
    },
    prevTrackWrap() {
      this.sharedState.audioPlayerControls.prev()
    },
    setDurationOnClickWrap(event) {
      this.sharedState.audioPlayerControls.setDurationOnClick(event)
    },
    closePlayer(){
      this.sharedState.audioPlayerControls.stop()
      this.show = false
    },
  },
}
</script>

<style scoped>
.audio-player-container {
  display: flex;
  flex-direction: column;

}

.audio-player-main {
  background-color: rgba(0, 0, 0, 0.5);
  width: 100%;
  height: 54px;
  z-index: 5;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 18px;
}

.audio-player-item {
  cursor: pointer;
}

.audio-player-progress {
  transition: height 0.1s;
  height: 6px;
  position: relative;
  background-color: rgba(0, 0, 0, 0.8);
}

.audio-player-progress:hover {
  height: 24px;
}

.audio-player-progressbar,
.audio-player-buffered {
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 100%;
  transition: width 0.1s;
}

.audio-player-progressbar {
  background-color: #F8F272;
  z-index: 1;
}

.audio-player-buffered {
  background-color: #999999;
}

.audio-player-item svg {
  max-height: 20px;
  max-width: 20px;
  height: 20px;
  width: 20px;
  fill: white;
}




@media screen and (prefers-color-scheme: dark) {
  .audio-player-progressbar{
    background-color: #9D6A89;
  }
  .audio-player-buffered{
    background-color: #999999;
  }
}

@media screen and (prefers-color-scheme: light) {
  .audio-player-progressbar{
    background-color: #9D6A89;
  }
  .audio-player-buffered{
    background-color: #999999;
  }
}
</style>
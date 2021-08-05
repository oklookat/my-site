<template>
<div class="audio-player-container">
  <div class="audio-player-item audio-player-progress" v-on:click="setDurationOnClickWrap">
    <div class="audio-player-progressbar"></div>
  </div>
  <div class="audio-player-main">
    <div class="audio-player-item audio-player-prev" v-on:click="prevTrackWrap">
      <svg id="player_prev_button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 417.11 190.55"><path d="M33.58,245.64a3.76,3.76,0,0,0,0,6.51L191.8,343.5a3.76,3.76,0,0,0,5.64-3.26V157.55a3.76,3.76,0,0,0-5.64-3.26ZM191.87,337.1,39.09,248.89l152.78-88.2Z" transform="translate(-31.7 -153.79)"/><polygon points="7.39 95.1 160.17 183.31 160.17 6.9 7.39 95.1"/><path d="M285,246a3.77,3.77,0,0,0,0,6.52l158.22,91.34a3.72,3.72,0,0,0,3.76,0,3.77,3.77,0,0,0,1.88-3.25V157.88a3.76,3.76,0,0,0-5.64-3.26Zm158.29,91.47L290.46,249.23,443.24,161Z" transform="translate(-31.7 -153.79)"/><polygon points="258.76 95.44 411.54 183.65 411.54 7.23 258.76 95.44"/></svg>
    </div>
    <div class="audio-player-item audio-player-play-pause" v-on:click="playPauseWrap">
      <svg id="player_play_button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 243 278.89"><path d="M369.24,218.06a5.52,5.52,0,0,1,0,9.56l-232,133.92a5.52,5.52,0,0,1-8.27-4.77V88.91a5.52,5.52,0,0,1,8.27-4.78ZM137.17,352.17l224-129.33-224-129.33Z" transform="translate(-129 -83.39)"/><polygon points="232.18 139.44 8.17 268.77 8.17 10.12 232.18 139.44"/></svg>
    </div>
    <div class="audio-player-item audio-player-next" v-on:click="nextTrackWrap">
      <svg id="player_prev_button" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 417.11 190.55"><polygon points="256.94 6.9 409.73 95.1 256.94 183.31 256.94 6.9"/><path d="M461.43,248.89l-152.78-88.2V337.1ZM308.72,343.5a3.77,3.77,0,0,1-5.65-3.26V157.55a3.77,3.77,0,0,1,5.65-3.26l158.21,91.35a3.75,3.75,0,0,1,0,6.51Z" transform="translate(-51.7 -153.79)"/><polygon points="158.35 95.44 5.57 183.65 5.57 7.23 158.35 95.44"/><path d="M57.27,161V337.44l152.79-88.21Zm158.29,91.47L57.34,343.83a3.72,3.72,0,0,1-3.76,0,3.79,3.79,0,0,1-1.88-3.25V157.88a3.76,3.76,0,0,1,5.64-3.26L215.56,246a3.77,3.77,0,0,1,0,6.52Z" transform="translate(-51.7 -153.79)"/></svg>
    </div>
  </div>
</div>
</template>

<script>
import {methods, state} from "@/common/plugins/ElvenPlayer/core/Shared";

export default {
  name: "ElvenPlayerC",
  data() {
    return {
      sharedState: state,
    }
  },
  mounted() {
    methods.mounted()
  },
  methods: {
    playPauseWrap(){
      this.sharedState.audioPlayer.playPause()
    },
    nextTrackWrap(){
      this.sharedState.audioPlayer.next()
    },
    prevTrackWrap(){
      this.sharedState.audioPlayer.prev()
    },
    setDurationOnClickWrap(event){
      this.sharedState.audioPlayer.setDurationOnClick(event)
    },
    test(){
      console.log('MOUSE UP')
    }
  },
}
</script>

<style scoped>
.audio-player-container {
 display: flex;
  flex-direction: column;

}
.audio-player-main{
  background-color: var(--color-level-1);
  width: 100%;
  height: 54px;
  z-index: 5;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 18px;
}
.audio-player-progress{
  height: 6px;
}
.audio-player-progressbar{
  height: 100%;
  width: 0;
  background-color: red;
  transition: width 0.8s;
}
.audio-player-item {
  cursor: pointer;
}
.audio-player-item svg {
  max-height: 25px;
  max-width: 25px;
  height: 25px;
  width: 25px;
}
.audio-player-play-pause {
}

@media screen and (prefers-color-scheme: light) {
  .audio-player-item svg {
    fill: black;
  }
}
@media screen and (prefers-color-scheme: dark) {
  .audio-player-item svg {
    fill: white;
  }
}
</style>
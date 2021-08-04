<template>
<div class="audio-player-container">
  <div class="audio-player" v-on:click="playPauseWrap">PLAY / PAUSE</div>
</div>
</template>

<script>
import ElvenPlayerCore from "@/common/plugins/ElvenPlayer/ElvenPlayerCore"
export default {
  name: "ElvenPlayerC",
  data() {
    return {
      // service vars start //
      SERVICE: 'ELVEN_PLAYER_C',
      audioPlayer: undefined,
      addSource: {default: 'null', type: String},
      // service vars end //
      sources: [],
      currentSource: '',
    }
  },
  // mounted() {
  //   this.setSourceFunc()
  // },
  // unmounted() {
  //
  // },
  watch: {
    // wraps from app.$elvenPlayer
    // setSources: function () {
    //
    // },
    addSource: function () {
      this.currentSource = this.addSource
      this.sources.push(this.currentSource)
      this.setSourceFunc()
    },
  },
  methods: {
    setSourceFunc(){
      this.audioPlayer = new ElvenPlayerCore(this.currentSource)
    },
    playPauseWrap(){
      this.audioPlayer.playPause()
    },
  },
}
</script>

<style scoped>
.audio-player-container{
  background-color: red;
  width: 100%;
  height: 54px;
  z-index: 5;
}
.audio-player {
  width: 100%;
  height: 100%;
}
.audio-player:hover{
  background-color: rgba(0, 0, 0, 0.2);
}
</style>
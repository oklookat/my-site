<template>
  <div class="audio-overlay-container" ref="audioOverlayContainer">
    <div class="overlay">
      <div class="content">
        <div class="audio-overlay">

          <div class="player-controls">
            <div class="close-overlay" v-on:click="deactivate">
              <svg data-name="close-butt" xmlns="http://www.w3.org/2000/svg"
                   viewBox="0 0 255.07 295.91">
                <path d="M135,390.18h0a33.69,33.69,0,0,0,48-4.29L369.93,159.34a35.1,35.1,0,0,0-4.19-48.86h0a33.69,33.69,0,0,0-48,4.29L130.83,341.33A35.09,35.09,0,0,0,135,390.18Z" transform="translate(-122.85 -102.37)"/>
                <path d="M365.74,390.18h0a33.68,33.68,0,0,1-48-4.29L130.83,159.34A35.1,35.1,0,0,1,135,110.48h0a33.7,33.7,0,0,1,48,4.29L369.93,341.33A35.09,35.09,0,0,1,365.74,390.18Z" transform="translate(-122.85 -102.37)"/>
              </svg>
            </div>
          </div>

          <div class="cover">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 300 300">
              <g id="back">
                <rect width="300" height="300" rx="18.75"/>
              </g>
              <g id="melody">
                <rect x="78" y="163" width="90" height="80" rx="25" style="fill:#fcfcfc"/>
                <rect x="132.64" y="56" width="35" height="166" rx="11.5" style="fill:#fcfcfc"/>
                <rect x="133" y="56" width="100" height="40" rx="20" style="fill:#fcfcfc"/>
              </g>
            </svg>
          </div>

          <div class="time-control">
            <div class="time-sliders">
              <Progress class="time-buffered"
                        v-bind:percents="player.audioPlayer.currentPlaying.percentsBuffered"></Progress>
              <Slider class="time-slider"
                      :emit-after-up="true"
                      v-on:slide="onProgressSliderTriggered"
                      v-bind:percents="player.audioPlayer.currentPlaying.percentsReached">
              </Slider>
            </div>
            <div class="time-duration">
              <div class="time-current">
                {{ player.audioPlayer.currentPlaying.currentTime }}
              </div>
              <div class="time-total">
                {{ player.audioPlayer.currentPlaying.duration }}
              </div>
            </div>
          </div>

          <div class="volume-control">
            <div class="volume-sliders">
              <Slider
                  :percents="player.audioPlayer.percentsVolume"
                  :emit-after-up="false"
                  v-on:slide="onVolumeSliderTriggered">
              </Slider>
            </div>
          </div>


          <PlaybackControls class="playback-controls" :player="player">
          </PlaybackControls>


        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Slider from "../Slider/Component.vue"
import Progress from '../Progress/Component.vue'
import {onMounted, onUnmounted, ref, toRefs, watch} from "vue"
import Composition from "@/common/plugins/ElvenPlayer/components/Player/Composition";
import PlaybackControls from "@/common/plugins/ElvenPlayer/components/Player/PlaybackControls.vue";

const props = defineProps({
  player: Composition,
})

const {player} = toRefs(props)

const emit = defineEmits<{
  (e: 'deactivated'): void
  (e: 'slideProgress', percents: number): void
  (e: 'slideVolume', percents: number): void
}>()

const noScroll = ref(false)
watch(noScroll, () => {
  switchScroll()
})

const audioOverlayContainer = ref(null)
defineExpose({audioOverlayContainer})
onMounted(() => {
  document.body.appendChild(audioOverlayContainer.value)
  noScroll.value = true
})

onUnmounted(() => {
  document.body.style.overflow = null
})


function switchScroll() {
  if (noScroll.value) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = null
  }
}

function deactivate() {
  noScroll.value = false
  emit('deactivated')
}

function onProgressSliderTriggered(percents) {
  player.value.setTimeByPercents(percents)
}

function onVolumeSliderTriggered(percents) {
  player.value.setVolumeByPercents(percents)
}
</script>

<style scoped>
.audio-overlay-container,
.overlay{
  display: flex;
  align-items: center;
  justify-content: center;
}

.audio-overlay-container {
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
}

.overlay {
  background-color: white;
  border-radius: 12px;
  z-index: 9999;
  overflow: auto;
  max-width: 365px;
  max-height: 365px;
  width: 294px;
  height: 524px;
}

.content {
  box-sizing: border-box;
  font-size: 1rem;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.audio-overlay {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  align-self: center;
  gap: 18px;
}

.player-controls{
  height: 19px;
  margin-left: auto;
  margin-right: 12px;
  margin-top: 4px;
  display: flex;
  flex-direction: row;
  gap: 16px;
}

.close-overlay{
  cursor: pointer;
}
.close-overlay,
.close-overlay svg,
.close-player svg{
  width: 15px;
  height: 15px;
}

.cover {
  height: 150px;
  width: 150px;
}

.time-control,
.volume-control{
  height: max-content;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.time-control {
  width: 80%;
}

.time-sliders,
.volume-sliders{
  position: relative;
  border-radius: 4px;
  background-color: rgba(0, 0, 0, 0.5);
  height: 18px;
  width: 100%;
}

.time-slider,
.time-buffered {
  border-radius: inherit;
  position: absolute;
  width: 100%;
  height: 100%;
}

.time-duration {
  display: flex;
  flex-direction: row;
}

.time-total {
  margin-left: auto;
}

.volume-control{
  height: 24px;
  width: 50%;
}

.playback-controls{
  width: 128px;
  height: 30px;
  fill: black;
  justify-content: center;
}
</style>
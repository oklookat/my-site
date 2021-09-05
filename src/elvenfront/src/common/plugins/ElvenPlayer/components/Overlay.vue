<template>
  <div class="audio-overlay-container" v-on:click.self="deactivate" ref="audioOverlayContainer">
    <div class="overlay">
      <div class="content">

        <div class="audio-overlay">
          <div class="audio-player-item time-container">
            <Progress class="time-buffered"
                v-bind:percents="percentsBuffered"></Progress>
            <Slider class="time-slider"
                    :emit-after-up="true"
                    v-on:slide="onProgressSliderTriggered"
                    v-bind:percents="percentsTime">
            </Slider>
          </div>

          <div class="audio-player-volume-control">
            <!--        <div class="audio-player-mute">Mute</div>-->
            <div id="volume-slider-container">
              <!--            <Slider class="volume-slider"></Slider>-->
            </div>
          </div>

          <div class="audio-player-item audio-player-close">
            <div id="audio-player-close-butt">закрыть плеер</div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Slider from "./Slider/Component.vue"
import Progress from './Progress/Component.vue'
import {onMounted, onUnmounted, ref, toRefs, watch} from "vue"

const props = defineProps({
  percentsTime: Number,
  percentsBuffered: Number,
})

const emit = defineEmits<{
  (e: 'deactivated'): void
  (e: 'slideProgress', percents: number): void
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

function onProgressSliderTriggered(percents){
  emit('slideProgress', percents)
}
</script>

<style scoped>
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
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.overlay {
  background-color: white;
  border-radius: 12px;
  z-index: 9999;
  overflow: auto;
  max-width: 365px;
  max-height: 365px;
  width: 250px;
  height: 250px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  border: 1px solid #E6E6E6;
}

.content {
  box-sizing: border-box;
  font-size: 1.2rem;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.audio-overlay {
  width: 100%;
  height: 100%;
  margin-top: 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  align-self: center;
  gap: 24px;
}

.time-container {
  position: relative;
  margin-top: 48px;
  border-radius: 4px;
  background-color: rgba(0, 0, 0, 0.5);
  height: 16px;
  width: 164px;
}

.time-slider,
.time-buffered{
  position: absolute;
  width: 100%;
  height: 100%;
}
</style>
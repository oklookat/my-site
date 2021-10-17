<template>
  <div class="player__controls">
    <svg
      v-on:click="prev"
      class="player__controls-item player__prev-butt"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 406.76 169.37"
    >
      <path
        d="M180.76,320.07V179.49A14.38,14.38,0,0,0,159.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,180.76,320.07Z"
        transform="translate(-28 -165.09)"
      />
      <path
        d="M434.76,320.07V179.49A14.38,14.38,0,0,0,413.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,434.76,320.07Z"
        transform="translate(-28 -165.09)"
      />
    </svg>
    <div class="player__controls-item player__play-pause">
      <svg
        v-if="!player?.audioPlayer.isPlaying"
        v-on:click="play"
        class="player__play-butt"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 152.76 169.37"
      >
        <path
          d="M173.62,320.29V179.71a14.38,14.38,0,0,1,21.47-12.51l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.38,14.38,0,0,1,173.62,320.29Z"
          transform="translate(-173.62 -165.31)"
        />
      </svg>
      <svg
        v-if="player?.audioPlayer.isPlaying"
        v-on:click="pause"
        class="player__controls-item player__pause-butt"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 207 310"
      >
        <rect width="70" height="310" />
        <rect x="137" width="70" height="310" />
      </svg>
    </div>
    <svg
      v-on:click="next"
      class="player__controls-item player__next-butt"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 406.76 169.37"
    >
      <path
        d="M308,320.07V179.49A14.37,14.37,0,0,1,329.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,308,320.07Z"
        transform="translate(-54 -165.09)"
      />
      <path
        d="M54,320.07V179.49A14.37,14.37,0,0,1,75.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,54,320.07Z"
        transform="translate(-54 -165.09)"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
import Composition from "./Composition"

const props = defineProps({
  player: Composition,
})

function next() {
  props.player?.next()
}

function prev() {
  props.player?.prev()
}

function play() {
  props.player?.play()
}

function pause() {
  props.player?.pause()
}
</script>

<style scoped>
.player__controls-item {
  cursor: pointer;
}
.player__controls {
  z-index: 5;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 14px;
}
.player__controls svg {
  height: 20px;
  width: 20px;
  max-height: 20px;
  max-width: 20px;
}
</style>
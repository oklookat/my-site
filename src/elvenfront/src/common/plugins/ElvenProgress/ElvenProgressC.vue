<template>
  <div class="progressbar__container">
    <div id="progressbar__line" :style="{ width: `${percents}%`, height: height}"></div>
  </div>
</template>

<script setup lang="ts">
import { onUnmounted, ref } from '@vue/runtime-core'
const SERVICE = 'ELVEN_PROGRESS_C'

//// settings
let height = ref('3px')
let moveSpeed = 100
let loadingStartTo = 45
let loadingStartSpeed = 30
let loadingFinishSpeed = 5

let percents = ref(0)

onUnmounted(() => {
  destroy()
})

function destroy() {
  percents.value = 0
}

function move() {
  const intervalID = setInterval(() => {
    if (percents.value < 100) {
      percents.value++
    } else {
      clearInterval(intervalID)
    }
  }, moveSpeed)
}

function loading(start: boolean) {
  if (start) {
    const intervalID = setInterval(() => {
      // freeze progress at loadingStartTo
      if (percents.value < loadingStartTo) {
        percents.value++
      } else {
        clearInterval(intervalID)
      }
    }, loadingStartSpeed)
    return
  }
  // finish (go to 100 and destroy)
  percents.value = loadingStartTo
  const intervalID = setInterval(() => {
    if (percents.value < 100) {
      percents.value++
    } else {
      clearInterval(intervalID)
      destroy()
    }
  }, loadingFinishSpeed)
  return
}
</script>

<style scoped>
.progressbar__container {
  cursor: default;
  position: absolute;
  width: 100%;
  height: 6px;
}

#progressbar__line {
  height: 100%;
  width: 0;
  background-color: #a097dc;
}
</style>
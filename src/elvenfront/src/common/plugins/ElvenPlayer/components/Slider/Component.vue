<template>
  <div class="slider-container" ref="sliderContainer">
    <div class="sliders" v-if="ready">
      <div class="sl-slider"
           :style="{width: `${elementPercents}`}">
      </div>
      <div
          class="sl-slider-bubble"
          :style="{left: `calc(${elementPercents} - 3px)`}">
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, onUnmounted, reactive, ref, toRefs, watch, toRef} from "vue"
import Composition from "@/common/plugins/ElvenPlayer/components/Slider/Composition"

interface IProps {
  emitAfterUp: boolean // send percents only after mouse up
  percents?: number // percents in slider
}

const props = withDefaults(defineProps<IProps>(), {
  emitAfterUp: false,
  percents: 0
})

const emit = defineEmits<{
  (e: 'slide', percents: number): void
}>()

const {emitAfterUp, percents} = toRefs(props)
// if true - we can render slider data in component
const ready = ref(false)
// slider container element used by composition class. Need for set events.
const sliderContainer = ref(null)
defineExpose({sliderContainer})
const slide = reactive(new Composition())
// percents by slider dragging be here
const slidePercents = ref(0)

onMounted(() => {
  if (percents.value > 0) {
    calcStyles(percents.value)
  }
  slide.init(sliderContainer.value)
  ready.value = true
})

onUnmounted(() => {
  slide.destroy()
})

// convert slider values to ref for using in watch
const isSlideMouseDown = toRef(slide.slider, 'isMouseDown')
const percentsChanged = toRef(slide.slider, 'percents')
// watch changes (like dragging) for emit percents and set slider style
watch([isSlideMouseDown, percentsChanged], (newValues) => {
  slidePercents.value = newValues[1]
  const isSlideMouseUp = !isSlideMouseDown.value
  if (emitAfterUp.value) {
    // emit only when mouse up after dragging
    if (isSlideMouseUp) {
      // if mouse up after slider dragging
      emit('slide', slidePercents.value)
    }
  } else {
    // every drag - emit
    emit('slide', slidePercents.value)
  }
  calcStyles(slidePercents.value)
})

// on slider percents changed by prop
watch(percents, newValue => {
  if (isSlideMouseDown.value) {
    // block change percents by prop, because now user dragging slider
    return
  }
  calcStyles(newValue)
})

// used only for style slider
const elementPercents = ref('0%')

// set style on slider element
function calcStyles(percents: number) {
  elementPercents.value = `${percents}%`
}
</script>

<style scoped>
.slider-container {
  width: 100%;
  height: 100%;
}

.sliders {
  width: calc(100% - 3px);
  height: 100%;
  position: relative;
  display: flex;
  align-items: center;
}

.slider-container,
.sliders,
.sl-slider,
.sl-slider-bubble {
  cursor: pointer;
}

.sl-slider {
  position: absolute;
  border-top-left-radius: 2px;
  border-bottom-left-radius: 2px;
  width: 0;
  height: 100%;
  background-color: #9D6A89;
  top: auto;
  bottom: 0;
}

.sl-slider-bubble {
  position: absolute;
  background-color: #9D6A89;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  box-sizing: border-box;
  left: 0;
  margin-left: -6px;
  transition: width, height 100ms;
}
</style>
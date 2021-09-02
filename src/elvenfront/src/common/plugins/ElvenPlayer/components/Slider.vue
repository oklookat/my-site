<template>
  <div class="slider-container">
    <div class="sl-buffered" v-if="withBuffered"></div>
    <div class="sl-slider"></div>
    <div class="sl-slider-bubble"></div>
  </div>
</template>

<script>
export default {
  name: "Slider",
  props: {
    withBuffered: false
  }
}
</script>

<style scoped>
.slider-container{
  position: relative;
  width: calc(100% - 3px);
  height: 100%;
  display: flex;
  align-items: center;
}

.slider-container,
.sl-slider,
.sl-slider-bubble{
  cursor: pointer;
}

.sl-slider{
  position: absolute;
  border-top-left-radius: 2px;
  border-bottom-left-radius: 2px;
  width: 0;
  height: 100%;
  background-color: #9D6A89;
  top: auto;
  bottom: 0;
}

.sl-slider-bubble{
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

.sl-slider-bubble:hover{
  width: 25px;
  height: 25px;
}
</style>
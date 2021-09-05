<template>
<div class="progress-container">
  <div class="ps-progress" :style="{width: `${percents}%`}">
  </div>
</div>
</template>

<script setup lang="ts">
const props = defineProps({
  percents: Number
})
</script>

<style scoped>
.progress-container{
  width: 100%;
  height: 100%;
}

.ps-progress{
  height: 100%;
  width: 0;
  background-color: #404040;
}
</style>
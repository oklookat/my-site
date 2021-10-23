
<div class="progress__container">
  <div class="progress__line" style="width: {percents}%">
  </div>
</div>


<script lang="ts">
export let percents: number
</script>

<style scoped>
.progress__container{
  border-radius: inherit;
  width: 100%;
  height: 100%;
}

.progress__line{
  border-radius: inherit;
  height: 100%;
  width: 0;
  background-color: #404040;
}
</style>
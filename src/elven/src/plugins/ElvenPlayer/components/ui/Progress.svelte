<script lang="ts">
  export let percents: number;
</script>

<div class="progress">
  <div class="progress__line" style="width: {percents}%" />
</div>

<style lang="scss">
  .progress {
    border-radius: inherit;
    width: 100%;
    height: 100%;
    &__line {
      border-radius: inherit;
      height: 100%;
      width: 0;
      background-color: #404040;
    }
  }
</style>

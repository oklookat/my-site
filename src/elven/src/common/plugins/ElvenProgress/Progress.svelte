<script lang="ts">
  import { onDestroy } from "svelte";

  // controlled by plugin
  export let isLoading = false;
  $: watchLoading(isLoading);
  function watchLoading(value) {
    loading(value);
  }

  //// settings
  let height = "3px";
  let moveSpeed = 100;
  let loadingStartTo = 45;
  let loadingStartSpeed = 30;
  let loadingFinishSpeed = 5;
  // main
  let percents = 0;

  onDestroy(() => {
    destroy();
  });

  function destroy() {
    percents = 0;
  }

  function loading(start: boolean) {
    if (start) {
      const intervalID = setInterval(() => {
        // freeze progress at loadingStartTo
        if (percents < loadingStartTo) {
          percents++;
        } else {
          clearInterval(intervalID);
        }
      }, loadingStartSpeed);
      return;
    }
    // finish (go to 100 and destroy)
    percents = loadingStartTo;
    const intervalID = setInterval(() => {
      if (percents < 100) {
        percents++;
      } else {
        clearInterval(intervalID);
        destroy();
      }
    }, loadingFinishSpeed);
    return;
  }
</script>

<div class="progressbar__container">
  <div id="progressbar__line" style="width: {percents}%; height: {height}" />
</div>

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

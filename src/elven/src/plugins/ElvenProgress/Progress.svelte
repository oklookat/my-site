<script lang="ts">
  import type { IElvenProgress } from "@/plugins/ElvenProgress/types";
  import { onDestroy } from "svelte";

  // plugin controls
  class Plugin implements IElvenProgress {
    public startBasic() {
      startBasic();
    }
    public finishBasic() {
      finishBasic();
    }
    public setPercents(value: number) {
      percents = value;
    }
    public resetPercents() {
      percents = 0;
      destroy();
    }
  }
  window.$elvenProgress = new Plugin();

  // settings
  let height = "3px";
  let loadingStartTo = 45;
  let loadingStartSpeed = 30;
  let loadingFinishSpeed = 5;

  // element percents
  let percents = 0;

  onDestroy(() => {
    destroy();
  });

  function destroy() {
    percents = 0;
  }

  // freeze progress at loadingStartTo
  function startBasic() {
    const intervalID = setInterval(() => {
      if (percents < loadingStartTo) {
        percents++;
      } else {
        clearInterval(intervalID);
      }
    }, loadingStartSpeed);
  }

  // finish (go to 100 and destroy)
  function finishBasic() {
    percents = loadingStartTo;
    const intervalID = setInterval(() => {
      if (percents < 100) {
        percents++;
      } else {
        clearInterval(intervalID);
        destroy();
      }
    }, loadingFinishSpeed);
  }
</script>

<div class="progressbar__container">
  <div class="progressbar__line" style="width: {percents}%; height: {height}" />
</div>

<style lang="scss">
  .progressbar {
    &__container {
      cursor: default;
      position: absolute;
      width: 100%;
      height: 6px;
    }
    &__line {
      height: 100%;
      width: 0;
      background-color: #a097dc;
    }
  }
</style>

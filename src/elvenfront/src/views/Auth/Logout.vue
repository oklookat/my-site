<template>
  <div class="logout">
    <div class="logout-cube-wrap">
      <div class="logout-cube">
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import AuthAdapter from "@/common/adapters/Main/AuthAdapter";

export default defineComponent({
  name: 'Logout',
  async mounted() {
    await AuthAdapter.logout()
  }
})
</script>

<style scoped lang="scss">
.logout{
  height: inherit;
  width: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
}

.logout-cube-wrap {
  width: inherit;
  height: inherit;
  perspective: 1000px;
  perspective-origin: 50% 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-cube {
  display: flex;
  background-color: var(--color-text);
  width: 15vh;
  height: 15vh;
  transform: scaleX(1) scaleY(1) scaleZ(1) rotateX(0deg) rotateY(0deg) rotateZ(0deg) translateX(0px) translateY(0px) translateZ(0px) skewX(0deg) skewY(0deg);
  animation: loading infinite 2s;
}

@keyframes loading {
  0%{
    transform: rotateY(0deg)
  }
  100%{
    transform: rotateY(180deg)
  }
}

</style>
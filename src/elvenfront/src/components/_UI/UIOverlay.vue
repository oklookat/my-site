<template>
  <div class="ui-overlay-container" v-if="active" v-on:click.self="deactivate">
    <div class="ui-overlay">
      <div class="ui-overlay-content">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue";

export default defineComponent({
  name: 'UIOverlay',
  unmounted() {
    document.body.classList.remove('no-scroll')
  },
  data(){
    return{
    }
  },
  props: {
    title: String,
    active: Boolean,
  },
  watch: {
    active: function () {
      if(this.active){
        document.body.classList.add('no-scroll')
      } else{
        document.body.classList.remove('no-scroll')
      }
    }
  },
  methods: {
    deactivate(){
      document.body.classList.remove('no-scroll')
      this.$emit('deactivated')
    },
  },
})
</script>

<style scoped lang="scss">
.ui-overlay-container{
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: rgba(0,0,0,0.5);
  z-index: 9998;
  display: flex;
  flex-direction: column;
}
.ui-overlay{
  z-index: 9999;
  background-color: var(--color-level-1);
  border-radius: 12px;
  position: absolute;
  transform: translate(-50%,-50%);
  top: 50%;
  left: 50%;
  margin-left: auto;
  margin-right: auto;
  width: 80%;
  height: 50%;
  max-width: 365px;
  max-height: 365px;
  overflow: auto;
}
.ui-overlay-content{
  margin: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  font-size: 1.2rem;
}
</style>
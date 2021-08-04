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

<style scoped>
.ui-overlay-container{
  background-color: rgba(0,0,0,0.7);
  z-index: 9998;
  max-width: 100vw;
  width: 100%;
  height: 100%;
  position: fixed;
  top: 0; right: 0; bottom: 0; left: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.ui-overlay{
  background-color: var(--color-body);
  border-radius: 12px;
  z-index: 9999;
  overflow: auto;
  max-width: 365px;
  max-height: 365px;
  width: 90%;
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  border: var(--color-border) 1px solid;
}
.ui-overlay-content{
  box-sizing: border-box;
  /*padding-top: 12px;*/
  /*margin-left: 12px;*/
  /*margin-right: 12px;*/
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  font-size: 1.2rem;
  height: 100%;
  width: 100%;
}
</style>
<template>
<div class="audio-overlay-container" v-show="active" v-on:click.self="deactivate">
  <div class="overlay">
    <div class="content">
      <div class="audio-overlay">

        <div class="audio-player-item audio-player-progress-container">
          <Slider class="playback-progressbar" :with-buffered="true"></Slider>
        </div>

        <div class="audio-player-volume-control">
          <!--        <div class="audio-player-mute">Mute</div>-->
          <div id="volume-slider-container">
            <Slider class="volume-slider"></Slider>
          </div>
        </div>

        <div class="audio-player-item audio-player-close">
          <div id="audio-player-close-butt">закрыть плеер</div>
        </div>

      </div>
    </div>
  </div>
</div>
</template>

<script>
import Slider from "./Slider.vue"
export default {
  name: "Overlay",
  components: {Slider},
  unmounted() {
    document.body.classList.remove('no-scroll')
  },
  data(){
    return{
      scroll: false,
    }
  },
  props: {
    active: Boolean,
  },
  watch: {
    active: function () {
      this.scroll = this.active
    },
    scroll: function (){
      this.switchScroll()
    }
  },
  methods: {
    switchScroll(){
      if(this.scroll){
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = null
      }
    },
    deactivate(){
      this.scroll = false
      this.$emit('deactivated')
    },
  },
  mounted() {
    document.body.appendChild(this.$el)
  }
}
</script>

<style scoped>
.audio-overlay-container{
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
.overlay{
  background-color: white;
  border-radius: 12px;
  z-index: 9999;
  overflow: auto;
  max-width: 365px;
  max-height: 365px;
  width: 250px;
  height: 250px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  border: 1px solid #E6E6E6;
}

.content{
  box-sizing: border-box;
  font-size: 1.2rem;
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.audio-overlay{
  width: 100%;
  height: 100%;
  margin-top: 14px;
  display: flex;
  flex-direction: column;
  align-items: center;
  align-self: center;
  gap: 24px;
}

.audio-player-progress-container {
  margin-top: 48px;
  border-radius: 4px;
  background-color: rgba(0, 0, 0, 0.5);
  height: 16px;
  width: 90%;
}


#volume-slider-container{
  width: 164px;
  height: 12px;
}
</style>
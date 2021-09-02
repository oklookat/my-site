<template>
  <div id="audio-player-container">
    <div class="audio-player-item" id="audio-player-progress-container">
      <div id="audio-player-progressbar"></div>
      <div id="audio-player-buffered"></div>
    </div>
    <div class="audio-player-main">
      <div class="audio-player-volume-control">
<!--        <div class="audio-player-mute">Mute</div>-->
        <div class="volume-slider-container" id="volume-container">
          <div id="volume-slider"></div>
          <div id="volume-slider-bubble"></div>
        </div>
      </div>
      <div class="audio-player-item" id="audio-player-prev">
        <svg id="audio-player-prev-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 406.76 169.37">
          <path
              d="M180.76,320.07V179.49A14.38,14.38,0,0,0,159.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,180.76,320.07Z"
              transform="translate(-28 -165.09)"/>
          <path
              d="M434.76,320.07V179.49A14.38,14.38,0,0,0,413.29,167l-124,70.29a14.38,14.38,0,0,0,0,25l124,70.29A14.38,14.38,0,0,0,434.76,320.07Z"
              transform="translate(-28 -165.09)"/>
        </svg>
      </div>
      <div class="audio-player-item" id="audio-player-play-pause">
        <svg id="audio-player-pause-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 207 310">
          <rect width="70" height="310"/>
          <rect x="137" width="70" height="310"/>
        </svg>
        <svg id="audio-player-play-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 152.76 169.37">
          <path
              d="M173.62,320.29V179.71a14.38,14.38,0,0,1,21.47-12.51l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.38,14.38,0,0,1,173.62,320.29Z"
              transform="translate(-173.62 -165.31)"/>
        </svg>
      </div>
      <div class="audio-player-item" id="audio-player-next">
        <svg id="audio-player-next-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 406.76 169.37">
          <path
              d="M308,320.07V179.49A14.37,14.37,0,0,1,329.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,308,320.07Z"
              transform="translate(-54 -165.09)"/>
          <path
              d="M54,320.07V179.49A14.37,14.37,0,0,1,75.46,167l124,70.29a14.38,14.38,0,0,1,0,25l-124,70.29A14.37,14.37,0,0,1,54,320.07Z"
              transform="translate(-54 -165.09)"/>
        </svg>
      </div>
      <div class="audio-player-item" id="audio-player-close">
        <svg id="audio-player-close-butt" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 245 270" style="transform: scale(0.8)">
          <path
              d="M127.21,365h0c-12.2-15.67-11.28-40,2-54.33L308.81,117.62c13.33-14.34,34-13.26,46.24,2.41h0c12.2,15.67,11.28,40-2.05,54.33L173.45,367.44C160.11,381.78,139.41,380.7,127.21,365Z"
              transform="translate(-118.63 -107.53)"/>
          <path
              d="M351.24,368.93h0A35.14,35.14,0,0,0,355,319L181.4,119.8A36,36,0,0,0,131,116.13h0a35.15,35.15,0,0,0-3.7,49.9L300.86,365.26A36,36,0,0,0,351.24,368.93Z"
              transform="translate(-118.63 -107.53)"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ElvenPlayerC",
  data() {
    return {
      SERVICE: 'ELVEN_PLAYER_C',
    }
  },
}
</script>

<style scoped>
#audio-player-container {
  display: none;
  flex-direction: column;
}

.audio-player-main {
  background-color: rgba(0, 0, 0, 0.5);
  width: 100%;
  height: 54px;
  z-index: 5;

  display: grid;
  grid-template-columns: 1fr repeat(3, auto) 1fr;
  grid-column-gap: 5px;
  grid-gap: 14px;
  justify-items: center;
  align-items: center;
}

.audio-player-item {
  cursor: pointer;
}

.audio-player-item:nth-child(1) {
  grid-column-start: 2;
}

.audio-player-item:nth-child(4) {
  margin-left: auto;
  padding-right: 24px;
}

#audio-player-progress-container {
  transition: height 0.1s;
  height: 6px;
  position: relative;
  background-color: rgba(0, 0, 0, 0.8);
}

#audio-player-progress-container:hover {
  height: 24px;
}

#audio-player-progressbar,
#audio-player-buffered {
  position: absolute;
  top: 0;
  left: 0;
  width: 0;
  height: 100%;
}

#audio-player-progressbar {
  background-color: #9D6A89;
  z-index: 1;
}

#audio-player-buffered {
  background-color: #999999;
}

.volume-slider-container{
  position: relative;
  width: 164px;
  height: 12px;
  display: flex;
  align-items: center;
}
#volume-slider{
  position: absolute;
  border-radius: 2px;
  width: 100%;
  height: 100%;
  background-color: red;
  top: auto;
  bottom: 0;
}
#volume-slider-bubble{
  position: absolute;
  background-color: blueviolet;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  left: 100%;
  margin-left: -3px;
}

.audio-player-item svg {
  max-height: 20px;
  max-width: 20px;
  height: 20px;
  width: 20px;
  fill: white;
}


@media screen and (max-width: 912px){
  #audio-player-progress-container {
    height: 24px;
  }
}
</style>
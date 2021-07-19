<template>
  <div class="elven-progressbar">
    <div id="elven-progressbar-line">
    </div>
  </div>
</template>

<script>
export default {
  name: "ElvenProgressC",
  data() {
    return {
      // service vars start //
      SERVICE: 'ELVEN_PROGRESS_C',
      progressBar: undefined,
      isProgressBarActive: true,
      basicLoading: false, // wrap
      setPercents: 0, // wrap
      percents: 0,
      // service vars end //

      //// SETTINGS START ////
      // default settings
      progressBarHeight: '2px',
      progressBarColor: 'blue',
      moveSpeed: 100,
      basicLoadingStartSpeed: 30,
      basicLoadingFinishSpeed: 5,
      basicLoadingStartTo: 45,
      //// SETTINGS END ////
    }
  },
  watch: {
    // wraps from app.$elvenProgress
    setPercents: function () {
      this.percents = this.setPercents
      this.progressBar.style.width = `${this.percents}%`
    },
    basicLoading: function () {
      this.basicLoadingFunc(this.basicLoading)
    },
  },
  mounted() {
    this.progressBar = document.getElementById('elven-progressbar-line')
    this.openBar()

    // user settings are applied here
    this.progressBar.style.height = this.progressBarHeight
    this.progressBar.style.backgroundColor = this.progressBarColor
  },
  unmounted() {
    this.closeBar()
  },
  methods: {
    openBar() {
      this.percents = 0
      this.progressBar.style.width = '0'
      this.isProgressBarActive = true
    },
    closeBar() {
      this.percents = 0
      this.progressBar.style.width = '0'
      this.isProgressBarActive = false
      this.basicLoading = false
    },
    move() {
      const intervalID = setInterval(() => {
        if (this.percents < 100) {
          this.percents++
          this.progressBar.style.width = `${this.percents}%`
        } else {
          clearInterval(intervalID)
        }
      }, this.moveSpeed)
    },

    // BASIC LOADING START //
    basicLoadingFunc(isStart) {
      if (isStart) {
        const intervalID = setInterval(() => {
          if (this.percents < this.basicLoadingStartTo) {
            this.percents++
            this.progressBar.style.width = `${this.percents}%`
          } else {
            clearInterval(intervalID)
            // we don't not close bar, because waiting for user call the finish function
          }
        }, this.basicLoadingStartSpeed)
      } else {
        // finish function
        this.percents = this.basicLoadingStartTo
        const intervalID = setInterval(() => {
          if (this.percents < 100) {
            this.percents++
            this.progressBar.style.width = `${this.percents}%`
          } else {
            clearInterval(intervalID)
            this.closeBar()
          }
        }, this.basicLoadingFinishSpeed)
      }
    },
    // BASIC LOADING END //
  }
}
</script>

<style scoped>
.elven-progressbar {
  cursor: default;
  position: absolute;
  width: 100%;
  height: 2px;
}

#elven-progressbar-line {
  height: 100%;
  width: 0;
  background-color: red;
}
</style>
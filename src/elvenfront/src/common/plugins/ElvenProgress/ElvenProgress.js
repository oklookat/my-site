import ElvenProgressC from './ElvenProgressC.vue'

class ElvenProgress {

    static componentData = null

    static install(app, options) {
        app.component('elven-progress', ElvenProgressC)
        const _this = this
        app.mixin({
            created() {
                if (this.SERVICE === 'ELVEN_PROGRESS_C') {
                    _this.componentData = this
                    theLogic.options = options
                    theLogic.init()
                    window.$elvenProgress = theLogic
                }
            },
        })
    }
}

class theLogic {

    static options = null

    static init() {
        if (this.options) {
            if (this.options.height) {
                ElvenProgress.componentData.height = this.options.height
            }
            if (this.options.loadingStartSpeed) {
                ElvenProgress.componentData.loadingStartSpeed = this.options.loadingStartSpeed
            }
            if (this.options.loadingFinishSpeed) {
                ElvenProgress.componentData.loadingFinishSpeed = this.options.loadingFinishSpeed
            }
            if (this.options.loadingStartTo) {
                ElvenProgress.componentData.loadingStartTo = this.options.loadingStartTo
            }
        }
    }

    static setPercents(percents) {
        ElvenProgress.componentData.setPercents(percents)
    }

    static loadingStart() {
        ElvenProgress.componentData.loading(true)
    }

    static loadingFinish() {
        ElvenProgress.componentData.loading(false)
    }

    static close() {
        ElvenProgress.componentData.destroy()
    }
}

export default ElvenProgress
// Elven Progressbar (Vue 3)
// https://github.com/oklookat

import ElvenProgressC from './ElvenProgressC'

class ElvenProgress {
    static componentData = null

    static install(app, options) {
        app.component('elven-progress', ElvenProgressC)
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const instance = this
        app.mixin({
            created() {
                if (this.SERVICE === 'ELVEN_PROGRESS_C') {
                    instance.componentData = this
                    theLogic.options = options
                    theLogic.init()
                    app.config.globalProperties.$elvenProgress = theLogic
                }
            },
        })
    }
}

class theLogic {
    static options = null

    static init() {
        if (this.options) {
            if (this.options.progressBarHeight) {
                ElvenProgress.componentData.progressBarHeight = this.options.progressBarHeight
            }
            if (this.options.basicLoadingStartSpeed) {
                ElvenProgress.componentData.basicLoadingStartSpeed = this.options.basicLoadingStartSpeed
            }
            if (this.options.basicLoadingFinishSpeed) {
                ElvenProgress.componentData.basicLoadingFinishSpeed = this.options.basicLoadingFinishSpeed
            }
            if (this.options.basicLoadingStartTo) {
                ElvenProgress.componentData.basicLoadingStartTo = this.options.basicLoadingStartTo
            }
        }
    }

    static setPercents(percents) {
        ElvenProgress.componentData.setPercents = percents
    }

    static loadingStart() {
        ElvenProgress.componentData.basicLoading = true
    }

    static loadingFinish() {
        ElvenProgress.componentData.basicLoading = false
    }

    static close() {
        ElvenProgress.componentData.closeBar = true
    }
}

export default ElvenProgress
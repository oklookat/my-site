class ElvenProgressLogic {

    private options = null
    public events = new EventManager()

    // static init() {
    //     if (this.options) {
    //         if (this.options.height) {
    //             ElvenProgress.componentData.height = this.options.height
    //         }
    //         if (this.options.loadingStartSpeed) {
    //             ElvenProgress.componentData.loadingStartSpeed = this.options.loadingStartSpeed
    //         }
    //         if (this.options.loadingFinishSpeed) {
    //             ElvenProgress.componentData.loadingFinishSpeed = this.options.loadingFinishSpeed
    //         }
    //         if (this.options.loadingStartTo) {
    //             ElvenProgress.componentData.loadingStartTo = this.options.loadingStartTo
    //         }
    //     }
    // }


    public loadingStart() {
        this.events.fire('loading', true)
    }

    public loadingFinish() {
        this.events.fire('loading', false)
    }

    public close() {
        this.events.fire('close', null)
    }
}

export default ElvenProgressLogic
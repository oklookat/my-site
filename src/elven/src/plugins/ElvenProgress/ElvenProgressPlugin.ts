import Progress from "./Progress.svelte"


export default class ElvenProgressPlugin {

    private progress

    constructor(element: HTMLElement) {
        this.progress = new Progress({
            target: element
        })
        window.$elvenProgress = this
    }

    // start loading
    public start() {
        this.progress.$set({ isLoading: true })
    }

    // finish loading
    public finish() {
        this.progress.$set({ isLoading: false })
    }

    // destroy plugin and element
    public destroy() {
        this.progress.$destroy()
    }
}
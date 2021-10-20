import ElvenProgressC from "./ElvenProgressC.svelte"


export default class ElvenProgress {

    private progress

    constructor() {
        const el = document.querySelector('#elven__progress')
        if (!el) {
            throw Error('elvenProgress: element not found')
        }
        this.progress = new ElvenProgressC({
            target: el
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
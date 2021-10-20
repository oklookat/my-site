import ElvenNotifyC from './ElvenNotifyC.svelte'

export default class ElvenNotify {

    private notify: ElvenNotifyC

    constructor() {
        const el = document.querySelector('#elven__notify')
        if (!el) {
            throw Error('elvenNotify: element not found')
        }
        this.notify = new ElvenNotifyC({
            target: el
        })
        window.$elvenNotify = this
    }

    public add(message: string) {
        this.notify.$set({ addNot: message })
    }

    // destroy plugin and element
    public destroy() {
        this.notify.$destroy()
    }
}
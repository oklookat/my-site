import Notify from './Notify.svelte'

export default class ElvenNotify {

    private notify: Notify

    constructor(element: HTMLElement) {
        this.notify = new Notify({
            target: element
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
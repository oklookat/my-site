/**
 * resize textarea depending on text size or window resize
 */
export default class TextareaResizer {

    private textarea: HTMLElement
    private _onInput = this.onInput.bind(this)
    private _onWindowResize = this.onWindowResize.bind(this)
    private resizeTimeout: NodeJS.Timeout

    constructor(textarea: HTMLTextAreaElement) {
        if (!(textarea instanceof HTMLTextAreaElement)) {
            throw Error(`TextareaResizer: wrong element provided. Provide textarea element`)
        }
        this.textarea = textarea
        this.manageEvents(true)
        this.onInput()
    }

    public destroy() {
        this.manageEvents(false)
    }

    /** add or remove events */
    private manageEvents(add: boolean) {
        const action = add ? 'addEventListener' : 'removeEventListener'
        this.textarea[action]('input', this._onInput)
        this.textarea[action]('change', this._onInput)
        window[action]('resize', this._onWindowResize)
    }

    private onInput() {
        this.textarea.style.height = `0`
        const height = this.textarea.clientHeight
        const scrollHeight = this.textarea.scrollHeight
        const correct = scrollHeight === height
        if (correct) {
            return
        }
        this.textarea.style.height = `${scrollHeight}px`
    }

    private onWindowResize() {
        if (this.resizeTimeout) {
            clearTimeout(this.resizeTimeout)
        }
        this.resizeTimeout = setTimeout(() => {
            this.onInput()
        }, 500)
    }

}
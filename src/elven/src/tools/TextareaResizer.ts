/**
 * resize textarea depending on text size or window resize
 */
export default class TextareaResizer {

    private element: HTMLElement
    private _onInput = this.onInput.bind(this)
    private _onWindowResize = this.onWindowResize.bind(this)
    private resizeTimeout: NodeJS.Timeout

    constructor(element: HTMLElement) {
        if (!element.style) {
            throw Error(`TextareaResizer: wrong element`)
        }
        this.element = element
        this.start()
    }

    private start() {
        this.onInput()
        this.element.addEventListener('input', this._onInput)
        window.addEventListener('resize', this._onWindowResize)
    }

    public destroy() {
        this.element.removeEventListener('input', this._onInput)
        window.removeEventListener('resize', this._onWindowResize)
    }

    private onInput() {
        this.element.style.height = `0px`
        const height = this.element.clientHeight
        const scrollHeight = this.element.scrollHeight
        const isEqual = scrollHeight === height
        if (!isEqual) {
            this.element.style.height = `${scrollHeight}px`
        }
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
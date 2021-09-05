import {reactive} from "vue";

interface ISlider {
    element: HTMLElement | null
    isMouseDown: boolean
    percents: number
}

export default class Composition {

    public slider: ISlider
    private readonly _onMovingMouseDown = this.onMovingMouseDown.bind(this)
    private readonly _onDocumentMouseMove = this.onDocumentMouseMove.bind(this)
    private readonly _onDocumentMouseUp = this.onDocumentMouseUp.bind(this)

    constructor() {
        this.slider = {
            element: null,
            isMouseDown: false,
            percents: 100,
        }
        this.slider = reactive(this.slider)
    }

    public init(element: HTMLElement) {
        this.slider.element = element
        this.slider.element.addEventListener('mousedown', this._onMovingMouseDown, {passive: false})
        this.slider.element.addEventListener('touchstart', this._onMovingMouseDown, {passive: false})
    }

    public destroy() {
        this.slider.element.removeEventListener('mousedown', this._onMovingMouseDown)
        this.slider.element.removeEventListener('touchstart', this._onMovingMouseDown)
    }

    private onMovingMouseDown(event) {
        this.slider.isMouseDown = true
        // pre-compute moving, because user already clicked
        this.onDocumentMouseMove(event)
        // setup document events, not local, because more comfortable control all this stuff when you moving as you like
        document.addEventListener("mousemove", this._onDocumentMouseMove, {passive: false})
        document.addEventListener('mouseup', this._onDocumentMouseUp, {passive: false})
        document.addEventListener("touchmove", this._onDocumentMouseMove, {passive: false})
        document.addEventListener("touchend", this._onDocumentMouseUp, {passive: false})
        document.addEventListener("touchcancel", this._onDocumentMouseUp, {passive: false})
    }

    private movingCleanup() {
        document.removeEventListener("mousemove", this._onDocumentMouseMove)
        document.removeEventListener('mouseup', this._onDocumentMouseUp)
        document.removeEventListener("touchmove", this._onDocumentMouseMove)
        document.removeEventListener("touchend", this._onDocumentMouseUp)
        document.removeEventListener("touchcancel", this._onDocumentMouseUp)
        this.slider.isMouseDown = false
    }

    private onDocumentMouseMove(event) {
        this.computeMoving(event)
    }

    private onDocumentMouseUp() {
        // cleanup and set stuff
        // when the mouse button or touch is up after pressing
        if (this.slider.isMouseDown) {
            this.slider.isMouseDown = false
            this.movingCleanup()
        }
    }

    private computeMoving(event) {
        event.preventDefault()
        const rect = this.slider.element.getBoundingClientRect()
        const clientX = this.getClientX(event) - rect.left
        if (!clientX) {
            return
        }
        if (this.slider.isMouseDown) {
            this.computeView(clientX)
        }
    }

    private computeView(clientX: number) {
        let clickPosition = this.getClickPosition(clientX, this.slider.element)
        if (clickPosition < 0) {
            clickPosition = 0
        } else if (clickPosition > 1) {
            clickPosition = 1
        }
        const _at = clickPosition * 100
        this.slider.percents = this.computePercents(_at, 100)
    }


    //////////// SERVICE
    // get clientX (horizontal mouse position) by touch or mouse event
    private getClientX(event, staticRect): number {
        let clientX = 0
        const isMovedByTouchscreen = event.type.includes('touch') && event.touches && event.touches.length > 0
        if (isMovedByTouchscreen) {
            for (const touch of event.touches) {
                if (touch.clientX) {
                    clientX = touch.clientX
                    break
                }
            }
        } else if (event.type.includes('mouse')) {
            // if moved by mouse
            clientX = event.clientX
        } else {
            // if moved by unknown - reset
            this.movingCleanup()
        }
        if (!clientX) {
            return 0
        }
        return clientX
    }

    // get click position by clientX
    private getClickPosition(clientX: number, element: HTMLElement): number {
        return (clientX - element.offsetLeft) / element.offsetWidth
    }

    // get percents of current param by setting the total param
    private computePercents(current: number, total: number): number {
        current = Math.round(current)
        let percents = (current / total) * 100
        if (percents >= 100) {
            percents = 100
        } else if (total < 1) {
            percents = 0
        }
        return this.round(percents, 4)
    }

    // round the number to the specific number of decimal places
    private round(value: number, precision: number): number {
        // https://stackoverflow.com/a/7343013/16762009
        const multiplier = Math.pow(10, precision || 0)
        return Math.round(value * multiplier) / multiplier
    }
}
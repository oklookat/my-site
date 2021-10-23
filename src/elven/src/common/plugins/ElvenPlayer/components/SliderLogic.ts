import Service from "@/common/plugins/ElvenPlayer/tools/Service"
import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store'

interface ISlider {
    element: any
    isMouseDown: Writable<boolean>
    percents: Writable<number>
}

export default class SliderLogic {

    public slider: ISlider
    private readonly _onMovingMouseDown = this.onMovingMouseDown.bind(this)
    private readonly _onDocumentMouseMove = this.onDocumentMouseMove.bind(this)
    private readonly _onDocumentMouseUp = this.onDocumentMouseUp.bind(this)

    constructor() {
        this.slider = {
            element: null,
            isMouseDown: writable(false),
            percents: writable(100),
        }
    }

    // init not in constructor, because composition api specific thing
    public init(element: HTMLElement) {
        this.slider.element = element
        this.slider.element.addEventListener('mousedown', this._onMovingMouseDown, { passive: false })
        this.slider.element.addEventListener('touchstart', this._onMovingMouseDown, { passive: false })
    }

    public destroy() {
        this.slider.element.removeEventListener('mousedown', this._onMovingMouseDown)
        this.slider.element.removeEventListener('touchstart', this._onMovingMouseDown)
    }

    // when user click / pressed / touch
    private onMovingMouseDown(event) {
        this.slider.isMouseDown.set(true)
        // pre-compute moving, because user already clicked
        this.onDocumentMouseMove(event)
        // setup document events, not local, because more comfortable control all this stuff when you moving as you like
        document.addEventListener("mousemove", this._onDocumentMouseMove, { passive: false })
        document.addEventListener('mouseup', this._onDocumentMouseUp, { passive: false })
        document.addEventListener("touchmove", this._onDocumentMouseMove, { passive: false })
        document.addEventListener("touchend", this._onDocumentMouseUp, { passive: false })
        document.addEventListener("touchcancel", this._onDocumentMouseUp, { passive: false })
    }

    // remove moving events
    private movingCleanup() {
        this.slider.isMouseDown.set(false)
        document.removeEventListener("mousemove", this._onDocumentMouseMove)
        document.removeEventListener('mouseup', this._onDocumentMouseUp)
        document.removeEventListener("touchmove", this._onDocumentMouseMove)
        document.removeEventListener("touchend", this._onDocumentMouseUp)
        document.removeEventListener("touchcancel", this._onDocumentMouseUp)
    }

    // all move events be here
    private onDocumentMouseMove(event) {
        this.computeMoving(event)
    }

    // when user unpressed mouse button or finger
    private onDocumentMouseUp() {
        // cleanup and set stuff
        // when the mouse button or touch is up after pressing
        this.slider.isMouseDown.set(false)
        this.movingCleanup()
    }

    // get click position
    private computeMoving(event) {
        event.preventDefault()
        const rect = this.slider.element.getBoundingClientRect()
        let pageX = Service.getPageX(event)
        if (pageX === null) {
            // reset if no pageX
            this.onDocumentMouseUp()
            return
        }
        pageX = pageX - rect.left
        if (this.slider.isMouseDown) {
            this.computeView(pageX)
        }
    }

    // final get click position relatively of slider element in percents
    private computeView(pageX: number) {
        let clickPosition = Service.getClickPosition(pageX, this.slider.element)
        if (clickPosition < 0) {
            clickPosition = 0
        } else if (clickPosition > 1) {
            clickPosition = 1
        }
        const _at = clickPosition * 100
        this.slider.percents.set(Service.computePercents(_at, 100))
    }

}
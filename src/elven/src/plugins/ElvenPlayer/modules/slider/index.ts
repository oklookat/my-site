import Utils from "./utils"
import Store from "./store"
import type { Store as IStore } from "./types"

/** slider logic */
export default class Core {

    public store: IStore
    public container: HTMLDivElement
    /** @see https://developer.mozilla.org/en-US/docs/Web/API/PointerEvent#browser_compatibility */
    private mode: 'new' | 'legacy'
    // legacy
    private _containerMouseDown = null
    private _documentMouseMove = null
    private _documentMouseUp = null

    constructor(container: HTMLDivElement) {
        this.store = new Store()
        this.container = container
        if (!window.PointerEvent) {
            this.mode = 'new'
            this.startNew(true)
        } else {
            this.mode = 'legacy'
            this.startLegacy(true)
        }
    }

    /** new slider with pointerCapture in one method */
    private startNew(active: boolean) {
        // click
        const beginSliding = (e: PointerEvent) => {
            // disallow dragging slider with any mouse buttons except LMB
            if (e.pointerType === 'mouse' && e.button !== 0) {
                return
            }
            this.store.isMouseDown = true
            this.container.onpointermove = slide
            this.container.setPointerCapture(e.pointerId)
            // start slide because user already clicked
            slide(e)
        }
        // unclick
        const stopSliding = (e: PointerEvent) => {
            this.store.isMouseDown = false
            this.container.onpointermove = null
            this.container.releasePointerCapture(e.pointerId)
        }
        // sliding now
        const slide = (e: PointerEvent) => {
            const containerWidth = this.container.clientWidth
            const rect = this.container.getBoundingClientRect()
            const position = e.clientX - rect.left
            const percents = Utils.computePercents(position, containerWidth)
            this.store.percents = percents
        }
        if (active) {
            this.container.onpointerdown = beginSliding
            this.container.onpointerup = stopSliding
            this.container.onpointercancel = stopSliding
            this.container.oncontextmenu = (e) => {e.preventDefault()}
            return
        }
        this.store.isMouseDown = false
        this.container.onpointerdown = null
        this.container.onpointerup = null
        this.container.onpointermove = null
        this.container.onpointercancel = null
    }

    /**
    * setup events on container
    * @param add true = add events; false = remove events
    */
    private startLegacy(add: boolean) {
        this._containerMouseDown = add ? this.containerMouseDown.bind(this) : null
        this._documentMouseMove = add ? this.computeMoving.bind(this) : null
        this._documentMouseUp = add ? this.documentEvents.bind(this, false) : null
        const action = add ? 'addEventListener' : 'removeEventListener'
        const notPassive = { passive: false }
        this.container[action]('mousedown', this._containerMouseDown, notPassive)
        this.container[action]('touchstart', this._containerMouseDown, notPassive)
    }


    /**
     * setup events on document
     * @param add true = add events; false = remove events
     */
    private documentEvents(add: boolean) {
        if (!add) { this.store.isMouseDown = false }
        const notPassive = { passive: false }
        const action = add ? 'addEventListener' : 'removeEventListener'
        document[action]("mousemove", this._documentMouseMove, notPassive)
        document[action]('mouseup', this._documentMouseUp, notPassive)
        document[action]("touchmove", this._documentMouseMove, notPassive)
        document[action]("touchend", this._documentMouseUp, notPassive)
        document[action]("touchcancel", this._documentMouseUp, notPassive)
    }

    /** cleanup */
    public destroy() {
        if (this.mode === 'new') {
            this.startNew(false)
            return
        }
        this.startLegacy(false)
        this.documentEvents(false)
    }

    /** when user clicked / pressed / touched */
    private containerMouseDown(e: MouseEvent | TouchEvent) {
        // disallow dragging slider with any mouse buttons except LMB
        if (e instanceof MouseEvent && e.button !== 0) {
            return
        }
        this.store.isMouseDown = true
        // pre-compute moving, because user already clicked
        this.computeMoving(e)
        // setup document events. We no add events to container 
        // because more comfortable control slider by dragging it all over the window
        this.documentEvents(true)
    }


    /** get click position */
    private computeMoving(e: MouseEvent | TouchEvent) {
        e.preventDefault()
        const rect = this.container.getBoundingClientRect()
        let pageX = Utils.getPageX(e)
        pageX = pageX - rect.left
        if (this.store.isMouseDown) {
            // final get click position relatively of slider element in percents
            this.store.percents = Utils.getClickPercentsWidth(pageX, this.container)
        }
    }

}
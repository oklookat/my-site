import Utils from "./utils"
import Store from "./store"
import type { IStore } from "./types"

export default class Core {

    public store: IStore
    public container: HTMLDivElement
    private readonly _containerMouseDown = this.containerMouseDown.bind(this)
    private readonly _documentMouseMove = this.computeMoving.bind(this)
    private readonly _documentMouseUp = this.documentEvents.bind(this, false)

    constructor(container: HTMLDivElement) {
        this.store = new Store()
        this.container = container
        this.containerEvents(true)
    }

    /**
     * setup events on container
     * @param add true = add events; false = remove events
     */
    private containerEvents(add: boolean) {
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
        this.containerEvents(false)
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
            this.computeView(pageX)
        }
    }

    /** final get click position relatively of slider element in percents */
    private computeView(pageX: number) {
        let clickPosition = Utils.getClickPercentsWidth(pageX, this.container)
        if (clickPosition < 0) {
            clickPosition = 0
        } else if (clickPosition > 1) {
            clickPosition = 1
        }
        const percents = clickPosition * 100
        this.store.percents = percents
    }

}
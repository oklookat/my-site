import type { Events } from "../types"


/** controls audio element and adds helpful things */
export default class DOM {

    private element: HTMLAudioElement
    private events: Events
    private _onPlaying: (e?: Event) => void
    private _onPause: (e?: Event) => void
    private _onEnded: (e?: Event) => void
    private _onTimeUpdate: (e?: Event) => void
    private _onError: (e?: Event) => void


    constructor(events: Events, source?: string, volume?: number) {
        this.element = new Audio(source)
        this.volume = volume
        this.events = events
        this._onPlaying = this.events.onPlaying.bind(this.events)
        this._onPause = this.events.onPause.bind(this.events)
        this._onEnded = this.events.onEnded.bind(this.events)
        this._onTimeUpdate = this.events.onTimeUpdate.bind(this.events)
        this._onError = this.events.onError.bind(this.events)
        this.manageEvents(true)
    }

    /** add (true) / remove (false) events */
    private manageEvents(add: boolean) {
        const action = add ? 'addEventListener' : 'removeEventListener'
        this.element[action]('playing', this._onPlaying)
        this.element[action]('pause', this._onPause)
        this.element[action]('ended', this._onEnded)
        this.element[action]('timeupdate', this._onTimeUpdate)
        this.element[action]('error', this._onError)
    }

    public destroy() {
        const badElement = !this.element || !this.events
        if (badElement) {
            return
        }
        this.manageEvents(false)
        this.stop()
        this.element = null
    }

    public async play(): Promise<void> {
        try {
            await this.element.play()
            return Promise.resolve()
        } catch (err) {
        }
    }

    public pause() {
        this.element.pause()
    }

    public stop() {
        this.pause()
        this.currentTime = 0
    }

    public isCanPlay(): boolean {
        const state = this.element.readyState
        const err = this.element.error
        const notCanPlay = state === 0 || err != null
        return !notCanPlay
    }

    public get source(): string {
        return this.element.currentSrc
    }

    public set source(src: string) {
        this.element.src = src
    }

    public get currentTime(): number {
        return this.element.currentTime
    }

    public set currentTime(seconds: number) {
        if (!this.isCanPlay()) {
            return
        }
        const duration = this.duration
        // set max or min time if v bigger or lower than duration
        if (seconds > duration) {
            seconds = duration
        } else if (seconds < 0) {
            seconds = 0
        }
        this.element.currentTime = seconds
    }

    public get volume(): number {
        return this.element.volume
    }

    public set volume(vol: number) {
        // correct volume if volume > 1; < 0 or falsy but not 0
        if (vol > 1 || !vol && vol !== 0) {
            vol = 1.0
        } else if (vol < 0) {
            vol = 0
        }
        this.element.volume = vol
    }

    public get duration(): number {
        return this.element.duration
    }

    public get buffered(): TimeRanges {
        return this.element.buffered
    }

}
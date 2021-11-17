import type { IEvents } from "../types"

/** controls one audio element */
export default class DOM {

    private element: HTMLAudioElement
    private _source: string
    private events: IEvents

    constructor(source: string, events: IEvents, volume?: number) {
        this.source = source
        this.element = new Audio(this.source)
        this.element.volume = volume ? volume : 1.0
        this.events = events
        this.manageEvents(true)
    }

    /** destroy audio element */
    public destroy() {
        if (!this.element || !this.events) {
            return
        }
        this.manageEvents(false)
        this.element = null
    }

    /** add (true) / remove (false) events */
    private manageEvents(add: boolean) {
        const action = add ? 'addEventListener' : 'removeEventListener'
        this.element[action]('playing', this.events.onPlaying.bind(this.events))
        this.element[action]('pause', this.events.onPause.bind(this.events))
        this.element[action]('ended', this.events.onEnded.bind(this.events))
        this.element[action]('timeupdate', this.events.onTimeUpdate.bind(this.events))
        this.element[action]('error', this.events.onError.bind(this.events))
    }

    /** play audio */
    public async play(): Promise<void> {
        try {
            await this.element.play()
            return Promise.resolve()
        } catch (err) {
            return Promise.reject(err)
        }
    }

    /** pause audio */
    public pause() {
        this.element.pause()
    }

    /** stop audio */
    public stop() {
        this.element.pause()
        this.element.currentTime = 0
    }

    public get source(): string {
        return this._source
    }

    public set source(s: string) {
        if (this.element) {
            this.element.src = s
        }
        this._source = s
    }

    /** get current position of audio in seconds */
    public get currentTime(): number {
        return this.element.currentTime
    }

    /** set current position of audio in seconds */
    public set currentTime(v: number) {
        this.element.currentTime = v
    }

    /** set current time by percents. Where 100 - audio ends. */
    public set currentTimePercents(percents: number) {
        const total = this.element.duration
        const percToTime = Math.floor((total / 100) * percents)
        // set maximum time when percents > total time, otherwise set converted time
        this.element.currentTime = percToTime > total ? total : percToTime
    }

    /** get volume in float */
    public get volume(): number {
        return this.element.volume
    }

    /** set volume in float */
    public set volume(v: number) {
        this.element.volume = v
    }

    /** set volume in percents */
    public set volumePercents(percents: number) {
        let vol = (percents / 100)
        // if volume > 1 or < 0 set min or max volume
        vol = vol > 1.0 ? 1.0 : vol < 0 ? 0 : vol
        this.element.volume = vol
    }

    /** get volume in percents */
    public get volumePercents(): number {
        let vol = this.element.volume * 100
        vol = vol > 100 ? 100 : vol < 0 ? 0 : vol
        return vol
    }

    /** get total time of audio in seconds */
    public get duration(): number {
        return this.element.duration
    }

    /** get audio buffered TimeRanges */
    public get buffered(): TimeRanges {
        return this.element.buffered
    }

}
import Utils from "./utils"
import type { IEvents } from "../types"


/** controls audio element and adds helpful things */
export default class DOM {

    private element: HTMLAudioElement
    private events: IEvents

    constructor(events: IEvents, source?: string, volume?: number) {
        this.element = new Audio(source)
        this.volume = volume
        this.events = events
        this.manageEvents(true)
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

    /** destroy audio element and events */
    public destroy() {
        const badElement = !this.element || !this.events
        if (badElement) {
            return
        }
        this.manageEvents(false)
        this.stop()
        this.element = null
    }

    /** play audio */
    public async play(): Promise<void> {
        try {
            await this.element.play()
            return Promise.resolve()
        } catch (err) {
        }
    }

    /** pause audio */
    public pause() {
        this.element.pause()
    }

    /** stop audio */
    public stop() {
        this.pause()
        this.currentTime = 0
    }

    /** is audio can played */
    public isCanPlay(): boolean {
        const state = this.element.readyState
        const err = this.element.error
        const notCanPlay = state === 0 || err != null
        return !notCanPlay
    }

    /** convert current time percents to seconds depending on duration */
    public convertPercentsToCurrentTime(perc: number): number {
        const dur = this.duration
        return Utils.convertPercentsToCurrentTime(perc, dur)
    }

    /** convert current time percents to string like '01:11' or '01:11:11' */
    public convertPercentsToCurrentTimePretty(perc: number): string {
        const dur = this.duration
        const time = this.convertPercentsToCurrentTime(perc)
        return Utils.convertCurrentTimePretty(time, dur)
    }

    /** get source */
    public get source(): string {
        return this.element.currentSrc
    }

    /** set source */
    public set source(src: string) {
        this.element.src = src
    }

    /** get current time of audio in seconds */
    public get currentTime(): number {
        return this.element.currentTime
    }

    /** set current time of audio in seconds */
    public set currentTime(v: number) {
        if (!this.isCanPlay()) {
            return
        }
        const dur = this.duration
        // set max or min time if v bigger or lower than duration
        if (v > dur) {
            v = dur
        } else if (v < 0) {
            v = 0
        }
        this.element.currentTime = v
    }

    /** get current time percents. Where 100 - audio ends */
    public get currentTimePercents(): number {
        return Utils.getPercents(this.currentTime, this.duration)
    }

    /** set current time by percents. Where 100 - audio ends */
    public set currentTimePercents(perc: number) {
        const dur = this.duration
        perc = Utils.convertPercentsToCurrentTime(perc, dur)
        this.currentTime = perc
    }

    /** get volume (0 - 1) */
    public get volume(): number {
        return this.element.volume
    }

    /** set volume (0 - 1) */
    public set volume(vol: number) {
        // correct volume if volume > 1; < 0 or falsy but not 0
        if (vol > 1 || !vol && vol !== 0) {
            vol = 1.0
        } else if (vol < 0) {
            vol = 0
        }
        this.element.volume = vol
    }

    /** set volume in percents (0 - 100) */
    public set volumePercents(perc: number) {
        this.volume = perc / 100
    }

    /** get volume in percents (0 - 100) */
    public get volumePercents(): number {
        return this.volume * 100
    }

    /** get audio duration in seconds */
    public get duration(): number {
        return this.element.duration
    }

    /** get audio buffered */
    public get buffered(): TimeRanges {
        return this.element.buffered
    }

}
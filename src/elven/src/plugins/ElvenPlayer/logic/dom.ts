import Utils from "./utils"
import type { IEvents } from "../types"


/** controls audio element and adds helpful things */
export default class DOM {

    private element: HTMLAudioElement
    private events: IEvents

    constructor(events: IEvents, source?: string, volume?: number) {
        this.element = new Audio(source)
        this.volume = volume ? volume : 1.0
        this.events = events
        this.manageEvents(true)
    }

    /** get source */
    public get source(): string {
        return this.element.currentSrc
    }

    /** set source */
    public set source(src: string) {
        this.element.src = src
    }

    /** get current position of audio in seconds */
    public get currentTime(): number {
        return this.element.currentTime
    }

    /** set current position of audio in seconds */
    public set currentTime(v: number) {
        const dur = this.duration
        // set max or min time if v bigger or lower then duration
        v = v > dur ? dur : v < 0 ? 0 : v
        this.element.currentTime = v
    }

    /** get current time percents. Where 100 - audio ends. */
    public get currentTimePercents(): number {
        return Utils.getPercents(this.currentTime, this.duration)
    }

    /** set current time by percents. Where 100 - audio ends. */
    public set currentTimePercents(perc: number) {
        const dur = this.duration
        perc = Utils.percentsToCurrentTime(perc, dur)
        this.currentTime = perc
    }

    /** get volume (0 - 1) */
    public get volume(): number {
        let vol = this.element.volume
        vol = vol > 1 ? 1.0 : vol < 0 ? 0 : vol
        return this.element.volume
    }

    /** set volume (0 - 1) */
    public set volume(vol: number) {
        // if volume > 1 or < 0 set min or max volume
        vol = vol > 1 ? 1.0 : vol < 0 ? 0 : vol
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

    /** convert percents position to seconds */
    public getCurrentTimeByPercents(perc: number): number {
        const dur = this.duration
        return Utils.percentsToCurrentTime(perc, dur)
    }

    /** convert percents position to string like '01:11' */
    public getCurrentTimePrettyByPercents(perc: number): string {
        const dur = this.duration
        const time = this.getCurrentTimeByPercents(perc)
        return Utils.getPositionPretty(time, dur)
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
        if (!this.element || !this.events) {
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
        this.element.pause()
        this.currentTime = 0
    }
}
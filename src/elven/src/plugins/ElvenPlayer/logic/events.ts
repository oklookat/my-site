import type { IEvents, IState } from "../types";
import Utils from "./utils";
import Logger from "./logger";


/** updates playback state */
export default class Events implements IEvents {

    private state: IState

    constructor(state: IState) {
        this.state = state
    }

    public onPlaying() {
        this.state.playing = true
        this.state.ended = false
    }

    public onPause() {
        this.state.playing = false
        this.state.ended = false
    }

    public onEnded() {
        this.state.playing = false
        this.state.ended = true
    }

    public onTimeUpdate(e: Event) {
        const el = e.target as HTMLAudioElement
        const currentTime = el.currentTime
        const buffered = el.buffered
        let duration = el.duration
        const badDuration = !duration || isNaN(duration) || duration === Infinity
        duration = badDuration ? 0 : duration
        this.state.bufferedPercents = Utils.getBufferedPercents(currentTime, duration, buffered)
        this.state.durationNum = duration
        this.state.durationPretty = Utils.getPretty(duration, 'auto')
        this.state.currentTimeNum = currentTime
        this.state.currentTimePretty = Utils.convertCurrentTimePretty(currentTime, duration)
        this.state.currentTimePercents = Utils.getPercents(currentTime, duration)
    }

    public onError(e: Event) {
        // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/networkState
        const valid = e.target && e.target instanceof HTMLMediaElement
        if (!valid) {
            Logger.error('unknown error, invalid event target')
            return
        }
        const target = e.target as HTMLMediaElement
        const err = target.error
        const msg = err.message ? ` ${err.message}` : ''
        switch (err.code) {
            case MediaError.MEDIA_ERR_ABORTED:
                Logger.error(`aborted.${msg}`)
                break
            case MediaError.MEDIA_ERR_NETWORK:
                Logger.error(`network error.${msg}`)
                break
            case MediaError.MEDIA_ERR_DECODE:
                Logger.error(`decode error. Maybe audio damaged or something?${msg}`)
                break
            case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
                Logger.error(`not supported.${msg}`)
                break
            default:
                Logger.error(`unknown error.${msg}`)
                break
        }
    }

}
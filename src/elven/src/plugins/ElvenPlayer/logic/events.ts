import type { IEvents, IStore } from "../types";
import Utils from "./utils";
import Logger from "./logger";


/** updates playback state */
export default class Events implements IEvents {

    private store: IStore

    constructor(store: IStore) {
        this.store = store
    }

    public onPlaying() {
        this.store.playing = true
        this.store.ended = false
    }

    public onPause() {
        this.store.playing = false
        this.store.ended = false
    }

    public onEnded() {
        this.store.playing = false
        this.store.ended = true
    }

    public onTimeUpdate(e: Event) {
        const el = e.target as HTMLAudioElement
        const currentTime = el.currentTime
        const buffered = el.buffered
        let duration = el.duration
        const badDuration = !duration || isNaN(duration) || duration === Infinity
        duration = badDuration ? 0 : duration
        this.store.bufferedPercents = Utils.getBufferedPercents(currentTime, duration, buffered)
        this.store.durationNum = duration
        this.store.durationPretty = Utils.getPretty(duration, 'auto')
        this.store.currentTimeNum = currentTime
        this.store.currentTimePretty = Utils.convertCurrentTimePretty(currentTime, duration)
        this.store.currentTimePercents = Utils.getPercents(currentTime, duration)
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
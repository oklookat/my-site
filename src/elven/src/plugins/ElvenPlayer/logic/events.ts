import Utils from "./utils";
import type { IEvents, IState, TConvertSecondsMode } from "../types";


export default class Events implements IEvents {
    private state: IState
    private next: () => void

    constructor(state: IState, next: () => void) {
        this.state = state
        this.next = next
    }

    public onPlaying() {
        this.state.isPlaying = true
    }
    public onPause() {
        this.state.isPlaying = false
    }
    public onEnded() {
        this.state.isPlaying = false
        this.next()
    }

    public onTimeUpdate(e: Event) {
        const el = e.target as HTMLAudioElement
        const currentTime = el.currentTime
        const buffered = el.buffered
        let duration = el.duration
        const badDuration = !duration || isNaN(duration) || duration === Infinity
        duration = badDuration ? 0 : duration
        this.state.bufferedPercents = Utils.computeBuffered(currentTime, duration, buffered)
        this.state.durationNum = duration
        this.state.durationPretty = Utils.convertSeconds(duration, 'auto')
        this.state.positionNum = currentTime
        const mode: TConvertSecondsMode = duration < 3600 ? 'minutes' : 'hours'
        this.state.positionPretty = Utils.convertSeconds(currentTime, mode)
        this.state.positionPercents = Utils.computePercents(currentTime, duration)
    }

    public onError(event) {
        // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/networkState
        switch (event.target.error.code) {
            case event.target.error.MEDIA_ERR_ABORTED:
                console.error('elvenPlayer: aborted')
                break
            case event.target.error.MEDIA_ERR_NETWORK:
                console.error('elvenPlayer: network error')
                break
            case event.target.error.MEDIA_ERR_DECODE:
                console.error('elvenPlayer: decode error. Maybe audio damaged or something?')
                break
            case event.target.error.MEDIA_ERR_SRC_NOT_SUPPORTED:
                // if (this.player.element.src === 'https://null/') {
                //     return
                // }
                if (event.target.error && event.target.error.message) {
                    console.error(`elvenPlayer: not supported. Error: ${event.target.error.message}`)
                } else {
                    console.error(`elvenPlayer: not supported.`)
                }
                break
            default:
                console.error('elvenPlayer: unknown error')
                break
        }
    }
}
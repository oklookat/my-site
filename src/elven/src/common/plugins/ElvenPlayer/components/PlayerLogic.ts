import { writable } from 'svelte/store';
import type { Writable } from 'svelte/store'

type TConvertSecondsMode = 'auto' | 'hours' | 'minutes'
interface IPlayer {
    // is player component active
    active: boolean
    initialized: boolean
    element: HTMLAudioElement
    isPlaying: Writable<boolean>
    // volume in float (1.0)
    volume: number
    // volume in percents
    percentsVolume: Writable<number>
    currentPlaying: {
        // total time of current audio like '04:23'
        duration: Writable<string>,
        // current time of playing audio like '01:23'
        currentTime: Writable<string>,
        // percents of audio reached
        percentsReached: Writable<number>
        // percents of buffered audio
        percentsBuffered: Writable<number>
        playlist: string[]
        index: number
    }
    // moving
    isProgressMouseDown: boolean
    progressPreviewTime: number
    isVolumeMouseDown: boolean
}

export default class PlayerLogic {

    private eventsInitialized: boolean = false
    public player: IPlayer
    private playPromise: Promise<void> | undefined = undefined
    // saved event links (for cleanup)
    private _onPlaying = this.onPlaying.bind(this)
    private _onPause = this.onPause.bind(this)
    private _onEnded = this.onEnded.bind(this)
    private _onTimeUpdate = this.onTimeUpdate.bind(this)
    private _onError = this.onError.bind(this)
    private readonly nullSrc = 'https://null/'

    constructor() {
        this.player = {
            active: false,
            initialized: false,
            element: new Audio(''),
            isPlaying: writable(false),
            volume: 1.0,
            percentsVolume: writable(100),
            currentPlaying: {
                duration: writable('00:00'),
                currentTime: writable('00:00'),
                percentsReached: writable(0),
                percentsBuffered: writable(0),
                playlist: [],
                index: 0
            },
            // moving
            isProgressMouseDown: false,
            progressPreviewTime: 0,
            isVolumeMouseDown: false,
        }
        this.player.element.volume = this.player.volume
    }

    public init() {
        if (this.eventsInitialized) {
            return
        }
        this.player.element.addEventListener('playing', this._onPlaying)
        this.player.element.addEventListener('pause', this._onPause)
        this.player.element.addEventListener('ended', this._onEnded)
        this.player.element.addEventListener('timeupdate', this._onTimeUpdate)
        this.player.element.addEventListener('error', this._onError)
        this.eventsInitialized = true
    }

    public async destroy() {
        if (!this.eventsInitialized) {
            return
        }
        this.player.element.removeEventListener('playing', this._onPlaying)
        this.player.element.removeEventListener('pause', this._onPause)
        this.player.element.removeEventListener('ended', this._onEnded)
        this.player.element.removeEventListener('timeupdate', this._onTimeUpdate)
        this.player.element.removeEventListener('error', this._onError)
        this.eventsInitialized = false
        await this.stop()
    }


    // playback controls
    public async play() {
        if (!this.player.initialized) {
            try {
                this.setCurrentAudio()
            } catch (err) {
                console.error(err)
                return
            }
        }
        this.playPromise = this.player.element.play()
        if (!this.playPromise) {
            return
        }
        this.playPromise
            .then(() => {
                this.player.initialized = true
            })
            .catch((err) => { console.error(err); this.stop() })
    }

    public pause() {
        if (!this.playPromise) {
            return
        }
        this.player.element.pause()
    }

    public async next() {
        if (!this.isHasNextAudio()) {
            this.player.element.currentTime = 0
            return
        }
        this.player.currentPlaying.index++
        this.setCurrentAudio()
        await this.play()
    }

    public async prev() {
        if (!this.isHasPrevAudio() || this.isAudioNotInStart()) {
            this.player.element.currentTime = 0
            return
        }
        this.player.currentPlaying.index--
        this.setCurrentAudio()
        await this.play()
    }

    public async stop() {
        this.player.element.src = this.nullSrc
        this.player.currentPlaying.index = 0
        this.player.currentPlaying.duration.set('00:00')
        this.player.currentPlaying.currentTime.set('00:00')
        this.player.currentPlaying.percentsReached.set(0)
        this.player.currentPlaying.percentsBuffered.set(0)
        this.player.isProgressMouseDown = false
        this.player.progressPreviewTime = 0
        this.player.isVolumeMouseDown = false
        this.player.initialized = false
        this.player.isPlaying.set(false)
        this.playPromise = undefined
    }

    // playback management
    public setCurrentAudio(playlistIndex = this.player.currentPlaying.index) {
        if (this.player.currentPlaying.playlist.length < 1) {
            return Error('elvenPlayer: empty playlist')
        }
        this.player.element.src = this.player.currentPlaying.playlist[playlistIndex]
    }

    public async setPlaylist(playlist: string[]) {
        await this.stop()
        this.player.currentPlaying.index = 0
        this.player.currentPlaying.playlist = playlist
    }

    public addToPlaylist(url: string) {
        this.player.currentPlaying.playlist.push(url)
    }

    // is audio not in start. Start = audio duration / 4.
    private isAudioNotInStart(): boolean {
        const isNotInStart = this.player.element.duration / 4
        return this.player.element.currentTime > isNotInStart
    }

    private isHasNextAudio(): boolean {
        const isHas = this.player.currentPlaying.playlist[this.player.currentPlaying.index + 1] !== undefined
        return isHas
    }

    private isHasPrevAudio(): boolean {
        const isHas = this.player.currentPlaying.playlist[this.player.currentPlaying.index - 1] !== undefined
        return isHas
    }

    public setTimeByPercents(percents: number) {
        const duration = this.player.element.duration
        const timeConverted = Math.floor((duration / 100) * percents)
        // set maximum time when percents > total time
        if (timeConverted > duration) {
            this.player.element.currentTime = duration
            return
        }
        this.player.element.currentTime = timeConverted
    }

    public setVolumeByPercents(percents: number) {
        let percentsVolume = (percents / 100)
        if (percentsVolume > 1.0) {
            percentsVolume = 1.0
        }
        if (percentsVolume < 0) {
            percentsVolume = 0
        }
        this.player.percentsVolume.set(percents)
        this.player.element.volume = percentsVolume
    }

    // events
    private onPlaying() {
        this.player.isPlaying.set(true)
    }

    private onPause() {
        this.player.isPlaying.set(false)
    }

    private async onEnded() {
        this.player.isPlaying.set(false)
        await this.next()
    }

    private async onError(event) {
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
                if (this.player.element.src === 'https://null/') {
                    return
                }
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

    private onTimeUpdate() {
        if (!this.player.isProgressMouseDown) {
            this.player.currentPlaying.percentsReached.set(PlayerLogic.computePercents(this.player.element.currentTime, this.player.element.duration))
            this.player.currentPlaying.percentsBuffered.set(PlayerLogic.computeBuffered(this.player.element))
        }
        if (this.player.element.duration) {
            this.player.currentPlaying.duration.set(PlayerLogic.convertSeconds(this.player.element.duration, 'auto'))
            let mode: TConvertSecondsMode
            if (this.player.element.duration < 3600) {
                mode = 'minutes'
            } else {
                mode = 'hours'
            }
            this.player.currentPlaying.currentTime.set(PlayerLogic.convertSeconds(this.player.element.currentTime, mode))
        } else {
            this.player.currentPlaying.duration.set('00:00')
            this.player.currentPlaying.currentTime.set('00:00')
        }
    }

    private static computePercents(current: number, total: number): number {
        let percents = (current / total) * 100
        percents = Math.round(percents)
        if (percents >= 100) {
            percents = 100
        } else if (total < 1) {
            percents = 0
        }
        return percents
    }

    // convert seconds to string like '01:23'
    private static convertSeconds(seconds: number, mode: TConvertSecondsMode): string {
        // https://stackoverflow.com/a/1322771/16762009
        switch (mode) {
            case 'auto':
                if (seconds < 3600) {
                    // like 00:01
                    return returnMinutes(seconds)
                } else {
                    // like 01:23:12
                    return returnHours(seconds)
                }
            case 'hours':
                // like 01:23:12
                return returnHours(seconds)
            case 'minutes':
                // like 00:01
                return returnMinutes(seconds)
        }

        function returnHours(seconds: number) {
            return new Date(seconds * 1000).toISOString().substr(11, 8)
        }

        function returnMinutes(seconds: number) {
            return new Date(seconds * 1000).toISOString().substr(14, 5)
        }
    }

    private static computeBuffered(playerEL: HTMLAudioElement): number {
        const currentTime = Math.round(playerEL.currentTime)
        const duration = playerEL.duration
        if (duration > 0) {
            for (let i = 0; i < playerEL.buffered.length; i++) {
                const len = playerEL.buffered.length - 1 - i
                if (playerEL.buffered.start(len) < currentTime) {
                    return Math.round(this.computePercents(playerEL.buffered.end(len), duration))
                }
            }
        }
        return 0
    }

}
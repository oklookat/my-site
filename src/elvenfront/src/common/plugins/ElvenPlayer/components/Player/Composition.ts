import { reactive } from "vue";
import Service, { TConvertSecondsMode } from "@/common/plugins/ElvenPlayer/tools/Service"

interface IAudioPlayer {
    // is player component active
    active: boolean
    initialized: boolean
    element: HTMLAudioElement
    isPlaying: boolean
    // volume in float (1.0)
    volume: number
    // volume in percents
    percentsVolume: number
    currentPlaying: {
        // total time of current audio like '04:23'
        duration: string,
        // current time of playing audio like '01:23'
        currentTime: string,
        // percents of audio reached
        percentsReached: number
        // percents of buffered audio
        percentsBuffered: number
        playlist: string[]
        index: number
    }
    // moving
    isProgressMouseDown: boolean
    progressPreviewTime: number
    isVolumeMouseDown: boolean
}

export default class Composition {

    public audioPlayer: IAudioPlayer
    private playPromise: Promise<void> | undefined = undefined
    // saved event links (for cleanup)
    private _onPlaying = this.onPlaying.bind(this)
    private _onPause = this.onPause.bind(this)
    private _onEnded = this.onEnded.bind(this)
    private _onTimeUpdate = this.onTimeUpdate.bind(this)
    private _onError = this.onError.bind(this)
    private readonly nullSrc = 'https://null/'

    constructor() {
        this.audioPlayer = {
            active: false,
            initialized: false,
            element: new Audio(''),
            isPlaying: false,
            volume: 1.0,
            percentsVolume: 100,
            currentPlaying: {
                duration: '00:00',
                currentTime: '00:00',
                percentsReached: 0,
                percentsBuffered: 0,
                playlist: [],
                index: 0
            },
            // moving
            isProgressMouseDown: false,
            progressPreviewTime: 0,
            isVolumeMouseDown: false,
        }
        this.audioPlayer.element.volume = this.audioPlayer.volume
        this.audioPlayer = reactive(this.audioPlayer)
    }

    public init() {
        // audio element events
        this.audioPlayer.element.addEventListener('playing', this._onPlaying)
        this.audioPlayer.element.addEventListener('pause', this._onPause)
        this.audioPlayer.element.addEventListener('ended', this._onEnded)
        this.audioPlayer.element.addEventListener('timeupdate', this._onTimeUpdate)
        this.audioPlayer.element.addEventListener('error', this._onError)
    }

    public async destroy() {
        // audio element events
        this.audioPlayer.element.removeEventListener('playing', this._onPlaying)
        this.audioPlayer.element.removeEventListener('pause', this._onPause)
        this.audioPlayer.element.removeEventListener('ended', this._onEnded)
        this.audioPlayer.element.removeEventListener('timeupdate', this._onTimeUpdate)
        this.audioPlayer.element.removeEventListener('error', this._onError)
        await this.stop()
    }


    //////////// PLAYBACK CONTROLS
    public async play() {
        if (!this.audioPlayer.initialized){
            try {
                this.setCurrentAudio()
            } catch (err) {
                console.error(err)
                return
            }
        }
        this.playPromise = this.audioPlayer.element.play()
        if (!this.playPromise) {
            return
        }
        this.playPromise
            .then(() => {
                this.audioPlayer.initialized = true
            })
            .catch((err) => { console.log(err); this.stop() })
    }

    public pause() {
        if (!this.playPromise) {
            return
        }
        this.audioPlayer.element.pause()
    }

    public async next() {
        if (!this.isHasNextAudio()) {
            this.audioPlayer.element.currentTime = 0
            return
        }
        this.audioPlayer.currentPlaying.index++
        this.setCurrentAudio()
        await this.play()
    }

    public async prev() {
        if (!this.isHasPrevAudio() || this.isAudioNotInStart()) {
            this.audioPlayer.element.currentTime = 0
            return
        }
        this.audioPlayer.currentPlaying.index--
        this.setCurrentAudio()
        await this.play()
    }

    public async stop() {
        this.audioPlayer.element.src = this.nullSrc
        this.audioPlayer.currentPlaying.index = 0
        this.audioPlayer.currentPlaying.duration = '00:00'
        this.audioPlayer.currentPlaying.currentTime = '00:00'
        this.audioPlayer.currentPlaying.percentsReached = 0
        this.audioPlayer.currentPlaying.percentsBuffered = 0
        this.audioPlayer.isProgressMouseDown = false
        this.audioPlayer.progressPreviewTime = 0
        this.audioPlayer.isVolumeMouseDown = false
        this.audioPlayer.initialized = false
        this.audioPlayer.isPlaying = false
        this.playPromise = undefined
    }

    //////////// PLAYBACK MANAGEMENT
    public setCurrentAudio(playlistIndex = this.audioPlayer.currentPlaying.index) {
        if (this.audioPlayer.currentPlaying.playlist.length < 1) {
            return Error('Audio: empty playlist')
        }
        this.audioPlayer.element.src = this.audioPlayer.currentPlaying.playlist[playlistIndex]
    }

    public async setPlaylist(playlist: string[]) {
        await this.stop()
        this.audioPlayer.currentPlaying.index = 0
        this.audioPlayer.currentPlaying.playlist = playlist
    }

    public addToPlaylist(url: string) {
        this.audioPlayer.currentPlaying.playlist.push(url)
    }

    // is audio not in start. Start = audio duration / 4.
    private isAudioNotInStart(): boolean {
        const isNotInStart = this.audioPlayer.element.duration / 4
        return this.audioPlayer.element.currentTime > isNotInStart
    }

    private isHasNextAudio(): boolean {
        const isHas = this.audioPlayer.currentPlaying.playlist[this.audioPlayer.currentPlaying.index + 1] !== undefined
        return isHas
    }

    private isHasPrevAudio(): boolean {
        const isHas = this.audioPlayer.currentPlaying.playlist[this.audioPlayer.currentPlaying.index - 1] !== undefined
        return isHas
    }

    public setTimeByPercents(percents: number) {
        const duration = this.audioPlayer.element.duration
        this.audioPlayer.element.currentTime = Service.round((duration / 100) * percents, 4)
    }

    public setVolumeByPercents(percents: number) {
        let percentsVolume = (percents / 100)
        if (percentsVolume > 1.0) {
            percentsVolume = 1.0
        }
        if (percentsVolume < 0) {
            percentsVolume = 0
        }
        this.audioPlayer.percentsVolume = percents
        this.audioPlayer.element.volume = percentsVolume
    }

    //////////// EVENTS
    private onPlaying() {
        this.audioPlayer.isPlaying = true
    }

    private onPause() {
        this.audioPlayer.isPlaying = false
    }

    private async onEnded() {
        this.audioPlayer.isPlaying = false
        await this.next()
    }

    private async onError(event) {
        // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/networkState
        switch (event.target.error.code) {
            case event.target.error.MEDIA_ERR_ABORTED:
                console.error('Audio: aborted')
                break
            case event.target.error.MEDIA_ERR_NETWORK:
                console.error('Audio: network error')
                break
            case event.target.error.MEDIA_ERR_DECODE:
                console.error('Audio: decode error. Maybe audio damaged or something?')
                break
            case event.target.error.MEDIA_ERR_SRC_NOT_SUPPORTED:
                console.log(this.audioPlayer.element.src)
                if (this.audioPlayer.element.src === 'https://null/') {
                    return
                }
                console.error('Audio: not supported')
                break
            default:
                console.error('Audio: unknown error')
                break
        }
    }

    private onTimeUpdate() {
        if (!this.audioPlayer.isProgressMouseDown) {
            this.audioPlayer.currentPlaying.percentsReached = Service.computePercents(this.audioPlayer.element.currentTime, this.audioPlayer.element.duration)
            this.audioPlayer.currentPlaying.percentsBuffered = Service.computeBuffered(this.audioPlayer.element)
        }
        if (this.audioPlayer.element.duration) {
            this.audioPlayer.currentPlaying.duration = Service.convertSeconds(this.audioPlayer.element.duration, 'auto')
            let mode: TConvertSecondsMode
            if (this.audioPlayer.element.duration < 3600) {
                mode = 'minutes'
            } else {
                mode = 'hours'
            }
            this.audioPlayer.currentPlaying.currentTime = Service.convertSeconds(this.audioPlayer.element.currentTime, mode)
        } else {
            this.audioPlayer.currentPlaying.duration = '00:00'
            this.audioPlayer.currentPlaying.currentTime = '00:00'
        }
    }

}
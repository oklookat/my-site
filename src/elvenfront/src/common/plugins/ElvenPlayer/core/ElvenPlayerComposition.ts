import {reactive} from "vue";

interface IAudioPlayer {
    active: boolean // is player component active
    initialized: boolean
    element: HTMLAudioElement
    isPlaying: boolean
    volume: number
    currentPlaying: {
        percentsReached: number // percents of audio reached
        percentsBuffered: number // percents of buffered audio
        playlist: string[]
        index: number
    }
    // moving
    isProgressMouseDown: boolean
    progressPreviewTime: number
    isVolumeMouseDown: boolean
}

export default class ElvenPlayerComposition {

    public audioPlayer: IAudioPlayer
    // saved events links (for cleanup)
    private _onPlaying = this.onPlaying.bind(this)
    private _onPause = this.onPause.bind(this)
    private _onEnded = this.onEnded.bind(this)
    private _onTimeUpdate = this.onTimeUpdate.bind(this)
    private _onError = this.onError.bind(this)

    constructor() {
        this.audioPlayer = {
            active: false,
            initialized: false,
            element: new Audio(''),
            isPlaying: false,
            volume: 1.0,
            currentPlaying: {
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

    public init(){
        // audio element events
        this.audioPlayer.element.addEventListener('playing', this._onPlaying)
        this.audioPlayer.element.addEventListener('pause', this._onPause)
        this.audioPlayer.element.addEventListener('ended', this._onEnded)
        this.audioPlayer.element.addEventListener('timeupdate', this._onTimeUpdate)
        this.audioPlayer.element.addEventListener('error', this._onError)
    }
    public async destroy(){
        // audio element events
        this.audioPlayer.element.removeEventListener('playing', this._onPlaying)
        this.audioPlayer.element.removeEventListener('pause', this._onPause)
        this.audioPlayer.element.removeEventListener('ended', this._onEnded)
        this.audioPlayer.element.removeEventListener('timeupdate', this._onTimeUpdate)
        this.audioPlayer.element.removeEventListener('error', this._onError)
        await this.stop()
        this.audioPlayer.initialized = false
    }


    //////////// PLAYBACK CONTROLS
    public async play() {
        if(!this.audioPlayer.initialized){
            this.setCurrentAudio()
            this.audioPlayer.initialized = true
        }
        await this.audioPlayer.element.play()
    }

    public pause() {
        this.audioPlayer.element.pause()
    }

    public async next() {
        if (!this.isHasNextAudio()) {
            await this.stop()
        } else {
            this.audioPlayer.currentPlaying.index++
            this.setCurrentAudio()
            await this.play()
        }
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
        await this.pause()
        this.audioPlayer.element.currentTime = 0
        this.audioPlayer.isPlaying = false
    }


    //////////// PLAYBACK MANAGEMENT
    public setCurrentAudio(playlistIndex = this.audioPlayer.currentPlaying.index) {
        if (this.audioPlayer.currentPlaying.playlist.length < 1) {
            return Error('E_PLAYLIST_EMPTY')
        }
        this.audioPlayer.element.src = this.audioPlayer.currentPlaying.playlist[playlistIndex]
        return true
    }

    public setPlaylist(playlist: string []) {
        this.audioPlayer.currentPlaying.index = 0
        this.audioPlayer.initialized = false
        this.audioPlayer.currentPlaying.playlist = playlist
    }

    public addToPlaylist(url: string) {
        this.audioPlayer.currentPlaying.playlist.push(url)
    }

    private isAudioNotInStart() {
        const isNotInStart = this.audioPlayer.element.duration / 4
        return this.audioPlayer.element.currentTime > isNotInStart
    }

    private isHasNextAudio() {
        const isHas = this.audioPlayer.currentPlaying.playlist[this.audioPlayer.currentPlaying.index + 1]
        return !!isHas
    }

    private isHasPrevAudio() {
        const isHas = this.audioPlayer.currentPlaying.playlist[this.audioPlayer.currentPlaying.index - 1]
        return !!isHas
    }

    //////////// EVENTS
    private onPlaying() {
        this.audioPlayer.isPlaying = true
    }

    private onPause() {
        this.audioPlayer.isPlaying = false
    }

    private async onEnded() {
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
                console.error('Audio: decode error. Audio damaged or not supported.')
                break
            case event.target.error.MEDIA_ERR_SRC_NOT_SUPPORTED:
                console.error('Audio: source not supported')
                break
            default:
                console.error('Audio: unknown error')
                break
        }
        await this.stop()
    }

    private onTimeUpdate() {
        if (!this.audioPlayer.isProgressMouseDown) {
            this.audioPlayer.currentPlaying.percentsReached = this.computePercents(this.audioPlayer.element.currentTime, this.audioPlayer.element.duration)
            this.audioPlayer.currentPlaying.percentsBuffered = this.computeBuffered(this.audioPlayer.element)
        }
    }




    //////////// SERVICE
    // computing how much buffered
    private computeBuffered(playerEL: HTMLAudioElement): number {
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

    public setTimeByPercents(percents: number){
        const duration = this.audioPlayer.element.duration
        this.audioPlayer.element.currentTime = this.round((duration / 100) * percents, 4)
    }

    // compute audio progress by duration and current time (percents)
    private computePercents(current: number, total: number): number {
        current = Math.round(current)
        let percents = (current / total) * 100
        if (percents >= 100) {
            percents = 100
        } else if (total < 1) {
            percents = 0
        }
        return this.round(percents, 4)
    }

    // round the number to the specific number of decimal places
    private round(value: number, precision: number): number {
        // https://stackoverflow.com/a/7343013/16762009
        const multiplier = Math.pow(10, precision || 0)
        return Math.round(value * multiplier) / multiplier
    }
}
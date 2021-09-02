import ElvenPlayerCore from "@/common/plugins/ElvenPlayer/core/ElvenPlayerCore";
import * as buffer from "buffer";

export default class ElvenPlayerWorker {

    private readonly core: ElvenPlayerCore

    // moving
    private isProgressMouseDown: boolean = false
    private progressPreviewTime: number = 0
    private isVolumeMouseDown: boolean = false
    private volume: number = 1.0

    // elements
    private isElementsInitialized: boolean = false
    public audioPlayerEL: HTMLAudioElement
    public containerEL: HTMLElement
    private progressContainerEL: HTMLElement
    private progressBarEL: HTMLElement
    private bufferedBarEL: HTMLElement
    private playEL: HTMLElement
    private pauseEL: HTMLElement
    private nextEL: HTMLElement
    private prevEL: HTMLElement
    private closeEL: HTMLElement
    private volumeContainerEL: HTMLElement
    private volumeSliderEL: HTMLElement
    private volumeSliderBubbleEL: HTMLElement

    // saved events links (for cleanup)
    private _documentMouseMove = this.onDocumentMouseMove.bind(this)
    private _documentMouseUp = this.onDocumentMouseUp.bind(this)


    constructor(coreContext) {
        this.core = coreContext
        this.initElements()
        this.audioPlayerEL = new Audio('')
    }

    // initialize

    public initElements() {
        const container = document.getElementById('audio-player-container')
        const progressContainer = document.getElementById('audio-player-progress-container')
        const progress = document.getElementById('audio-player-progressbar')
        const buffered = document.getElementById('audio-player-buffered')
        const play = document.getElementById('audio-player-play-butt')
        const pause = document.getElementById('audio-player-pause-butt')
        const next = document.getElementById('audio-player-next-butt')
        const prev = document.getElementById('audio-player-prev-butt')
        const close = document.getElementById('audio-player-close-butt')
        const volumeContainer = document.getElementById('volume-container')
        const volumeSlider = document.getElementById('volume-slider')
        const volumeSliderBubble = document.getElementById('volume-slider-bubble')
        const isElements = container && progressContainer && progress
            && buffered && play && pause && next && prev && close && volumeContainer
            && volumeSlider && volumeSliderBubble
        if (isElements) {
            this.containerEL = container
            this.progressContainerEL = progressContainer
            this.progressBarEL = progress
            this.bufferedBarEL = buffered
            this.playEL = play
            this.pauseEL = pause
            this.nextEL = next
            this.prevEL = prev
            this.closeEL = close
            this.volumeContainerEL = volumeContainer
            this.volumeSliderEL = volumeSlider
            this.volumeSliderBubbleEL = volumeSliderBubble
            this.isElementsInitialized = true
        } else {
            this.isElementsInitialized = false
            throw Error('E_INIT_ELEMENTS')
        }
    }

    public initEvents() {
        if (!this.isElementsInitialized) {
            this.initElements()
        }
        this.audioPlayerEL.addEventListener('playing', this.onPlaying.bind(this))
        this.audioPlayerEL.addEventListener('pause', this.onPause.bind(this))
        this.audioPlayerEL.addEventListener('ended', this.onEnded.bind(this))
        this.audioPlayerEL.addEventListener('timeupdate', this.onTimeupdate.bind(this))
        this.audioPlayerEL.addEventListener('error', this.onError.bind(this))
        // playback controls
        this.playEL.addEventListener('click', this.onPlayClick.bind(this))
        this.pauseEL.addEventListener('click', this.onPauseClick.bind(this))
        this.nextEL.addEventListener('click', this.onNextClick.bind(this))
        this.prevEL.addEventListener('click', this.onPrevClick.bind(this))
        this.closeEL.addEventListener('click', this.onCloseClick.bind(this))
        // rewind on progress click
        this.progressContainerEL.addEventListener('mousedown', this.onMovingMousedown.bind(this), {passive: false})
        this.progressContainerEL.addEventListener('touchstart', this.onMovingMousedown.bind(this), {passive: false})
        // volume on slider click
        this.volumeSliderEL.addEventListener('mousedown', this.onMovingMousedown.bind(this), {passive: false})
        this.volumeSliderEL.addEventListener('touchstart', this.onMovingMousedown.bind(this), {passive: false})
    }

    // cleanup

    public destroyEvents() {
        this.audioPlayerEL.removeEventListener('playing', this.onPlaying)
        this.audioPlayerEL.removeEventListener('pause', this.onPause)
        this.audioPlayerEL.removeEventListener('ended', this.onEnded)
        this.audioPlayerEL.removeEventListener('timeupdate', this.onTimeupdate)
        this.audioPlayerEL.removeEventListener('error', this.onError)
        this.playEL.removeEventListener('click', this.onPlayClick)
        this.pauseEL.removeEventListener('click', this.onPauseClick)
        this.nextEL.removeEventListener('click', this.onNextClick)
        this.prevEL.removeEventListener('click', this.onPrevClick)
        this.closeEL.removeEventListener('click', this.onCloseClick)
        this.progressContainerEL.removeEventListener('mousedown', this.onMovingMousedown)
        this.progressContainerEL.removeEventListener('touchstart', this.onMovingMousedown)
        this.volumeSliderEL.removeEventListener('mousedown', this.onMovingMousedown)
        this.volumeSliderEL.removeEventListener('touchstart', this.onMovingMousedown)
    }

    // playback

    private async onPlayClick() {
        await this.core.play()
        this.switchPlayPauseButtons()
    }

    private onPlaying() {
        this.core.isPlaying = true
        this.switchPlayPauseButtons()
    }

    private onPauseClick() {
        this.core.pause()
        this.switchPlayPauseButtons()
    }

    private onPause() {
        this.core.isPlaying = false
        this.switchPlayPauseButtons()
    }

    private async onNextClick() {
        await this.core.next()
    }

    private async onPrevClick() {
        await this.core.previous()
    }

    private async onCloseClick() {
        await this.core.destroy()
    }

    private async onEnded() {
        await this.core.next()
    }

    private onTimeupdate() {
        if (!this.isProgressMouseDown) {
            const percents = this.computePercents(this.audioPlayerEL.currentTime, this.audioPlayerEL.duration)
            this.progressBarEL.style.width = `${percents}%`
            this.computeBuffered()
        }
    }

    // service

    public switchPlayPauseButtons() {
        if (this.core.isPlaying) {
            this.playEL.style.display = "none"
            this.pauseEL.style.display = "block"
        } else {
            this.pauseEL.style.display = "none"
            this.playEL.style.display = "block"
        }
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
        await this.core.stop()
    }

    // moving stuff

    private onMovingMousedown(event: MouseEvent | TouchEvent) {
        // on user move audio progressbar or volume slider
        /// progress:
        // when user move mouse or finger down
        // we need to compute the preview, and when user up mouse or finger
        // we need to set position of the audio
        /// volume slider:
        // difference with progress - we need set volume immediately, without preview
        switch (event.target) {
            // audio progress triggered
            case this.progressContainerEL:
            case this.progressBarEL:
            case this.bufferedBarEL:
                this.isProgressMouseDown = true
                break
            // volume slider triggered
            case this.volumeSliderEL:
            case this.volumeSliderBubbleEL:
                this.isVolumeMouseDown = true
                break
            default:
                return
        }
        // pre-compute moving, because user already clicked
        this.onDocumentMouseMove(event)
        // setup document events, not local, because more comfortable control all this stuff when you moving as you like
        document.addEventListener("mousemove", this._documentMouseMove, {passive: false})
        document.addEventListener('mouseup', this._documentMouseUp, {passive: false})
        document.addEventListener("touchmove", this._documentMouseMove, {passive: false})
        document.addEventListener("touchend", this._documentMouseUp, {passive: false})
        document.addEventListener("touchcancel", this._documentMouseUp, {passive: false})
    }


    private onDocumentMouseMove(event) {
        this.computeMoving(event)
    }

    private onDocumentMouseUp() {
        // cleanup and set stuff
        if (this.isProgressMouseDown) {
            try {
                // on user up mouse or touch, we set audio time from preview
                this.audioPlayerEL.currentTime = this.progressPreviewTime
                this.isProgressMouseDown = false
            } catch (err) {
            }
        }
        this.movingCleanup()
    }

    private movingCleanup() {
        document.removeEventListener("mousemove", this._documentMouseMove)
        document.removeEventListener('mouseup', this._documentMouseMove)
        document.removeEventListener("touchmove", this._documentMouseMove)
        document.removeEventListener("touchend", this._documentMouseUp)
        document.removeEventListener("touchcancel", this._documentMouseUp)
        this.isProgressMouseDown = false
        this.isVolumeMouseDown = false
    }

    // get horizontal position of moving and compute preview
    private computeMoving(event) {
        event.preventDefault()
        let pageX = this.getPageX(event)
        if (!pageX) {
            return
        }
        if (this.isProgressMouseDown) {
            this.computeProgressPreview(pageX)
        } else if (this.isVolumeMouseDown) {
            this.computeVolume(pageX)
        }
    }

    private computeProgressPreview(pageX: number) {
        const clickPosition = this.getClickPosition(pageX, this.progressContainerEL)
        this.progressPreviewTime = clickPosition * this.audioPlayerEL.duration
        const percents = this.computePercents(this.progressPreviewTime, this.audioPlayerEL.duration)
        this.progressBarEL.style.width = `${percents}%`
    }

    private computeVolume(pageX: number) {
        const clickPosition = this.getClickPosition(pageX, this.volumeContainerEL)
        this.volume = Math.round(clickPosition * (this.audioPlayerEL.volume * 100))
        if(this.volume > 100){
            this.volume = 100
        }
        if(this.volume < 1){
            this.volume = 1
        }
        const percents = this.computePercents(this.volume, 100)
        this.volumeSliderEL.style.width = `${percents}%`
        this.volume = this.volume / 100
        if(this.volume < 0.1){
            this.volume = 0.1
        }
        this.audioPlayerEL.volume = this.volume
        console.log(`volume in player: ${this.volume}`)
    }


    // get pageX (horizontal mouse position) by touch or mouse event
    private getPageX(event): number | null {
        let pageX = null
        const isMovedByTouchscreen = event.type.includes('touch') && event.touches && event.touches.length > 0
        if (isMovedByTouchscreen) {
            for (const touch of event.touches) {
                if (touch.pageX) {
                    pageX = touch.pageX
                    break
                }
            }
        } else if (event.type.includes('mouse')) {
            // if moved by mouse
            pageX = event.pageX
        } else {
            // if moved by unknown - reset
            this.movingCleanup()
        }
        if (!pageX) {
            return null
        }
        return pageX
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
        return percents
    }

    // get click position by pageX
    private getClickPosition(pageX: number, element: HTMLElement): number {
        const clickPosition = (pageX - element.offsetLeft) / element.offsetWidth
        return clickPosition
    }


    // computing how much buffered
    private computeBuffered(currentTime = this.audioPlayerEL.currentTime) {
        currentTime = Math.round(currentTime)
        const duration = this.audioPlayerEL.duration
        if (duration > 0) {
            for (let i = 0; i < this.audioPlayerEL.buffered.length; i++) {
                const len = this.audioPlayerEL.buffered.length - 1 - i
                if (this.audioPlayerEL.buffered.start(len) < currentTime) {
                    const percents = this.computePercents(this.audioPlayerEL.buffered.end(len), duration)
                    this.bufferedBarEL.style.width = `${percents}%`
                    break
                }
            }
        }
    }


}
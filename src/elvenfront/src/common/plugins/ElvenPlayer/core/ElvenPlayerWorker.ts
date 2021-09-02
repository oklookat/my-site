import ElvenPlayerCore from "@/common/plugins/ElvenPlayer/core/ElvenPlayerCore";

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
    private playEL: HTMLElement
    private pauseEL: HTMLElement
    private nextEL: HTMLElement
    private prevEL: HTMLElement
    private closeEL: HTMLElement

    private progressSliderContainerEL: HTMLElement
    private progressSliderEL: HTMLElement
    private bufferedSliderEL: HTMLElement
    private progressSliderBubbleEL: HTMLElement

    private volumeSliderContainerEL: HTMLElement
    private volumeSliderEL: HTMLElement
    private volumeSliderBubbleEL: HTMLElement

    // saved events links (for cleanup)
    private _onPlaying = this.onPlaying.bind(this)
    private _onPause = this.onPause.bind(this)
    private _onEnded = this.onEnded.bind(this)
    private _onTimeUpdate = this.onTimeUpdate.bind(this)
    private _onError = this.onError.bind(this)
    private _onPlayClick = this.onPlayClick.bind(this)
    private _onPauseClick = this.onPauseClick.bind(this)
    private _onNextClick = this.onNextClick.bind(this)
    private _onPrevClick = this.onPrevClick.bind(this)
    private _onCloseClick = this.onCloseClick.bind(this)
    private _onMovingMouseDown = this.onMovingMouseDown.bind(this)
    private _documentMouseMove = this.onDocumentMouseMove.bind(this)
    private _documentMouseUp = this.onDocumentMouseUp.bind(this)


    constructor(coreContext) {
        this.core = coreContext
        this.initElements()
        this.audioPlayerEL = new Audio('')
        this.setVolume(this.audioPlayerEL.volume * 100)
    }

    // initialize

    public initElements() {
        const container = document.getElementById('audio-player-container')
        const play = document.getElementById('audio-player-play-butt')
        const pause = document.getElementById('audio-player-pause-butt')
        const next = document.getElementById('audio-player-next-butt')
        const prev = document.getElementById('audio-player-prev-butt')
        const close = document.getElementById('audio-player-close-butt')


        const progressContainer = document.querySelector('.playback-progressbar')
        if(progressContainer){
            this.progressSliderContainerEL = progressContainer
            const buffered = progressContainer.querySelector('.sl-buffered')
            const slider = progressContainer.querySelector('.sl-slider')
            const sliderBubble = progressContainer.querySelector('.sl-slider-bubble')
            if(buffered && slider && sliderBubble){
                this.bufferedSliderEL = buffered
                this.progressSliderEL = slider
                this.progressSliderBubbleEL = sliderBubble
            }
        }
        const volumeSliderContainer = document.querySelector('.volume-slider')
        if(volumeSliderContainer){
            this.volumeSliderContainerEL = volumeSliderContainer
            const slider = volumeSliderContainer.querySelector('.sl-slider')
            const sliderBubble = volumeSliderContainer.querySelector('.sl-slider-bubble')
            if(slider && sliderBubble){
                this.volumeSliderEL = slider
                this.volumeSliderBubbleEL = sliderBubble
            }
        }

        const isElements = container && play && pause && next && prev && close
        if (isElements) {
            this.containerEL = container
            this.playEL = play
            this.pauseEL = pause
            this.nextEL = next
            this.prevEL = prev
            this.closeEL = close
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
        // audio element events
        this.audioPlayerEL.addEventListener('playing', this._onPlaying)
        this.audioPlayerEL.addEventListener('pause', this._onPause)
        this.audioPlayerEL.addEventListener('ended', this._onEnded)
        this.audioPlayerEL.addEventListener('timeupdate', this._onTimeUpdate)
        this.audioPlayerEL.addEventListener('error', this._onError)
        // playback controls
        this.playEL.addEventListener('click', this._onPlayClick)
        this.pauseEL.addEventListener('click', this._onPauseClick)
        this.nextEL.addEventListener('click', this._onNextClick)
        this.prevEL.addEventListener('click', this._onPrevClick)
        this.closeEL.addEventListener('click', this._onCloseClick)

        // rewind on progress click / press
        this.progressSliderContainerEL.addEventListener('mousedown', this._onMovingMouseDown, {passive: false})
        this.progressSliderContainerEL.addEventListener('touchstart', this._onMovingMouseDown, {passive: false})
        this.progressSliderBubbleEL.addEventListener('mousedown', this._onMovingMouseDown, {passive: false})
        this.progressSliderBubbleEL.addEventListener('touchstart', this._onMovingMouseDown, {passive: false})
        // volume slider click / press
        this.volumeSliderContainerEL.addEventListener('mousedown', this._onMovingMouseDown, {passive: false})
        this.volumeSliderEL.addEventListener('mousedown', this._onMovingMouseDown, {passive: false})
        this.volumeSliderEL.addEventListener('touchstart', this._onMovingMouseDown, {passive: false})
        this.volumeSliderBubbleEL.addEventListener('mousedown', this._onMovingMouseDown, {passive: false})
        this.volumeSliderBubbleEL.addEventListener('touchstart', this._onMovingMouseDown, {passive: false})
    }

    // cleanup

    public destroyEvents() {
        // audio element events
        this.audioPlayerEL.removeEventListener('playing', this._onPlaying)
        this.audioPlayerEL.removeEventListener('pause', this._onPause)
        this.audioPlayerEL.removeEventListener('ended', this._onEnded)
        this.audioPlayerEL.removeEventListener('timeupdate', this._onTimeUpdate)
        this.audioPlayerEL.removeEventListener('error', this._onError)
        // playback controls
        this.playEL.removeEventListener('click', this._onPlayClick)
        this.pauseEL.removeEventListener('click', this._onPauseClick)
        this.nextEL.removeEventListener('click', this._onNextClick)
        this.prevEL.removeEventListener('click', this._onPrevClick)
        this.closeEL.removeEventListener('click', this._onCloseClick)
        // rewind on progress click / press
        this.progressSliderContainerEL.removeEventListener('mousedown', this._onMovingMouseDown)
        this.progressSliderContainerEL.removeEventListener('touchstart', this._onMovingMouseDown)
        // volume slider click / press
        this.volumeSliderContainerEL.removeEventListener('mousedown', this._onMovingMouseDown)
        this.volumeSliderEL.removeEventListener('mousedown', this._onMovingMouseDown)
        this.volumeSliderEL.removeEventListener('touchstart', this._onMovingMouseDown)
        this.volumeSliderBubbleEL.removeEventListener('mousedown', this._onMovingMouseDown)
        this.volumeSliderBubbleEL.removeEventListener('touchstart', this._onMovingMouseDown)
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

    private onTimeUpdate() {
        if (!this.isProgressMouseDown) {
            const percents = this.computePercents(this.audioPlayerEL.currentTime, this.audioPlayerEL.duration)
            this.progressSliderEL.style.width = `${percents}%`
            this.progressSliderBubbleEL.style.left = `${percents}%`
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

    private onMovingMouseDown(event: MouseEvent | TouchEvent) {
        // on user move audio progressbar or volume slider
        /// progress:
        // when user move mouse or finger down
        // we need to compute the preview, and when user up mouse or finger
        // we need to set position of the audio
        /// volume slider:
        // difference with progress - we need set volume immediately, without preview
        switch (event.target) {
            // audio progress triggered
            case this.progressSliderContainerEL:
            case this.progressSliderEL:
            case this.bufferedSliderEL:
            case this.progressSliderBubbleEL:
                this.isProgressMouseDown = true
                break
            // volume slider triggered
            case this.volumeSliderContainerEL:
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

    // progressbar, on user move mouse
    private computeProgressPreview(pageX: number) {
        const clickPosition = this.getClickPosition(pageX, this.progressSliderContainerEL)
        this.progressPreviewTime = clickPosition * this.audioPlayerEL.duration
        let percents = this.computePercents(this.progressPreviewTime, this.audioPlayerEL.duration)
        if(percents > 100){
            percents = 100
        }
        if(percents < 0){
            percents = 0
        }
        this.progressSliderEL.style.width = `${percents}%`
        this.progressSliderBubbleEL.style.left = `calc(${percents}% - 6px)`
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
                    this.bufferedSliderEL.style.width = `${percents}%`
                    break
                }
            }
        }
    }

    // VOLUME EVENTS START //
    private computeVolume(pageX: number) {
        const clickPosition = this.getClickPosition(pageX, this.volumeSliderContainerEL)
        this.volume = Math.round(clickPosition * 100)
        const percents = this.computePercents(this.volume, 100)
        this.setVolume(percents)
        console.log(`volume in player: ${this.volume}`)
    }

    private setVolume(percents: number){
        if(percents > 100){
            percents = 100
            this.volume = 100
        }
        if(percents < 0){
            percents = 0
            this.volume = 0
        }
        this.volumeSliderEL.style.width = `${percents}%`
        this.volumeSliderBubbleEL.style.left = `calc(${percents}% - 6px)` // -6 = margin
        this.volume = percents / 100
        if(this.volume > 1){
            this.volume = 1
        } else if(this.volume < 0){
            this.volume = 0
        }
        this.audioPlayerEL.volume = this.volume
    }
    // VOLUME EVENTS END //


    // SERVICE START //
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
    // SERVICE END //


}
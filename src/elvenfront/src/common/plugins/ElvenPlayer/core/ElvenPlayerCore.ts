let _this

export default class ElvenPlayerCore {

    private registered = {}
    public initialized = false

    // ELEMENTS START //
    private isElements = false
    private audioPlayerEL: HTMLAudioElement
    private containerEL: HTMLElement
    private progressContainerEL: HTMLElement
    private progressBarEL: HTMLElement
    private bufferedBarEL: HTMLElement
    private playEL: HTMLElement
    private pauseEL: HTMLElement
    private nextEL: HTMLElement
    private prevEL: HTMLElement
    private closeEL: HTMLElement
    // ELEMENTS END //

    // PLAYBACK START //
    public isPlaying = false
    private playlist: string[]
    private currentPlaying = {
        index: 0,
    }
    // PLAYBACK END //

    //////////////// SERVICE START
    private progressMouseDown = false
    private progressMouseDownTempTime = 0

    constructor() {
        _this = this
        this.initElements()
        this.audioPlayerEL = new Audio('')
        this.playlist = []
    }

    private initElements() {
        const container = document.getElementById('audio-player-container')
        const progressContainer = document.getElementById('audio-player-progress-container')
        const progress = document.getElementById('audio-player-progressbar')
        const buffered = document.getElementById('audio-player-buffered')
        const play = document.getElementById('audio-player-play-butt')
        const pause = document.getElementById('audio-player-pause-butt')
        const next = document.getElementById('audio-player-next-butt')
        const prev = document.getElementById('audio-player-prev-butt')
        const close = document.getElementById('audio-player-close-butt')
        if (container && progressContainer && progress && buffered && play && pause && next && prev && close) {
            this.containerEL = container
            this.progressContainerEL = progressContainer
            this.progressBarEL = progress
            this.bufferedBarEL = buffered
            this.playEL = play
            this.pauseEL = pause
            this.nextEL = next
            this.prevEL = prev
            this.closeEL = close
            this.isElements = true
        } else {
            this.isElements = false
            throw Error('E_INIT_ELEMENTS')
        }
    }

    private initEvents() {
        if (!this.isElements) {
            this.initElements()
        }
        this.progressContainerEL.addEventListener('mousedown', progressContainerELMousedown, {passive: false})
        this.progressContainerEL.addEventListener('touchstart', progressContainerELMousedown, {passive: false})
        this.audioPlayerEL.addEventListener('playing', audioPlayerELPlaying)
        this.audioPlayerEL.addEventListener('pause', audioPlayerELPause)
        this.audioPlayerEL.addEventListener('ended', audioPlayerELEnded)
        this.audioPlayerEL.addEventListener('timeupdate', audioPlayerELTimeupdate)
        this.audioPlayerEL.addEventListener('error', onError)
        this.playEL.addEventListener('click', playELClick)
        this.pauseEL.addEventListener('click', pauseELClick)
        this.nextEL.addEventListener('click', nextELClick)
        this.prevEL.addEventListener('click', prevELClick)
        this.closeEL.addEventListener('click', closeELClick)
    }

    private destroyEvents() {
        this.progressContainerEL.removeEventListener('mousedown', progressContainerELMousedown)
        this.progressContainerEL.removeEventListener('touchstart', progressContainerELMousedown)
        this.audioPlayerEL.removeEventListener('playing', audioPlayerELPlaying)
        this.audioPlayerEL.removeEventListener('pause', audioPlayerELPause)
        this.audioPlayerEL.removeEventListener('ended', audioPlayerELEnded)
        this.audioPlayerEL.removeEventListener('timeupdate', audioPlayerELTimeupdate)
        this.audioPlayerEL.removeEventListener('error', onError)
        this.playEL.removeEventListener('click', playELClick)
        this.pauseEL.removeEventListener('click', pauseELClick)
        this.nextEL.removeEventListener('click', nextELClick)
        this.prevEL.removeEventListener('click', prevELClick)
        this.closeEL.removeEventListener('click', closeELClick)
    }

    //////////////// SERVICE END


    private setCurrentAudio(playlistIndex = this.currentPlaying.index) {
        if (this.playlist.length < 1) {
            return Error('E_PLAYLIST_EMPTY')
        }
        this.audioPlayerEL.src = this.playlist[playlistIndex]
        return true
    }

    private isAudioNotInStart() {
        const isNotInStart = this.audioPlayerEL.duration / 4
        return this.audioPlayerEL.currentTime > isNotInStart
    }

    private isHasNextAudio() {
        const isHas = this.playlist[this.currentPlaying.index + 1]
        return !!isHas
    }

    private isHasPrevAudio() {
        const isHas = this.playlist[this.currentPlaying.index - 1]
        return !!isHas
    }


    // PLAYBACK CONTROLS START //

    public async play() {
        if (!this.initialized) {
            this.initEvents()
            this.containerEL.style.display = 'flex'
            this.setCurrentAudio()
            this.initialized = true
        }
        await this.audioPlayerEL.play()
            .catch(() =>{
               this.stop()
            })
    }

    public pause() {
        this.audioPlayerEL.pause()
    }

    public async next() {
        if (!this.isHasNextAudio()) {
            await this.stop()
        } else {
            this.currentPlaying.index++
            this.setCurrentAudio()
            await this.play()
        }
    }

    public async previous() {
        if (!this.isHasPrevAudio() || this.isAudioNotInStart()) {
            // repeat if current song not in start
            // repeat if no songs in prev
            this.audioPlayerEL.currentTime = 0
            return
        }
        this.currentPlaying.index--
        this.setCurrentAudio()
        await this.play()
    }

    public async stop() {
        await this.pause()
        this.audioPlayerEL.currentTime = 0
    }

    public destroy() {
        this.audioPlayerEL.pause()
        this.audioPlayerEL.currentTime = 0
        this.playlist = []
        this.audioPlayerEL.src = ''
        this.containerEL.style.display = 'none'
        this.destroyEvents()
        this.initialized = false
    }

    public setPlaylist(playlist: string []) {
        this.currentPlaying.index = 0
        this.initialized = false
        this.playlist = playlist
    }

    public addToPlaylist(url: string) {
        this.playlist.push(url)
    }

    // PLAYBACK CONTROLS END //

}

function onError(event) {
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
    _this.stop()
}

function computeProgress(currentTime = _this.audioPlayerEL.currentTime) {
    currentTime = Math.round(currentTime)
    const duration = _this.audioPlayerEL.duration
    let progressPercents = (currentTime / duration) * 100
    if (progressPercents >= 100) {
        progressPercents = 100
    }
    if (duration > 0) {
        _this.progressBarEL.style.width = `${progressPercents}%`
    }
}


function computeMoveProgress(event) {
    if (!_this.progressMouseDown) {
        return
    }
    let pageX
    if (event.type.includes('touch') && event.touches && event.touches.length > 0) {
        // if move by touchscreen
        for (const touch of event.touches) {
            if (touch.pageX) {
                pageX = touch.pageX
                break
            }
        }
    } else if (event.type.includes('mouse')) {
        // if move by mouse
        pageX = event.pageX
    } else {
        documentMouseup()
    }
    if (!pageX) {
        return
    }
    const clickPosition = (pageX - _this.progressContainerEL.offsetLeft) / _this.progressContainerEL.offsetWidth
    _this.progressMouseDownTempTime = clickPosition * _this.audioPlayerEL.duration
    computeProgress(_this.progressMouseDownTempTime)
}


function computeBuffered(currentTime = _this.audioPlayerEL.currentTime) {
    currentTime = Math.round(currentTime)
    const duration = _this.audioPlayerEL.duration
    if (duration > 0) {
        for (let i = 0; i < _this.audioPlayerEL.buffered.length; i++) {
            const len = _this.audioPlayerEL.buffered.length - 1 - i
            if (_this.audioPlayerEL.buffered.start(len) < currentTime) {
                let bufferPercents = (_this.audioPlayerEL.buffered.end(len) / duration) * 100
                if (bufferPercents >= 100) {
                    bufferPercents = 100
                }
                _this.bufferedBarEL.style.width = `${bufferPercents}%`
                break
            }
        }
    }
}

function progressContainerELMousedown(event) {
    event.preventDefault()
    _this.progressMouseDown = true
    computeMoveProgress(event)
    document.addEventListener("mousemove", documentMousemove, {passive: false})
    document.addEventListener('mouseup', documentMouseup, {passive: false})
    document.addEventListener("touchmove", documentMousemove, {passive: false})
    document.addEventListener("touchend", documentMouseup, {passive: false})
    document.addEventListener("touchcancel", documentMouseup, {passive: false})
}

function documentMousemove(event) {
    if (_this.progressMouseDown) {
        event.preventDefault()
        computeMoveProgress(event)
    }
}

function documentMouseup(event?) {
    if (_this.progressMouseDown) {
        if (event && event.cancelable) {
            event.preventDefault()
        }
        document.removeEventListener("mousemove", documentMousemove)
        document.removeEventListener('mouseup', documentMouseup)
        document.removeEventListener("touchmove", documentMousemove)
        document.removeEventListener("touchend", documentMouseup)
        _this.progressMouseDown = false
        _this.audioPlayerEL.currentTime = _this.progressMouseDownTempTime
    }
}

function switchPlayPauseButtons() {
    if (_this.isPlaying) {
        _this.playEL.style.display = "none"
        _this.pauseEL.style.display = "block"
    } else {
        _this.pauseEL.style.display = "none"
        _this.playEL.style.display = "block"
    }
}

function audioPlayerELPlaying() {
    _this.isPlaying = true
    switchPlayPauseButtons()
}

function audioPlayerELPause() {
    _this.isPlaying = false
    switchPlayPauseButtons()
}

function audioPlayerELEnded() {
    _this.next()
}

function audioPlayerELTimeupdate() {
    if (!_this.progressMouseDown) {
        computeProgress()
        computeBuffered()
    }
}


function playELClick() {
    _this.play()
    switchPlayPauseButtons()
}

function pauseELClick() {
    _this.pause()
    switchPlayPauseButtons()
}

function nextELClick() {
    _this.next()
}

function prevELClick() {
    _this.previous()
}

function closeELClick() {
    _this.destroy()
}
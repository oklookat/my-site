let _global_this

export default class ElvenPlayerCore {

    private registered = {}
    public initialized = false

    public isPlaying = false
    private playlist: string[]
    private currentPlaying = {
        index: 0,
    }

    // ELEMENTS START //
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

    private progressMouseDown = false
    private progressMouseDownTempTime = 0


    public addEventListener(name, callback) {
        if (!this.registered[name]) this.registered[name] = []
        this.registered[name].push(callback)
    }

    private triggerEvent(name, args) {
        this.registered[name]?.forEach(fnc => fnc.apply(this, args))
    }

    constructor() {
        this.initElements()
        this.audioPlayerEL = new Audio('')
        this.initEvents()
        this.playlist = []
        _global_this = this
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
        } else {
            throw Error('E_INIT_ELEMENTS')
        }
    }

    private initEvents() {
        // INIT STATE START //
        this.audioPlayerEL.addEventListener('playing', () => {
            this.isPlaying = true
            this.switchPlayPauseButtons()
        })
        this.audioPlayerEL.addEventListener('pause', () => {
            this.isPlaying = false
            this.switchPlayPauseButtons()
        })
        // INIT STATE END //

        // INIT PROGRESS AND ERROR HANDLING START //
        this.audioPlayerEL.addEventListener('ended', onEnded)
        this.audioPlayerEL.addEventListener('timeupdate', () =>{
            if(!this.progressMouseDown){
                computeProgress()
                computeBuffered()
            }
        })
        this.audioPlayerEL.addEventListener('error', onError)
        // INIT PROGRESS AND ERROR HANDLING END //

        // INIT UI CONTROLS START //
        this.progressContainerEL.addEventListener("mousedown", (event) =>{
            event.preventDefault()
            this.progressMouseDown = true
            this.computeClickProgress(event)
        })
        document.addEventListener("mousemove", (event) =>{
            if(this.progressMouseDown){
                this.computeClickProgress(event)
            }
        })
        document.addEventListener('mouseup', () =>{
            if(this.progressMouseDown){
                this.progressMouseDown = false
                this.audioPlayerEL.currentTime = this.progressMouseDownTempTime
            }
        })
        this.playEL.addEventListener('click', () => {
            this.play()
            this.switchPlayPauseButtons()
        })
        this.pauseEL.addEventListener('click', () => {
            this.pause()
            this.switchPlayPauseButtons()
        })
        this.nextEL.addEventListener('click', () => {
            this.next()
        })
        this.prevEL.addEventListener('click', () => {
            this.previous()
        })
        this.closeEL.addEventListener('click', () => {
            this.destroy()
        })
        // INIT UI CONTROLS END //
    }

    private computeClickProgress(event){
        if(this.progressMouseDown){
            const clickPosition = (event.pageX - this.progressContainerEL.offsetLeft) / this.progressContainerEL.offsetWidth
            this.progressMouseDownTempTime = clickPosition * this.audioPlayerEL.duration
            computeProgress(this.progressMouseDownTempTime)
        }
    }

    private switchPlayPauseButtons() {
        if (this.isPlaying) {
            this.playEL.style.display = "none"
            this.pauseEL.style.display = "block"
        } else {
            this.pauseEL.style.display = "none"
            this.playEL.style.display = "block"
        }
    }

    private setCurrentAudio(playlistIndex = this.currentPlaying.index) {
        if (this.playlist.length < 1) {
            return Error('E_PLAYLIST_EMPTY')
        }
        this.audioPlayerEL.src = this.playlist[playlistIndex]
        return true
    }

    private isAudioNotInStart(){
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
            this.containerEL.style.display = 'flex'
            this.setCurrentAudio()
            this.initialized = true
        }
        await this.audioPlayerEL.play()
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
            await this.audioPlayerEL.play()
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
        await this.audioPlayerEL.play()
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



function computeProgress(currentTime = _global_this.audioPlayerEL.currentTime){
    currentTime = Math.round(currentTime)
    const duration = _global_this.audioPlayerEL.duration
    let progressPercents = (currentTime / duration) * 100
    if (progressPercents >= 100) {
        progressPercents = 100
    }
    if(duration > 0){
        _global_this.progressBarEL.style.width = `${progressPercents}%`
    }
}

function computeBuffered(currentTime = _global_this.audioPlayerEL.currentTime){
    currentTime = Math.round(currentTime)
    const duration = _global_this.audioPlayerEL.duration
    if (duration > 0) {
        for (let i = 0; i < _global_this.audioPlayerEL.buffered.length; i++) {
            const len = _global_this.audioPlayerEL.buffered.length - 1 - i
            if (_global_this.audioPlayerEL.buffered.start(len) < currentTime) {
                let bufferPercents = (_global_this.audioPlayerEL.buffered.end(len) / duration) * 100
                if (bufferPercents >= 100) {
                    bufferPercents = 100
                }
                _global_this.bufferedBarEL.style.width = `${bufferPercents}%`
                break
            }
        }
    }
}

function onEnded() {
    _global_this.next()
}

function onError() {
    // https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/networkState
    const networkState = _global_this.audioPlayerEL.networkState
    if (networkState === 3) {
        return
    }
    throw Error('E_AUDIO_LOADING')
}
let _global_this

export default class ElvenPlayerCore {

    private registered = {}
    public initialized = false
    public isPlaying = false

    private playlist: string[]
    private currentPlaying = {
        index: 0,
    }

    private audioPlayer: HTMLAudioElement
    private progressBar: Element
    private bufferedBar: Element

    public addEventListener(name, callback) {
        if (!this.registered[name]) this.registered[name] = []
        this.registered[name].push(callback)
    }

    private triggerEvent(name, args) {
        this.registered[name]?.forEach(fnc => fnc.apply(this, args))
    }

    constructor() {
        this.initElements()
        this.audioPlayer = new Audio('')
        this.initEvents()
        this.playlist = []
        _global_this = this
    }

    private initElements(){
        const progress = document.querySelector('.audio-player-progressbar')
        const buffered = document.querySelector('.audio-player-buffered')
        if (progress && buffered) {
            this.progressBar = progress
            this.bufferedBar = buffered
        } else {
            throw Error('audio-player-progressbar or audio-player-buffered classes not found.')
        }
    }

    private initEvents() {
        this.audioPlayer.addEventListener('playing', () => {
            this.isPlaying = true
            this.triggerEvent('playPause', [this.isPlaying])
        })
        this.audioPlayer.addEventListener('pause', () => {
            this.isPlaying = false
            this.triggerEvent('playPause', [this.isPlaying])
        })
        this.audioPlayer.addEventListener('ended', onEnded)
        this.audioPlayer.addEventListener('timeupdate', onTimeUpdate)
        this.audioPlayer.addEventListener('error', onError)
    }

    public getPlaylistLength() {
        return this.playlist.length
    }

    public setPlaylist(playlist: string []) {
        this.currentPlaying.index = 0
        this.initialized = false
        this.playlist = playlist
    }

    public addToPlaylist(url: string) {
        this.playlist.push(url)
    }

    private setCurrentAudio(sourceIndex = this.currentPlaying.index) {
        if (this.playlist.length < 1) {
            return Error('PLAYLIST_EMPTY')
        }
        this.audioPlayer.src = this.playlist[sourceIndex]
        return true
    }

    private isHasNextAudio() {
        const isHas = this.playlist[this.currentPlaying.index + 1]
        return !!isHas
    }

    private async repeat() {
        this.setCurrentAudio()
        this.audioPlayer.currentTime = 0
        await this.play()
    }


    // PLAYBACK CONTROLS START //
    public async play() {
        if (!this.initialized) {
            this.setCurrentAudio()
            this.initialized = true
        }
        await this.audioPlayer.play()
    }

    public pause() {
        this.audioPlayer.pause()
    }

    public async next() {
        if (!this.isHasNextAudio()) {
            await this.repeat()
        } else {
            this.currentPlaying.index++
            this.setCurrentAudio()
            await this.audioPlayer.play()
        }
    }

    public async previous() {
        const durationCF = this.audioPlayer.duration / 4
        if (this.audioPlayer.currentTime > durationCF) {
            this.audioPlayer.currentTime = 0
            return
        }
        if (!this.playlist[this.currentPlaying.index - 1]) {
            this.audioPlayer.currentTime = 0
            return
        }
        this.currentPlaying.index--
        this.setCurrentAudio()
        await this.audioPlayer.play()
    }

    public stop(){
        this.playlist = []
        this.audioPlayer.src = ''
    }

    public async setDurationOnClick(event) {
        let parentWidth
        if (!this.progressBar.parentElement) {
            parentWidth = window.screenX
        } else {
            parentWidth = this.progressBar.parentElement.clientWidth
        }
        const percents = Math.ceil((event.screenX * 100) / parentWidth)
        this.audioPlayer.currentTime = Math.ceil((this.audioPlayer.duration * percents) / 100)
    }

    // PLAYBACK CONTROLS END //

}


function onTimeUpdate() {
    const currentTime = Math.ceil(_global_this.audioPlayer.currentTime)
    const duration = _global_this.audioPlayer.duration
    let progressPercents = (currentTime / duration) * 100
    if(progressPercents >= 100){
        progressPercents = 100
    }
    if (duration > 0) {
        _global_this.progressBar.style.width = `${progressPercents}%`
        for (let i = 0; i < _global_this.audioPlayer.buffered.length; i++) {
            const len = _global_this.audioPlayer.buffered.length - 1 - i
            if (_global_this.audioPlayer.buffered.start(len) < currentTime) {
                let bufferPercents = (_global_this.audioPlayer.buffered.end(len) / duration) * 100
                if(bufferPercents >= 100){
                    bufferPercents = 100
                }
                _global_this.bufferedBar.style.width = `${bufferPercents}%`
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
    const networkState = _global_this.audioPlayer.networkState
    if (networkState === 3) {
        return
    }
    throw Error('При загрузке аудио произошла ошибка')
}
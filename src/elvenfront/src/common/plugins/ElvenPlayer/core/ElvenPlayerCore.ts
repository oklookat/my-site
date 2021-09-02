import ElvenPlayerWorker from "@/common/plugins/ElvenPlayer/core/ElvenPlayerWorker";

export default class ElvenPlayerCore {

    public playerInitialized: boolean = false
    private playerWorker: ElvenPlayerWorker

    // playback
    public isPlaying: boolean = false
    private playlist: string[] = []
    private currentPlaying = {
        index: 0,
    }

    constructor() {
        this.playerWorker = new ElvenPlayerWorker(this)
    }

    private setCurrentAudio(playlistIndex = this.currentPlaying.index) {
        if (this.playlist.length < 1) {
            return Error('E_PLAYLIST_EMPTY')
        }
        this.playerWorker.audioPlayerEL.src = this.playlist[playlistIndex]
        this.playerWorker.switchPlayPauseButtons()
        return true
    }

    private isAudioNotInStart() {
        const isNotInStart = this.playerWorker.audioPlayerEL.duration / 4
        return this.playerWorker.audioPlayerEL.currentTime > isNotInStart
    }

    private isHasNextAudio() {
        const isHas = this.playlist[this.currentPlaying.index + 1]
        return !!isHas
    }

    private isHasPrevAudio() {
        const isHas = this.playlist[this.currentPlaying.index - 1]
        return !!isHas
    }


    // playback controls

    public async play() {
        if (!this.playerInitialized) {
            this.playerWorker.initEvents()
            this.playerWorker.containerEL.style.display = 'flex'
            this.setCurrentAudio()
            this.playerInitialized = true
        }
        await this.playerWorker.audioPlayerEL.play()
            .catch(() => {
                this.stop()
            })
    }

    public pause() {
        this.playerWorker.audioPlayerEL.pause()
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
            this.playerWorker.audioPlayerEL.currentTime = 0
            return
        }
        this.currentPlaying.index--
        this.setCurrentAudio()
        await this.play()
    }

    public async stop() {
        await this.pause()
        this.playerWorker.audioPlayerEL.currentTime = 0
        // reset all for recover in future.
        // If no set this, after an error you will have to reload the page to try again
        this.playerInitialized = false
        this.isPlaying = false
    }

    public async destroy() {
        this.playerWorker.destroyEvents()
        await this.stop()
        this.playerWorker.containerEL.style.display = 'none'
    }

    public setPlaylist(playlist: string []) {
        this.currentPlaying.index = 0
        this.playerInitialized = false
        this.playlist = playlist
    }

    public addToPlaylist(url: string) {
        this.playlist.push(url)
    }

    // PLAYBACK CONTROLS END //
}
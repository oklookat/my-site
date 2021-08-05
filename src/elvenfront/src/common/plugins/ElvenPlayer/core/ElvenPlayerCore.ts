let _global_this

export default class ElvenPlayerCore {

    private initialized = false
    public isPlaying = false
    private audioSources: string[]
    private audioPlayer: HTMLAudioElement
    private progressBar: Element

    private totalTracks = -1
    private currentPlaying = {
        index: 0,
    }

    constructor() {
        _global_this = this
        this.audioPlayer = new Audio('')
        this.audioPlayer.addEventListener('timeupdate', onPlaying)
        this.audioPlayer.addEventListener('playing', () =>{
            this.isPlaying = true
        })
        const progress = document.querySelector('.audio-player-progressbar')
        if(progress){
            this.progressBar = progress
        } else {
            throw Error('audio-player-progressbar class not found.')
        }
        this.audioSources = []
    }

    public addSource(url: string){
        this.audioSources.push(url)
        this.totalTracks++
    }

    private async setSource(sourceIndex = this.currentPlaying.index){
        this.audioPlayer.src = this.audioSources[sourceIndex]
    }

    // PLAYBACK CONTROLS START //

    public async play(){
        if(!this.initialized){
            await this.setSource()
            this.initialized = true
        }
        await this.audioPlayer.play()
    }

    public async pause(){
        await this.audioPlayer.pause()
        this.isPlaying = false
    }

    public async next(){
        console.log('NEXT')
        if(!this.audioSources[this.currentPlaying.index + 1]){
           return Promise.reject('NEXT_NO_TRACKS')
        }
        this.currentPlaying.index++
        await this.setSource()
        await this.audioPlayer.play()
    }

    public async previous(){
        console.log('PREV')
        if(!this.audioSources[this.currentPlaying.index - 1]){
            return Promise.reject('PREV_NO_TRACKS')
        }
        this.currentPlaying.index--
        await this.setSource()
        await this.audioPlayer.play()
    }

    public async setDurationOnClick(event){
        let parentWidth = this.progressBar.parentElement.clientWidth
        if(!parentWidth){
            parentWidth = window.screenX
        }
        const percents = Math.ceil((event.screenX * 100) / parentWidth)
        const _calc = Math.ceil((this.audioPlayer.duration * percents) / 100)
        this.audioPlayer.currentTime = _calc
    }
    // PLAYBACK CONTROLS END //

}



function onPlaying(){
    const currentTime = Math.ceil(_global_this.audioPlayer.currentTime)
    const duration = Math.ceil(_global_this.audioPlayer.duration)
    _global_this.progressBar.style.width = `${Math.ceil((currentTime * 100) / duration)}%`
}
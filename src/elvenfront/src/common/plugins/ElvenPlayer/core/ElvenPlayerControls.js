import ElvenPlayerCore from "@/common/plugins/ElvenPlayer/core/ElvenPlayerCore"

export default class ElvenPlayerControls {
    #registered = {}

    audioPlayer = undefined
    isPlaying = false

    addEventListener(name, callback){
        if (!this.#registered[name]) this.#registered[name] = []
        this.#registered[name].push(callback)
    }
    #triggerEvent(name, args) {
        this.#registered[name]?.forEach(fnc => fnc.apply(this, args))
    }

    constructor() {
        this.audioPlayer = new ElvenPlayerCore()
        this.audioPlayer.addEventListener('playPause', (event) =>{
            this.#triggerEvent('playPause', [event])
        })
    }

    addToPlaylist(url) {
        this.audioPlayer.addToPlaylist(url)
    }

    setPlaylist(playlist){
        this.audioPlayer.setPlaylist(playlist)
    }

    play(){
        this.audioPlayer.play()
    }

    playPause() {
        if(!this.audioPlayer){
             return null
         }
        if (!this.audioPlayer.isPlaying) {
             this.audioPlayer.play()
        } else {
             this.audioPlayer.pause()
        }
        return !this.audioPlayer.isPlaying
    }

    next() {
        this.audioPlayer.next()
    }

    prev() {
        this.audioPlayer.previous()
    }

    stop(){
        this.audioPlayer.stop()
    }

    setDurationOnClick(event){
        this.audioPlayer.setDurationOnClick(event)
    }

    getPlaylistLength(){
        return this.audioPlayer.getPlaylistLength()
    }
}
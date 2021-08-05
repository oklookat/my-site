import ElvenPlayerCore from "@/common/plugins/ElvenPlayer/core/ElvenPlayerCore"

export default class ElvenPlayerControls {

    audioPlayer = undefined

    constructor() {
        this.audioPlayer = new ElvenPlayerCore()
    }

    addSource(source) {
        this.audioPlayer.addSource(source)
    }

    playPause() {
        console.log('PLAY/PAUSE')
        if(!this.audioPlayer){
             return null
         }
        if (!this.audioPlayer.isPlaying) {
             this.audioPlayer.play()
        } else {
             this.audioPlayer.pause()
        }
    }

    next() {
        this.audioPlayer.next()
    }

    prev() {
        this.audioPlayer.previous()
    }

    setDurationOnClick(event){
        this.audioPlayer.setDurationOnClick(event)
    }
}
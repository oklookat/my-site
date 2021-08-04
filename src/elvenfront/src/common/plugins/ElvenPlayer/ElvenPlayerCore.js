import {Howl, Howler} from 'howler'

export default class ElvenPlayerCore {

    audio = undefined
    isPlaying = false

    constructor(url) {
        console.log(url)
        this.audio = new Howl({src: [url]})
        console.log(this.audio)
    }

    playPause(){
        if(this.isPlaying){
            this.audio.play()
        } else {
            this.audio.pause()
        }
        this.isPlaying = !this.isPlaying
    }
}
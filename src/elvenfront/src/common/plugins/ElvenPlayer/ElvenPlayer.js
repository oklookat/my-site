import ElvenPlayerC from './ElvenPlayerC.vue'
import ElvenPlayerCore from "@/common/plugins/ElvenPlayer/core/ElvenPlayerCore";

export default class ElvenPlayer {
    static componentData = null
    static audioPlayer = null

    static install(app, options) {
        app.component('elven-player', ElvenPlayerC)
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const instance = this
        app.mixin({
            created() {
                if (this.SERVICE === 'ELVEN_PLAYER_C') {
                    theLogic.options = options
                    theLogic.init()
                    instance.componentData = this
                    app.config.globalProperties.$elvenPlayer = theLogic
                }
            },
            mounted() {
                if (this.SERVICE === 'ELVEN_PLAYER_C') {
                    instance.audioPlayer = new ElvenPlayerCore()
                }
            },
        })
    }
}

export class theLogic {

    static options = null

    static init() {
        if (this.options) {
            return null
        }
    }

    static addToPlaylist(url) {
        ElvenPlayer.audioPlayer.addToPlaylist(url)
    }

    static setPlaylist(playlist){
        ElvenPlayer.audioPlayer.setPlaylist(playlist)
    }

    static play(url){
        theLogic.setPlaylist([url])
        ElvenPlayer.audioPlayer.play()
    }
}
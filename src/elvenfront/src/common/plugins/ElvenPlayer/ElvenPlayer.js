// Elven Player (Vue 3)
// https://github.com/oklookat

import ElvenPlayerC from './ElvenPlayerC'
import {methods} from './core/Shared'

export default class ElvenPlayer {
    static componentData = null

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
        })
    }
}

class theLogic {

    static options = null

    static init() {
        if (this.options) {
            return null
        }
    }

    static addToPlaylist(url){
        methods.addToPlaylist(url)
        const playlistLength = methods.getPlaylistLength()
        if(playlistLength > 0){
            ElvenPlayer.componentData.showFunc(true)
        }
    }

    static playSingleAudio(url){
        methods.setPlaylist(url)
        methods.play()
        ElvenPlayer.componentData.showFunc(true)
    }
}
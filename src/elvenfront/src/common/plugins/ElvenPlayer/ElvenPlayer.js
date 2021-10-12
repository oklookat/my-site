import ElvenPlayerC from './components/Player/Component.vue'

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
                    window.$elvenPlayer = theLogic
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

    static async addToPlaylist(url) {
        await ElvenPlayer.componentData.player.addToPlaylist(url)
    }

    static async setPlaylist(playlist){
        await ElvenPlayer.componentData.player.setPlaylist(playlist)
    }

    static async play(url){
        ElvenPlayer.componentData.player.audioPlayer.active = true
        await theLogic.setPlaylist([url])
        await ElvenPlayer.componentData.player.play()
    }
}
// Elven Player (Vue 3 and Howler)
// https://github.com/oklookat

import ElvenPlayerC from './ElvenPlayerC'

export default class ElvenPlayer {
    static componentData = null

    static install(app, options) {
        app.component('elven-player', ElvenPlayerC)
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const instance = this
        app.mixin({
            created() {
                if (this.SERVICE === 'ELVEN_PLAYER_C') {
                    instance.componentData = this
                    theLogic.options = options
                   // theLogic.init()
                    app.config.globalProperties.$elvenPlayer = theLogic
                }
            },
        })
    }
}

class theLogic {

    static options = null

    // static init() {
    //     if (this.options) {
    //         return null
    //     }
    // }

    static setSources(sources) {
        ElvenPlayer.componentData.sources = sources
        return true
    }

    static addSource(source){
        ElvenPlayer.componentData.addSource = source
        return true
    }
}
// Elven Player (Vue 3 and Howler)
// https://github.com/oklookat

import ElvenPlayerC from './ElvenPlayerC'
import {methods, state} from './core/Shared'

export default class ElvenPlayer {

    static install(app, options) {
        app.component('elven-player', ElvenPlayerC)
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        app.mixin({
            created() {
                if (state.SERVICE === 'ELVEN_PLAYER_C') {
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

    static addSource(source){
        methods.addSource(source)
        return true
    }
}
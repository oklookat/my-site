import ElvenNotifyC from './ElvenNotifyC.vue'

/*
Options:
timer = time in milliseconds when notification deleted
*/

export default class ElvenNotify {

    static componentData = null

    static install(app, options) {
        app.component('elven-notify', ElvenNotifyC)
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const _this = this
        app.mixin({
            created() {
                if (this.SERVICE === 'ELVEN_NOTIFY_C') {
                    _this.componentData = this
                    theLogic.options = options
                    theLogic.init()
                    window.$elvenNotify = theLogic
                }
            },
        })
    }
}

export class theLogic {

    static options = null

    static init() {
        if (this.options) {
            if (this.options.timer) {
               ElvenNotify.componentData.deletedIn = this.options.timer
            }
        }
    }

    static add(message){
      ElvenNotify.componentData.addNotification(message)
    }
}
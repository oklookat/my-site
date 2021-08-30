import ElvenNotifyC from './ElvenNotifyC.vue'

export default class ElvenNotify {
    static componentData = null

    static install(app, options) {
        app.component('elven-notify', ElvenNotifyC)
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const instance = this
        app.mixin({
            created() {
                if (this.SERVICE === 'ELVEN_NOTIFY_C') {
                    theLogic.options = options
                    theLogic.init()
                    instance.componentData = this
                    app.config.globalProperties.$elvenNotify = theLogic
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

    static error(message){
      // ElvenNotify.componentData.addError(message)
        let cnt = 0
        const inter = setInterval(() =>{
            if(cnt >= 10){
                clearInterval(inter)
            } else {
                cnt++
                ElvenNotify.componentData.addError(`сообщение: ${cnt}`)
            }
        }, 500)
    }

}
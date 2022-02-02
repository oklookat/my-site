import './style.scss'
import type { 
    ElvenNotify as IElvenNotify, 
    Notification as TNotification
} from "./types";
import Notification from './notification';

export class ElvenNotify implements IElvenNotify {

    private container: HTMLDivElement
    // used for set notification id
    private counter = 0
    private settings = {
        // when notification will be deleted (in ms)
        deletedIn: 5000,
        // max notifications on desktop
        maxNotificationsDesktop: 8,
        // max notifications on mobile
        maxNotificationsMobile: 2
    }

    constructor(container: HTMLDivElement) {
        this.container = container
        window.$notify = this
    }
 
    public add(n: TNotification) {
        const count = this.container.childNodes.length
        // clear counter if no notifications
        if (count < 1) {
            this.counter = 0
        }
        // adaptive
        let isMaxNotifications = false
        if (window.screen.width > 765) {
            isMaxNotifications = count > this.settings.maxNotificationsDesktop - 1
        } else {
            isMaxNotifications = count > this.settings.maxNotificationsMobile - 1
        }
        if (isMaxNotifications) {
            const first = this.container.lastElementChild
            const isDiv = first && first['tagName'] && first.tagName === 'DIV'
            if (isDiv) {
                first.remove()
            }
        }
        this.set(n);
    }

    /** create full notification and render */
    private set(n: TNotification) {
        const not = new Notification(this.container, this.counter, n.message)
        not.render()
        this.counter++
    }

}
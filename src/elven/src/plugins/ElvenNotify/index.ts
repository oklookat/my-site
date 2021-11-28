import {
    ClassName,
    ElvenNotify as IElvenNotify,
    Notification,
    NotificationFull,
} from "./types";
import './style.scss'

// TODO: add transition
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
        window.$elvenNotify = this
    }

    public add(n: Notification) {
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
            if (first && first instanceof HTMLDivElement) {
                this.delete(first)
            }
        }
        this.set(n);
    }

    /** create full notification and render */
    private set(notification: Notification) {
        const full: NotificationFull = {
            id: this.counter++,
            self: notification,
            percents: 0,
            timeWhenGone: 0,
            executed: false,
            timeoutID: null,
            intervalID: null
        };
        this.render(full)
    }

    private render(n: NotificationFull) {
        const container = document.createElement('div')
        container.className = ClassName.notification
        container.onclick = () => this.delete(container)
        const message = document.createElement('div')
        message.className = ClassName.message
        message.innerText = n.self.message
        container.appendChild(message)
        //
        const timerWrapper = document.createElement('div')
        timerWrapper.className = ClassName.timerWrapper
        const timer = document.createElement('div')
        timer.className = ClassName.timer
        timer.style.transform = `scaleX(${n.percents / 100})`
        timerWrapper.appendChild(timer)
        //
        container.appendChild(timerWrapper)
        this.execute(n, container, timer)
        this.container.appendChild(container)
    }

    private delete(node: HTMLDivElement) {
        node.remove()
    }

    private calcPercents(n: NotificationFull, deletedIn: number) {
        const now = new Date().getTime();
        // if date when item should be deleted
        if (now >= n.timeWhenGone) {
            clearInterval(n.intervalID);
        }
        // get the difference between current date and time when item deleted
        const diff = Math.abs(now - n.timeWhenGone);
        // get how much is left as a percentage. deletedIn = 100%
        n.percents = (diff / deletedIn) * 100;
    }

    /** init timers */
    private execute(n: NotificationFull, container: HTMLDivElement, timer: HTMLDivElement) {
        if (n.executed) {
            return n;
        }
        n.timeoutID = setTimeout(() => {
            // delete from array after time
            this.delete(container);
        }, this.settings.deletedIn);
        // calc time when notification be deleted
        const deletedIn = this.settings.deletedIn
        // set time when item deleted
        n.timeWhenGone = new Date().getTime() + deletedIn;
        n.intervalID = setInterval(() => {
            this.calcPercents(n, deletedIn);
            timer.style.transform = `scaleX(${n.percents / 100})`
        }, 100); // interval time = performance. timer transition time = this time + 20ms
        n.executed = true;
    }
}
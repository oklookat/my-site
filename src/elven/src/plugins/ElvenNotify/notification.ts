import { ClassName } from "./types"
import Animation from "../../tools/animation"

/** represents one notification */
export default class Notification {

    public id: number = 0
    /** user message */
    public message: string
    /** setting: when should be deleted (in ms) */
    private deletedIn: number = 5000
    /** elements */
    private el = {
        /** user container */
        root: HTMLDivElement = undefined,
        /** notification container */
        container: {
            /** container itself */
            self: HTMLDivElement = undefined,
            /** user message */
            message: HTMLDivElement = undefined,
            /** progress */
            progress: {
                /** wrapper */
                wrapper: HTMLDivElement = undefined,
                /** line */
                self: HTMLDivElement = undefined
            }
        }
    }
    /** timestamp when deleted */
    public timestampWhenDeleted: number = 0
    /** when deleted (in percents) */
    public progress: number = 100
    /** update percents */
    private percentsInterval: ReturnType<typeof setInterval>
    /** is timers active */
    public isExecuted: boolean = false


    constructor(root: HTMLDivElement, id: number, message: string, time: number = 5000) {
        this.deletedIn = time
        this.el.root = root
        this.id = id
        this.message = message
    }

    /** run timers, render notification */
    public render() {
        // create container
        const container = document.createElement('div')
        container.style.opacity = '0'
        container.className = ClassName.notification
        container.onclick = () => this.delete()
        // create message
        const message = document.createElement('div')
        message.className = ClassName.message
        message.innerText = this.message
        // create progress wrapper
        const progressWrapper = document.createElement('div')
        progressWrapper.className = ClassName.progressWrapper
        // create progress
        const progress = document.createElement('div')
        progress.className = ClassName.progress
        // append
        container.appendChild(message)
        container.appendChild(progressWrapper)
        progressWrapper.appendChild(progress)
        // assign
        this.el.container.self = container
        this.el.container.message = message
        this.el.container.progress.wrapper = progressWrapper
        this.el.container.progress.self = progress
        // final
        this.run()
        this.el.root.appendChild(this.el.container.self)
        Animation.fadeIn(this.el.container.self)
    }

    private delete() {
        Animation.fadeOut(this.el.container.self)
            .then(() => {
                this.el.container.self.remove()
            })
    }

    /** set up hooks */
    private run() {
        if (this.isExecuted) {
            return
        }
        this.isExecuted = true
        // set time when item deleted
        this.timestampWhenDeleted = this.getTimestamp() + this.deletedIn
        // when deleting
        setTimeout(() => {
            this.delete()
        }, this.deletedIn)
        // when shown
        this.percentsInterval = setInterval(() => {
            this.computeProgress()
        }, 100)
    }

    /** sets progress time depend on deletion timestamp */
    private computeProgress() {
        const now = this.getTimestamp()
        // if should be deleted
        if (now >= this.timestampWhenDeleted) {
            clearInterval(this.percentsInterval)
            return
        }
        // get the difference between current date and time when item deleted
        const diff = Math.abs(now - this.timestampWhenDeleted);
        // get how much is left as a percentage. deletedIn = 100%
        this.progress = (diff / this.deletedIn) * 100;
        this.el.container.progress.self.style.transform = `scaleX(${this.progress / 100})`
    }

    /** get current timestamp */
    private getTimestamp(): number {
        return new Date().getTime()
    }

}
/** notification element class names */
export enum ClassName {
    notification = 'notification',
    message = 'notification__message',
    progressWrapper = 'notification__progress-wrapper',
    progress = 'notification__progress'
}

/** plugin */
export interface ElvenNotify {
    add: (notification: Notification) => void
}

/** user notification */
export type Notification = {
    message: string;
}
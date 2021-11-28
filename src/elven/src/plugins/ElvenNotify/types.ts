/** notification element class names */
export enum ClassName {
    notification = 'notification',
    message = 'notification__message',
    timerWrapper = 'notification__timer-wrapper',
    timer = 'notification__timer'
}

/** plugin */
export interface ElvenNotify {
    add: (notification: Notification) => void
}

/** user notification */
export type Notification = {
    message: string;
}

/** represents full notification (used in internal) */
export type NotificationFull = {
    id: number;
    /** user notification object */
    self: Notification;
    /** delete notification in percents */
    percents: number;
    /** when must be deleted (unix timestamp; used for percents calc) */
    timeWhenGone: number;
    /** is notification already executed */
    executed: boolean;
    /** delete when timeout ends */
    timeoutID: NodeJS.Timeout | null;
    /** calculate percents with interval */
    intervalID: ReturnType<typeof setInterval> | null;
}
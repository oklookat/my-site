/** user notification */
export type TNotification = {
    message: string;
}

/** represents full notification (used in internal) */
export type TNotificationFull = {
    id: number;
    /** user notification object */
    self: TNotification;
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

/** notification plugin */
export interface IElvenNotify {
    add: (notification: TNotification) => void
}
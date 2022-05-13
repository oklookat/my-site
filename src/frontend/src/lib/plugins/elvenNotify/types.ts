/** notification element class names */
export enum ClassName {
	notification = 'notification',
	message = 'message'
}

/** plugin */
export interface ElvenNotify {
	add: (notification: Notification) => void;
	destroy(): void;
}

/** user notification */
export type Notification = {
	message: string;
};

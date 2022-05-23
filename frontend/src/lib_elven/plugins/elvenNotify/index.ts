import './style.scss';
import type { ElvenNotify as IElvenNotify, Notification as TNotification } from './types';
import Notification from './notification';

export class ElvenNotify implements IElvenNotify {
	private container: HTMLDivElement;

	private uid = 0;

	private settings = {
		/** when notification will be deleted (in ms) */
		deletedIn: 5000,
		/** max notifications on desktop */
		maxNotificationsDesktop: 8,
		/** max notifications on mobile */
		maxNotificationsMobile: 2
	};

	constructor(container: HTMLDivElement, time: number = 5000) {
		this.settings.deletedIn = time;
		this.container = container;
		window.$notify = this;
	}

	public destroy(): void {
		window.$notify = undefined;
	}

	public add(n: TNotification) {
		const count = this.container.childNodes.length;

		// clear counter if no notifications
		if (count < 1) {
			this.uid = 0;
		}

		// get max notifications depend on screen width
		let isMaxNotifications = false;
		if (window.screen.width > 765) {
			isMaxNotifications = count > this.settings.maxNotificationsDesktop - 1;
		} else {
			isMaxNotifications = count > this.settings.maxNotificationsMobile - 1;
		}

		// remove oldest notification if max notifications
		if (isMaxNotifications) {
			const first = this.container.lastElementChild;
			const isDiv = first instanceof HTMLDivElement;
			if (isDiv) {
				first.remove();
			}
		}

		this.set(n);
	}

	/** create notification and render */
	private set(n: TNotification) {
		const not = new Notification(this.container, this.uid, n.message, this.settings.deletedIn);
		not.render();
		this.uid++;
	}
}

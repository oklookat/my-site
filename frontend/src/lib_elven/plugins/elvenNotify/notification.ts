import { fadeIn, fadeOut } from '$elven/tools';
import { ClassName } from './types';

/** represents one notification */
export default class Notification {
	public id: number = 0;

	/** user message */
	public message: string;

	/** setting: when should be deleted (in ms) */
	private deletedIn: number = 5000;

	/** user container */
	private rootEL: HTMLDivElement;

	/** notification container */
	private containerEL: HTMLDivElement;

	/** when deleted (in percents) */
	public progress: number = 100;

	/** is timers active */
	public isExecuted: boolean = false;

	private deleteTimeout: NodeJS.Timeout | undefined;

	constructor(root: HTMLDivElement, id: number, message: string, time: number = 5000) {
		this.deletedIn = time;
		this.rootEL = root;
		this.id = id;
		this.message = message;

		// create container
		const container = document.createElement('div');
		container.style.opacity = '0';
		container.className = ClassName.notification;
		container.onclick = () => this.delete();

		// create message
		const msg = document.createElement('div');
		msg.className = ClassName.message;
		msg.innerText = this.message;
		container.appendChild(msg);

		// assign
		this.containerEL = container;
	}

	/** run timers, render notification */
	public render() {
		this.rootEL.appendChild(this.containerEL);
		fadeIn(this.containerEL).then(() => {
			this.runDeleteTimer();
		});
	}

	private delete() {
		if (this.deleteTimeout) {
			clearTimeout(this.deleteTimeout);
		}

		fadeOut(this.containerEL).then(() => {
			this.containerEL.remove();
		});
	}

	private runDeleteTimer() {
		if (this.isExecuted) {
			return;
		}

		this.isExecuted = true;

		// when deleting
		this.deleteTimeout = setTimeout(() => {
			this.delete();
		}, this.deletedIn);
	}
}

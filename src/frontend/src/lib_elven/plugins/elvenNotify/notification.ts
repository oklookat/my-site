import { ClassName } from './types';
import Animation from '../../tools/animation';

/** represents one notification */
export default class Notification {
	public id: number = 0;
	/** user message */
	public message: string;
	/** setting: when should be deleted (in ms) */
	private deletedIn: number = 5000;
	/** elements */
	private el = {
		/** user container */
		root: (HTMLDivElement = undefined),
		/** notification container */
		container: {
			/** container itself */
			self: (HTMLDivElement = undefined),
			/** user message */
			message: (HTMLDivElement = undefined)
		}
	};
	/** timestamp when deleted */
	public timestampWhenDeleted: number = 0;
	/** when deleted (in percents) */
	public progress: number = 100;
	/** is timers active */
	public isExecuted: boolean = false;

	constructor(root: HTMLDivElement, id: number, message: string, time: number = 5000) {
		this.deletedIn = time;
		this.el.root = root;
		this.id = id;
		this.message = message;
	}

	/** run timers, render notification */
	public render() {
		// create container
		const container = document.createElement('div');
		container.style.opacity = '0';
		container.className = ClassName.notification;
		container.onclick = () => this.delete();
		// create message
		const message = document.createElement('div');
		message.className = ClassName.message;
		message.innerText = this.message;
		// create progress wrapper
		const progressWrapper = document.createElement('div');
		progressWrapper.className = ClassName.progressWrapper;
		// append
		container.appendChild(message);
		container.appendChild(progressWrapper);
		// assign
		this.el.container.self = container;
		this.el.container.message = message;
		// final
		this.run();
		this.el.root.appendChild(this.el.container.self);
		Animation.fadeIn(this.el.container.self);
	}

	private delete() {
		Animation.fadeOut(this.el.container.self).then(() => {
			this.el.container.self.remove();
		});
	}

	/** set up hooks */
	private run() {
		if (this.isExecuted) {
			return;
		}
		this.isExecuted = true;
		// set time when item deleted
		this.timestampWhenDeleted = this.getTimestamp() + this.deletedIn;
		// when deleting
		setTimeout(() => {
			this.delete();
		}, this.deletedIn);
	}

	/** get current timestamp */
	private getTimestamp(): number {
		return new Date().getTime();
	}
}

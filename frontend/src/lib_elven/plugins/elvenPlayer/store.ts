import type Audio from './audio.svelte';
import {
	writable,
	type Subscriber,
	type Unsubscriber,
	type Updater,
	type Writable
} from 'svelte/store';
import type { Signal } from './types';

class SensitiveStore<T> implements Writable<T> {
	private uid = 0;
	private freeUids: number[] = [];
	private value: T;
	private subs: Record<number, Subscriber<T>> = {};

	constructor(value: T) {
		this.value = value;
	}

	public set(value: T): void {
		this.value = value;
		this.informSubs();
	}

	public update(updater: Updater<T>): void {
		this.value = updater(this.value);
		this.informSubs();
	}

	public subscribe(run: Subscriber<T>, invalidate?: (value?: T) => void): Unsubscriber {
		let newUid = 0;
		if (this.freeUids.length > 0) {
			newUid = this.freeUids[0];
			this.freeUids = this.freeUids.slice(1);
		} else {
			newUid = this.uid++;
		}

		this.subs[newUid] = run;
		return () => {
			delete this.subs[newUid];
			this.freeUids.push(newUid);
			if (invalidate) {
				invalidate(this.value);
			}
		};
	}

	private informSubs() {
		for (const key in this.subs) {
			const sub = this.subs[key];
			sub(this.value);
		}
	}
}

/** current playing source */
export const currentSource = writable<URL | undefined>(undefined);

/** current audio component */
export const audioComponent = writable<Audio | undefined>(undefined);

/** is playing now? */
export const isPlaying = writable(false);

/** volume in percents */
export const volumePercents = writable(100);

/** volume like: 1.0 */
export const volumeFloat = writable(1.0);

/** set current time */
export const setCurrentTime = new SensitiveStore(0);

/** current time like: 4921.001 */
export const currentTime = new SensitiveStore(0);

/** current time like: 03:21 */
export const currentTimePretty = writable('00:00');

/** current time in percents */
export const currentTimePercents = writable(0);

/** duration like: 0 */
export const duration = writable(0);

/** duration like: 03:21 */
export const durationPretty = writable('00:00');

/** buffered like: 69 */
export const bufferedPercents = writable(0);

/** signal to Audio component */
export const signal = new SensitiveStore<Signal>(0);

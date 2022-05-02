import type Store from './store';
import { convertCurrentTimePretty, getBufferedPercents, getPercents, getPretty } from '$lib_elven/plugins/elvenPlayer/core/utils';
import type { Events as IEvents} from '../types';
import Logger from './logger';

/** updates playback state */
export default class Events implements IEvents {
	private store: Store;

	constructor(store: Store) {
		this.store = store;
	}

	public onPlaying() {
		this.store.playing = true;
	}

	public onPause() {
		this.store.playing = false;
	}

	public onEnded() {
		this.store.playing = false;
		this.store.ended = true;
		this.store.ended = false;
	}

	public onTimeUpdate(e?: Event) {
		if(!e) {
			return
		}
		const el = e.target as HTMLAudioElement;
		const currentTime = el.currentTime;
		const buffered = el.buffered;
		let duration = el.duration;
		const badDuration = !duration || duration === Infinity;
		if (badDuration) {
			duration = 0;
		}
		this.store.bufferedPercents = getBufferedPercents(currentTime, duration, buffered);
		this.store.durationNum = duration;
		this.store.durationPretty = getPretty(duration);
		this.store.currentTimeNum = currentTime;
		this.store.currentTimePretty = convertCurrentTimePretty(currentTime, duration);
		this.store.currentTimePercents = getPercents(currentTime, duration);
	}

	public onError(e?: Event) {
		this.onPause();
		if (!e || !('target' in e)) {
			Logger.error(`unknown error.`);
			return
		}

		const target = e.target as HTMLMediaElement;
		const err = target.error;
		if (!err) {
			Logger.error(`unknown error.`);
			return
		}

		const msg = err.message ? ` ${err.message}` : '';

		// https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/networkState
		switch (err.code) {
			case MediaError.MEDIA_ERR_ABORTED:
				Logger.error(`aborted.${msg}`);
				break;
			case MediaError.MEDIA_ERR_NETWORK:
				Logger.error(`network error.${msg}`);
				break;
			case MediaError.MEDIA_ERR_DECODE:
				Logger.error(`decode error. Maybe audio damaged or something?${msg}`);
				break;
			case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
				Logger.error(`not supported.${msg}`);
				if (!msg) {
					Logger.error(
						`it is .flac? If yes, maybe in tags of your .flac exists non-ASCII chars? Idk why (browser problem maybe?), but we cannot play flacs with non-ASCII title/tags/etc.\nrelated: https://github.com/koel/koel/issues/869`
					);
				}
				break;
			default:
				Logger.error(`unknown error.${msg}`);
				break;
		}
	}
}

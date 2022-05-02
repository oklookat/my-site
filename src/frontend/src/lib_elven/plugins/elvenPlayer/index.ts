import type {
	Events as IEvents,
	Playlist,
	Source,
	ElvenPlayer as IElvenPlayer,
	Unsubscribe
} from './types';
import DOM from './core/dom';
import Store from './core/store';
import Events from './core/events';
import Logger from './core/logger';
import { getPercents, convertPercentsToCurrentTime as _convertPercentsToCurrentTime, convertCurrentTimePretty } from '$lib_elven/plugins/elvenPlayer/core/utils';

// warning: Chrome has bug with FLAC double-rewind to the end (PTS not defined). Idk how fix that, except decrease time of rewinding to end
/** controls audio player */
export default class ElvenPlayer implements IElvenPlayer {
	public store: Store;
	private dom: DOM;
	private unsubs: Unsubscribe[] = [];
	private events: IEvents;

	private _playlist: Playlist = {
		position: 0,
		sources: []
	};

	constructor() {
		this.store = new Store();
		this.events = new Events(this.store);
		this.dom = new DOM(this.events);
		this.subscribe();
		window.$player = this;
	}

	public destroy() {
		this.unsubscribe();
		this.dom.destroy();
		window.$player = undefined;
	}

	private subscribe() {
		this.unsubs.push(this.store.state.current.ended.onChange((v) => {
			if (v && this.store.playing) {
				this.next();
			}
		}));
	}

	private unsubscribe() {
		for (const unsub of this.unsubs) {
			unsub();
		}
		this.unsubs = [];
	}

	public addToPlaylist(source: Source) {
		const isPlaylistEmpty = this.playlist.sources.length < 1
		if (isPlaylistEmpty) {
			this.playlist.position = 0
		}

		this.playlist.sources.push(source);

		if (isPlaylistEmpty) {
			this.recheckSource(true)
		}
	}

	public clearPlaylist() {
		this.playlist = {
			position: 0,
			sources: []
		};
	}

	private recheckSource(force = false) {
		if (!this.dom.source || force) {
			this.stop()
			this.dom.source = this.playlist.sources[this.playlist.position]
		}
	}

	public async playPause() {
		this.recheckSource()

		if (this.isPlaying) {
			this.dom.pause()
			return
		}

		try {
			// @ts-ignore
			await this.dom.play();
		} catch (err) {
			Logger.error(`${err}`);
		}
	}

	public stop() {
		this.dom.stop();
	}

	private repeat() {
		this.stop();
		this.playPause();
	}

	public next() {
		const next = this.playlist.sources[this.playlist.position + 1];
		if (!next) {
			// stop if no source next
			this.stop();
			return;
		}
		this.playlist.position++;
		this.dom.source = this.playlist.sources[this.playlist.position];

		this.playPause();
	}

	public prev() {
		// if no source behind
		const prev = this.playlist.sources[this.playlist.position - 1];

		// if current time > 2% of total time - repeat
		const notInStart = this.currentTimePercents > 2;
		if (!prev || notInStart) {
			this.repeat();
			return;
		}

		this.playlist.position--;
		this.dom.source = this.playlist.sources[this.playlist.position];

		this.playPause();
	}

	public get isPlaying(): boolean {
		return this.store.playing;
	}

	public get playlist(): Playlist {
		return this._playlist;
	}

	public set playlist(playlist: Playlist) {
		const sources = playlist.sources;
		if (!sources || !(sources instanceof Array)) {
			Logger.error('wrong playlist');
			return;
		}

		this._playlist = playlist;

		let position = 0;
		if (this._playlist.position && sources[this._playlist.position]) {
			position = playlist.position;
		}

		this._playlist.position = position;
	}

	// utils

	public get source(): string {
		return this.dom.source;
	}

	public set volume(volume: number) {
		this.dom.volume = volume;
	}

	public set volumePercents(percents: number) {
		this.dom.volume = percents / 100;
		this.store.volumePercents = percents;
	}

	public get volume(): number {
		return this.dom.volume;
	}

	public get volumePercents(): number {
		return this.dom.volume * 100;
	}

	public set currentTime(seconds: number) {
		this.dom.currentTime = seconds;
	}

	public set currentTimePercents(percents: number) {
		const seconds = this.convertPercentsToCurrentTime(percents);
		this.dom.currentTime = seconds;
	}

	public get currentTime(): number {
		return this.dom.currentTime;
	}

	public get currentTimePercents(): number {
		return getPercents(this.dom.currentTime, this.dom.duration);
	}

	public get duration(): number {
		return this.dom.duration;
	}

	public convertPercentsToCurrentTime(percents: number): number {
		return _convertPercentsToCurrentTime(percents, this.dom.duration);
	}

	public convertPercentsToCurrentTimePretty(percents: number): string {
		const seconds = this.convertPercentsToCurrentTime(percents);
		return convertCurrentTimePretty(seconds, this.dom.duration);
	}
}

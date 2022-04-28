import type {
	Events as IEvents,
	Store as IStore,
	Playlist,
	Source,
	ElvenPlayer as IElvenPlayer,
	Unsubscribe
} from './types';
import DOM from './core/dom';
import Store from './core/store';
import Events from './core/events';
import Logger from './core/logger';
import Utils from './core/utils';

// warning: Chrome has bug with FLAC double-rewind to the end (PTS not defined). Idk how fix that, except decrease time of rewinding to end
/** controls audio player */
export default class ElvenPlayer implements IElvenPlayer {
	public store: IStore;
	private dom: DOM;
	private initialized: boolean = false;
	private unsubs: Unsubscribe[] = [];
	private events: IEvents;

	private _playlist: Playlist = {
		position: 0,
		sources: []
	};

	constructor() {
		this.init();
	}

	public init() {
		this.store = new Store();
		this.events = new Events(this.store);
		this.dom = new DOM(this.events);
		this.subscribe();
		this.initialized = true;
		window.$player = this;
	}

	public destroy() {
		this.unsubscribe();
		this.dom.destroy();
		this.dom = null;
		this.initialized = false;
		window.$player = undefined;
	}

	private subscribe() {
		const u1 = this.store.state.current.ended.onChange((v) => {
			if (v && this.store.playing) {
				this.next();
			}
		});
		this.unsubs.push(u1);
	}

	private unsubscribe() {
		for (const unsub of this.unsubs) {
			unsub();
		}
		this.unsubs = [];
	}

	public addToPlaylist(source: Source) {
		this.playlist.sources.push(source);
	}

	public clearPlaylist() {
		this.playlist = {
			position: 0,
			sources: []
		};
	}

	public async play() {
		this.initIfNotInit();
		try {
			const currentSource = this._playlist.sources[this._playlist.position];
			await this.dom.play(currentSource);
		} catch (err) {
			Logger.error(err);
		}
	}

	public pause() {
		this.initIfNotInit();
		this.dom.pause();
	}

	public stop() {
		this.initIfNotInit();
		this.dom.stop();
	}

	private repeat() {
		this.stop();
		this.play();
	}

	public next() {
		const next = this.playlist.sources[this.playlist.position + 1];
		if (!next) {
			// stop if no source next
			this.stop();
			return;
		}
		this.playlist.position++;
		this.play();
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
		this.play();
	}

	private initIfNotInit() {
		if (this.initialized) {
			return;
		}
		this.init();
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
		return Utils.getPercents(this.dom.currentTime, this.dom.duration);
	}

	public get duration(): number {
		return this.dom.duration;
	}

	public convertPercentsToCurrentTime(percents: number): number {
		return Utils.convertPercentsToCurrentTime(percents, this.dom.duration);
	}

	public convertPercentsToCurrentTimePretty(percents: number): string {
		const seconds = this.convertPercentsToCurrentTime(percents);
		return Utils.convertCurrentTimePretty(seconds, this.dom.duration);
	}
}

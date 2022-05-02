import type Store from "./core/store";

export type ConvertSecondsMode = 'auto' | 'hours' | 'minutes';
export type Source = string;
export type Sources = Array<string>;
/** playlist */
export type Playlist = {
	/** current index of source */
	position: number;

	/** list of playable (mostly) audios */
	sources: Sources;
};

/** unsubscribe from hook */
export interface Unsubscribe {
	(): void;
}

/** state management */
export interface State<T> {
	/** get value */
	get(): T;

	/** set value */
	set(value?: T): void;

	/** hook when value changed. Returns unsubscribe function */
	onChange(hook: (value?: T) => void): Unsubscribe;
}

/** local copy of IStore.state (almost),
 * thats updating by subscription or by external changes.
 * Used in component */
export type ComponentState = {
	playing: boolean;
	volume: {
		percents: number;
	};
	current: {
		buffered: {
			percents: number;
		};
		time: {
			/** is user dragging time slider */
			draggingNow: boolean;
			percents: number;
			pretty: string;
		};
		duration: {
			pretty: string;
		};
	};
};

/** plugin */
export interface ElvenPlayer {
	/** player store / state */
	store: Store;

	/** is playing? */
	get isPlaying(): boolean;

	/** set playlist */
	set playlist(playlist: Playlist);

	/** get playlist */
	get playlist(): Playlist;

	/** get audio source */
	get source(): string;

	/** set audio source */
	set source(src: string);

	/** set volume (0 - 1) */
	set volume(volume: number);

	/** set volume in percents (0 - 100) */
	set volumePercents(percents: number);

	/** get volume (0 - 1) */
	get volume(): number;

	/** get volume in percents (0 - 100) */
	get volumePercents(): number;

	/** set current time of audio in seconds */
	set currentTime(seconds: number);

	/** set current time by percents. Where 100 - audio ends */
	set currentTimePercents(percents: number);

	/** get current time of audio in seconds */
	get currentTime(): number;

	/** get current time percents. Where 100 - audio ends */
	get currentTimePercents(): number;

	/** get audio duration in seconds */
	get duration(): number;

	/** convert current time percents to seconds depending on duration */
	convertPercentsToCurrentTime(percents: number): number;

	/** convert current time percents to string like '01:11' or '01:11:11' */
	convertPercentsToCurrentTimePretty(percents: number): string;

	/** play audio */
	playPause(): void;

	/** stop audio */
	stop(): void;

	/** next audio */
	next(): void;

	/** previous audio */
	prev(): void;

	/** add source to playlist */
	addToPlaylist: (source: Source) => void;

	/** remove all tracks from playlist */
	clearPlaylist(): void;
}

/** audio element events */
export interface Events {
	/** when audio playing */
	onPlaying: (e?: Event) => void;

	/** when audio paused */
	onPause: (e?: Event) => void;

	/** when audio ended */
	onEnded: (e?: Event) => void;

	/** when audio time updated */
	onTimeUpdate: (e?: Event) => void;

	/** when error */
	onError: (e?: Event) => void;
}

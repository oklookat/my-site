import State from './state';

/** manage audio player store */
export default class Store {
	public state = {
		/** is audio playing now? */
		playing: new State(false),

		/** in float */
		volume: {
			/** in float */
			num: new State(1.0),

			/** in percents */
			percents: new State(100)
		},

		/** current playing */
		current: {
			/** is audio ended */
			ended: new State(false),

			duration: {
				/** in seconds */
				num: new State(0),

				/** in string like '04:20' */
				pretty: new State('00:00')
			},
			time: {
				/** in seconds */
				num: new State(0),

				/** in string '01:37' */
				pretty: new State('00:00'),

				/** in percents */
				percents: new State(0)
			},

			buffered: {
				/** in percents */
				percents: new State(0)
			}
		}
	};

	public set playing(v: boolean) {
		this.state.playing.set(v);
	}

	public get playing(): boolean {
		return this.state.playing.get();
	}

	public set ended(v: boolean) {
		this.state.current.ended.set(v);
	}

	public get ended(): boolean {
		return this.state.current.ended.get();
	}

	public set bufferedPercents(v: number) {
		this.state.current.buffered.percents.set(v);
	}

	public get bufferedPercents(): number {
		return this.state.current.buffered.percents.get();
	}

	public set durationNum(v: number) {
		this.state.current.duration.num.set(v);
	}

	public get durationNum(): number {
		return this.state.current.duration.num.get();
	}

	public set durationPretty(v: string) {
		this.state.current.duration.pretty.set(v);
	}

	public get durationPretty(): string {
		return this.state.current.duration.pretty.get();
	}

	public set currentTimeNum(v: number) {
		this.state.current.time.num.set(v);
	}

	public get currentTimeNum(): number {
		return this.state.current.time.num.get();
	}

	public set currentTimePercents(v: number) {
		this.state.current.time.percents.set(v);
	}

	public get currentTimePercents(): number {
		return this.state.current.time.percents.get();
	}

	public set currentTimePretty(v: string) {
		this.state.current.time.pretty.set(v);
	}

	public get currentTimePretty(): string {
		return this.state.current.time.pretty.get();
	}

	public set volumePercents(v: number) {
		this.state.volume.percents.set(v);
	}

	public get volumePercents(): number {
		return this.state.volume.percents.get();
	}

	public set volumeNum(v: number) {
		this.state.volume.num.set(v);
	}

	public get volumeNum(): number {
		return this.state.volume.num.get();
	}
}

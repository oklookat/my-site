import { writable } from 'svelte/store';
import type { IState, TStore } from "../types";

export default class State implements IState {

    public store: TStore = {
        isPlaying: writable(false),
        volume: {
            num: writable(1.0),
            percents: writable(100)
        },
        current: {
            duration: {
                num: writable(0),
                pretty: writable('00:00')
            },
            position: {
                num: writable(0),
                pretty: writable('00:00'),
                percents: writable(0)
            },
            buffered: {
                percents: writable(0)
            }
        }
    }

    public set isPlaying(v: boolean) {
        this.store.isPlaying.set(v)
    }

    public set positionNum(v: number) {
        this.store.current.position.num.set(v)
    }

    public set positionPercents(v: number) {
        this.store.current.position.percents.set(v)
    }

    public set positionPretty(v: string) {
        this.store.current.position.pretty.set(v)
    }

    public set bufferedPercents(v: number) {
        this.store.current.buffered.percents.set(v)
    }

    public set durationNum(v: number) {
        this.store.current.duration.num.set(v)
    }

    public set durationPretty(v: string) {
        this.store.current.duration.pretty.set(v)
    }

    public set volumePercents(v: number) {
        this.store.volume.percents.set(v)
    }

    public set volumeNum(v: number) {
        this.store.volume.num.set(v)
    }
}
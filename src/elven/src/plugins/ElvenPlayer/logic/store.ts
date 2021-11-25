import Statex from '../modules/state';
import type { IStore } from "../types";


/** manage audio player store */
export default class Store implements IStore {

    public state = {
        playing: new Statex(false),
        volume: {
            num: new Statex(1.0),
            percents: new Statex(100)
        },
        current: {
            ended: new Statex(false),
            duration: {
                num: new Statex(0),
                pretty: new Statex('00:00')
            },
            time: {
                num: new Statex(0),
                pretty: new Statex('00:00'),
                percents: new Statex(0)
            },
            buffered: {
                percents: new Statex(0)
            }
        }
    }

    public set playing(v: boolean) {
        this.state.playing.set(v)
    }

    public set ended(v: boolean) {
        this.state.current.ended.set(v)
    }

    public set bufferedPercents(v: number) {
        this.state.current.buffered.percents.set(v)
    }

    public set durationNum(v: number) {
        this.state.current.duration.num.set(v)
    }

    public set durationPretty(v: string) {
        this.state.current.duration.pretty.set(v)
    }

    public set currentTimeNum(v: number) {
        this.state.current.time.num.set(v)
    }

    public set currentTimePercents(v: number) {
        this.state.current.time.percents.set(v)
    }

    public set currentTimePretty(v: string) {
        this.state.current.time.pretty.set(v)
    }

    public set volumePercents(v: number) {
        this.state.volume.percents.set(v)
    }

    public set volumeNum(v: number) {
        this.state.volume.num.set(v)
    }
}
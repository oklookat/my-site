import BasicStore from '../modules/basicStore';
import type { IState, TStore } from "../types";


/** manage audio player store */
export default class State implements IState {

    public store: TStore = {
        playing: new BasicStore(false),
        volume: {
            num: new BasicStore(1.0),
            percents: new BasicStore(100)
        },
        current: {
            ended: new BasicStore(false),
            duration: {
                num: new BasicStore(0),
                pretty: new BasicStore('00:00')
            },
            time: {
                num: new BasicStore(0),
                pretty: new BasicStore('00:00'),
                percents: new BasicStore(0)
            },
            buffered: {
                percents: new BasicStore(0)
            }
        }
    }

    public set playing(v: boolean) {
        this.store.playing.set(v)
    }

    public set ended(v: boolean) {
        this.store.current.ended.set(v)
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

    public set currentTimeNum(v: number) {
        this.store.current.time.num.set(v)
    }

    public set currentTimePercents(v: number) {
        this.store.current.time.percents.set(v)
    }

    public set currentTimePretty(v: string) {
        this.store.current.time.pretty.set(v)
    }

    public set volumePercents(v: number) {
        this.store.volume.percents.set(v)
    }

    public set volumeNum(v: number) {
        this.store.volume.num.set(v)
    }
}
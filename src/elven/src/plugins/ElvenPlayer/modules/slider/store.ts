import type { Store as IStore } from "./types";
import State from '../state';

export default class Store implements IStore {

    public state = {
        isMouseDown: new State(false),
        percents: new State(100),
    }

    public get isMouseDown(): boolean {
        return this.state.isMouseDown.get()
    }

    public set isMouseDown(v: boolean) {
        this.state.isMouseDown.set(v)
    }

    public set percents(v: number) {
        this.state.percents.set(v)
    }

}
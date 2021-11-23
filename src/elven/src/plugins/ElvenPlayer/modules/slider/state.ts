import { writable } from 'svelte/store';
import type { TStore } from "./types";
import type { IState } from "./types";

export default class State implements IState {
    
    public store: TStore = {
        isMouseDown: writable(false),
        percents: writable(100),
    }

    public get isMouseDown(): boolean {
        let val = false
        const unsub = this.store.isMouseDown.subscribe(v => {
            val = v
        })
        unsub()
        return val
    }

    public set isMouseDown(v: boolean) {
        this.store.isMouseDown.set(v)
    }

    public set percents(v: number) {
        this.store.percents.set(v)
    }

}
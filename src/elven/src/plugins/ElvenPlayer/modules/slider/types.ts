import type { Writable } from 'svelte/store'

export type TStore = {
    isMouseDown: Writable<boolean>
    percents: Writable<number>
}

export interface IState {
    store: TStore
    get isMouseDown(): boolean
    set isMouseDown(v: boolean)
    set percents(v: number)
}
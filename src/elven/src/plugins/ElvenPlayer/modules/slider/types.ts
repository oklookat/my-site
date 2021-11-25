export interface IState<T> {
    get(): T
    set(value?: T): void
    /** hook when value changed. Returns unsubscribe function */
    onChange(hook: (value?: T) => void): () => void
}

export interface IStore {
    state: {
        isMouseDown: IState<boolean>
        percents: IState<number>
    }
    get isMouseDown(): boolean
    set isMouseDown(v: boolean)
    set percents(v: number)
}
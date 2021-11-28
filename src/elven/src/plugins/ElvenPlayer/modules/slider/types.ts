export interface State<T> {
    get(): T
    set(value?: T): void
    /** hook when value changed. Returns unsubscribe function */
    onChange(hook: (value?: T) => void): () => void
}

export interface Store {
    state: {
        isMouseDown: State<boolean>
        percents: State<number>
    }
    get isMouseDown(): boolean
    set isMouseDown(v: boolean)
    set percents(v: number)
}
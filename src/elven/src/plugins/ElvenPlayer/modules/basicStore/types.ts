/** available hooks */
export enum THookName {
    onChange = 'onChange'
}

/** remove hook */
export interface IHookRemove {
    (): void
}

/** user hook */
export interface IHook<T> {
    (value?: T): void
}

/** hooks list in specific hook type (internal) */
export type THooksList<T> = {
    /** counter of hooks */
    count: number
    /** hooks */
    items: {
        /** one hook */
        [id: number]: IHook<T>
    }
}

/** all hooks (internal) */
export type THooks<T> = {
    [name in THookName]: THooksList<T>
}
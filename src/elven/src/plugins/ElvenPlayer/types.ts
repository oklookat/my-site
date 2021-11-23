import type DOM from "./logic/dom"

export type TConvertSecondsMode = 'auto' | 'hours' | 'minutes'
export type TSource = string
export type TSources = Array<string>

/** store */
export interface IStore<T> {
    /** get value */
    get(value?: T): void
    /** set value */
    set(value?: T): void
    /** hook when value changed. Returns unsubscribe function */
    onChange(hook: (value?: T) => void): () => void
}

/** plugin */
export interface IElvenPlayer {
    state: IState
    dom: DOM
    play(): void
    pause(): void
    stop(): void
    next(): void
    prev(): void
    addToPlaylist: (source: TSource) => void
    set playlist(playlist: TPlaylist)
    get playlist(): TPlaylist
}

/** audio element events */
export interface IEvents {
    /** when audio playing */
    onPlaying: (e?: Event) => void
    /** when audio paused */
    onPause: (e?: Event) => void
    /** when audio ended */
    onEnded: (e?: Event) => void
    /** when audio time updated */
    onTimeUpdate: (e?: Event) => void
    /** when error */
    onError: (e?: ErrorEvent) => void
}

/** playlist */
export type TPlaylist = {
    /** current index of source */
    position: number
    /** list of playable (mostly) audios */
    sources: TSources
}

/** state of player */
export type TStore = {
    /** is audio playing now */
    playing: IStore<boolean>
    volume: {
        /** in float */
        num: IStore<number>
        /** in percents */
        percents: IStore<number>
    }
    /** current playing */
    current: {
        /** is audio ended */
        ended: IStore<boolean>
        /** buffered */
        buffered: {
            /** in percents */
            percents: IStore<number>
        }
        /** total time */
        duration: {
            /** in seconds */
            num: IStore<number>
            /** in string like '04:20' */
            pretty: IStore<string>
        }
        /** current time */
        time: {
            /** in seconds */
            num: IStore<number>
            /** in percents */
            percents: IStore<number>
            /** in string '01:37' */
            pretty: IStore<string>
        }
    }
}

/** local copy of TStore (almost), 
 * that updating by subscription or by external changes. 
 * Used in component */
export type TComponentStore = {
    playing: boolean
    volume: {
        percents: number
    }
    current: {
        buffered: {
            percents: number
        }
        time: {
            /** is user dragging time slider */
            draggingNow: boolean
            percents: number
            pretty: string
        }
        duration: {
            pretty: string
        }
    }
}

/** audio element state */
export interface IState {
    store: TStore
    set playing(v: boolean)
    set ended(v: boolean)

    set bufferedPercents(v: number)

    set durationNum(v: number)
    set durationPretty(v: string)

    set currentTimeNum(v: number)
    set currentTimePercents(v: number)
    set currentTimePretty(v: string)

    set volumePercents(v: number)
    set volumeNum(v: number)
}
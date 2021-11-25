import type DOM from "./logic/dom"

export type TConvertSecondsMode = 'auto' | 'hours' | 'minutes'
export type TSource = string
export type TSources = Array<string>

export interface IStateUnsubscriber {
    (): void
}

/** state management */
export interface IState<T> {
    /** get value */
    get(): T
    /** set value */
    set(value?: T): void
    /** hook when value changed. Returns unsubscribe function */
    onChange(hook: (value?: T) => void): IStateUnsubscriber
}


/** audio element store with state */
export interface IStore {
    state: {    
        /** is audio playing now */
        playing: IState<boolean>
        volume: {
            /** in float */
            num: IState<number>
            /** in percents */
            percents: IState<number>
        }
        /** current playing */
        current: {
            /** is audio ended */
            ended: IState<boolean>
            /** buffered */
            buffered: {
                /** in percents */
                percents: IState<number>
            }
            /** total time */
            duration: {
                /** in seconds */
                num: IState<number>
                /** in string like '04:20' */
                pretty: IState<string>
            }
            /** current time */
            time: {
                /** in seconds */
                num: IState<number>
                /** in percents */
                percents: IState<number>
                /** in string '01:37' */
                pretty: IState<string>
            }
        }
    }

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

/** local copy of IStore.state (almost), 
 * that updating by subscription or by external changes. 
 * Used in component */
export type TComponentState = {
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

/** plugin */
export interface IElvenPlayer {
    store: IStore
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
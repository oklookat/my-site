import type { Writable } from 'svelte/store'

export type TConvertSecondsMode = 'auto' | 'hours' | 'minutes'
export type TSource = string
export type TSources = Array<string>
export interface IElvenPlayer {
    addToPlaylist: (source: TSource) => void
    setPlaylist: (playlist: TPlaylist) => void
    play: (source: TSource) => void
}
export interface IEvents {
    onPlaying: (e?: Event) => void
    onPause: (e?: Event) => void
    onEnded: (e?: Event) => void
    onTimeUpdate: (e?: Event) => void
    onError: (e?: ErrorEvent) => void
}

export type TPlaylist = {
    position: number
    sources: TSources
}

export type TStore = {
    isPlaying: Writable<boolean>
    /** controls volume */
    volume: {
        /** in float (1.0) */
        num: Writable<number>,
        /** in percents (100) */
        percents: Writable<number>
    }
    current: {
        duration: {
            /** total time of audio in seconds */
            num: Writable<number>
            /** total time of audio like '04:23' */
            pretty: Writable<string>
        }
        position: {
            /** current time of audio in seconds */
            num: Writable<number>
            /** current time of audio like '01:23' */
            pretty: Writable<string>
            /** percents of audio reached */
            percents: Writable<number>
        }
        buffered: {
            /** percents of buffered audio */
            percents: Writable<number>
        }
    }
}

export interface IState {
    store: TStore
    set isPlaying(v: boolean)
    set positionNum(v: number)
    set positionPercents(v: number)
    set positionPretty(v: string)
    set bufferedPercents(v: number)
    set durationNum(v: number)
    set durationPretty(v: string)
    set volumePercents(v: number)
    set volumeNum(v: number)
}
import type { Writable } from 'svelte/store'

export type TConvertSecondsMode = 'auto' | 'hours' | 'minutes'
export type TSource = string
export type TSources = Array<string>

/** plugin */
export interface IElvenPlayer {
    addToPlaylist: (source: TSource) => void
    setPlaylist: (playlist: TPlaylist) => void
    play: (source: TSource) => void
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
    playing: Writable<boolean>
    /** controls volume */
    volume: {
        /** in float (1.0) */
        num: Writable<number>,
        /** in percents (100) */
        percents: Writable<number>
    }
    current: {
        /** is audio ended */
        ended: Writable<boolean>
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

/** audio element state */
export interface IState {
    store: TStore
    set playing(v: boolean)
    set ended(v: boolean)
    set positionNum(v: number)
    set positionPercents(v: number)
    set positionPretty(v: string)
    set bufferedPercents(v: number)
    set durationNum(v: number)
    set durationPretty(v: string)
    set volumePercents(v: number)
    set volumeNum(v: number)
}
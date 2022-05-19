export enum Signal {
    INITIAL = 0,
    PLAY = 1,
    PAUSE = 2
}

export interface ElvenPlayer {
    addToPlaylist: (src: URL) => void
    clearPlaylist: () => void
    play: () => Promise<void>
    pause: () => void
    next: () => Promise<void>
    prev: () => Promise<void>
}

export type Playlist = {
    currentPosition: 0,
    position: 0
    sources: URL[]
}
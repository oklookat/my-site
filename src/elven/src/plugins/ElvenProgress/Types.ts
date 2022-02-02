/** controls */
export interface ElvenProgress {
    get percents(): number
    set percents(percents: number)
    /** start basic loading (progress go to n percents and stop) */
    startBasic: () => void
    /** finish basic loading (progress go from n percents to 100 and resets) */
    finishBasic: () => void
    /** reset */
    reset: () => void
}

/** settings */
export type Settings = {
    /** height of progress line */
    height: string
    basicLoading: {
        /** go to n percents */
        startTo: number
        /** start with speed */
        startSpeed: number
        /** finish with speed */
        finishSpeed: number
    }
}
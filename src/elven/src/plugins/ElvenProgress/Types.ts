export interface IElvenProgress {
    /** start basic loading (progress go to n percents and stop) */
    startBasic: () => void
    /** finish basic loading (progress go from n percents to 100 and resets) */
    finishBasic: () => void
    setPercents: (percents: number) => void
    resetPercents: () => void
}
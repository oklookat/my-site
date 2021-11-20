import type { TConvertSecondsMode } from "../types"

export default class Utils {

    /** get percents of value */
    public static getPercents(value: number, total: number): number {
        let percents = (value / total) * 100
        percents = Math.round(percents)
        percents = percents > 100 ? 100 : percents < 0 ? 0 : percents
        return percents
    }

    /** convert seconds to string like '01:23' */
    public static getPretty(seconds: number, mode: TConvertSecondsMode): string {
        // https://stackoverflow.com/a/1322771/16762009
        const modes = {
            hours: [11, 8],
            minutes: [14, 5]
        }
        let sub = [0, 0]
        switch (mode) {
            case 'auto':
                sub = seconds < 3600 ? modes.minutes : modes.hours
                break
            case 'hours':
                sub = modes.hours
                break
            case 'minutes':
                sub = modes.minutes
                break
        }
        return new Date(seconds * 1000).toISOString().substr(sub[0], sub[1])
    }

    /** get buffered percents */
    public static getBufferedPercents(currentTime: number, duration: number, buffered: TimeRanges): number {
        if (duration > 0) {
            for (let i = 0; i < buffered.length; i++) {
                const len = buffered.length - 1 - i
                if (buffered.start(len) < currentTime) {
                    const perc = this.getPercents(buffered.end(len), duration)
                    return Math.round(perc)
                }
            }
        }
        return 0
    }

    /** get current position like '01:23' */
    public static getPositionPretty(currentTime: number, duration: number): string {
        const mode: TConvertSecondsMode = duration < 3600 ? 'minutes' : 'hours'
        return this.getPretty(currentTime, mode)
    }

    /** convert percents to position in seconds */
    public static percentsToCurrentTime(perc: number, duration: number): number {
        return Math.round((duration / 100) * perc)
    }
}
import type { TConvertSecondsMode } from "../types"

export default class Utils {

    public static computePercents(current: number, total: number): number {
        let percents = (current / total) * 100
        percents = Math.round(percents)
        if (percents >= 100) {
            percents = 100
        } else if (total < 1) {
            percents = 0
        }
        return percents
    }

    // convert seconds to string like '01:23'
    public static convertSeconds(seconds: number, mode: TConvertSecondsMode): string {
        const returnHours = (seconds: number) => {
            return new Date(seconds * 1000).toISOString().substr(11, 8)
        }

        const returnMinutes = (seconds: number) => {
            return new Date(seconds * 1000).toISOString().substr(14, 5)
        }

        // https://stackoverflow.com/a/1322771/16762009
        switch (mode) {
            case 'auto':
                if (seconds < 3600) {
                    // like 00:01
                    return returnMinutes(seconds)
                } else {
                    // like 01:23:12
                    return returnHours(seconds)
                }
            case 'hours':
                // like 01:23:12
                return returnHours(seconds)
            case 'minutes':
                // like 00:01
                return returnMinutes(seconds)
        }
    }

    public static computeBuffered(currentTime: number, duration: number, buffered: TimeRanges): number {
        currentTime = Math.round(currentTime)
        if (duration > 0) {
            for (let i = 0; i < buffered.length; i++) {
                const len = buffered.length - 1 - i
                if (buffered.start(len) < currentTime) {
                    return Math.round(this.computePercents(buffered.end(len), duration))
                }
            }
        }
        return 0
    }
}
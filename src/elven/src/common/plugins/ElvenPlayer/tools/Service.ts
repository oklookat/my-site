export type TConvertSecondsMode = 'auto' | 'hours' | 'minutes'

// only static methods be here
export default class Service {

    // get percents of current param by setting the total param
    public static computePercents(current: number, total: number): number {
        current = Math.round(current)
        let percents = (current / total) * 100
        if (percents >= 100) {
            percents = 100
        } else if (total < 1) {
            percents = 0
        }
        return this.round(percents, 4)
    }

    // round the number to the specific number of decimal places
    public static round(value: number, precision: number): number {
        // https://stackoverflow.com/a/7343013/16762009
        const multiplier = Math.pow(10, precision || 0)
        return Math.round(value * multiplier) / multiplier
    }

    // convert seconds to string like '01:23'
    public static convertSeconds(seconds: number, mode: TConvertSecondsMode): string {
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

        function returnHours(seconds: number) {
            return new Date(seconds * 1000).toISOString().substr(11, 8)
        }

        function returnMinutes(seconds: number) {
            return new Date(seconds * 1000).toISOString().substr(14, 5)
        }
    }

    // get pageX (horizontal mouse position) by touch or mouse event
    public static getPageX(event): number | null {
        let pageX = 0
        const isMovedByTouchscreen = event.type.includes('touch') && event.touches && event.touches.length > 0
        const isMovedByMouse = event.type.includes('mouse')
        if (isMovedByTouchscreen) {
            for (const touch of event.touches) {
                if (touch.pageX) {
                    pageX = touch.pageX
                    break
                }
            }
        } else if (isMovedByMouse) {
            // if moved by mouse
            pageX = event.pageX
        }
        const isMovedByUnknown = !isMovedByTouchscreen && !isMovedByMouse
        const isUnknownValue = isNaN(pageX)
        if (isMovedByUnknown || isUnknownValue) {
            return null
        }
        return pageX
    }


    // get click position by pageX
    public static getClickPosition(clientX: number, element: HTMLElement): number {
        return (clientX - element.offsetLeft) / element.offsetWidth
    }

    // computing how much is buffered
    public static computeBuffered(playerEL: HTMLAudioElement): number {
        const currentTime = Math.round(playerEL.currentTime)
        const duration = playerEL.duration
        if (duration > 0) {
            for (let i = 0; i < playerEL.buffered.length; i++) {
                const len = playerEL.buffered.length - 1 - i
                if (playerEL.buffered.start(len) < currentTime) {
                    return Math.round(this.computePercents(playerEL.buffered.end(len), duration))
                }
            }
        }
        return 0
    }
}
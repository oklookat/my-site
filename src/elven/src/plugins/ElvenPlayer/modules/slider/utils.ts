export default class Utils {

    /** get pageX (horizontal mouse position) by mouse or touch event */
    public static getPageX(e: MouseEvent | TouchEvent): number {
        const movedByTouch = e instanceof TouchEvent && e.touches.length > 0
        if (movedByTouch) {
            // get first touch
            return e.touches[0].pageX
        }
        const movedByMouse = e instanceof MouseEvent
        if (movedByMouse) {
            return e.pageX
        }
        return 0
    }

    /**
     * get click value (width)
     * @param pageX target pageX
     * @param container container where target placed
     * @returns width value. 
     * value < 0 - end of element (left); value > 1 end of element (right). Values between 0 and 1 means you inside target.
     * Multiply by 100 gives you percents.
     */
    public static getClickPercentsWidth(pageX: number, container: HTMLElement): number {
        return (pageX - container.offsetLeft) / container.offsetWidth
    }

    /** get percents of current param by setting the total param */
    public static computePercents(current: number, total: number): number {
        let percents = (current / total) * 100
        percents = Math.round(percents)
        percents = percents > 100 ? 100 : percents < 0 ? 0 : percents
        return percents
    }
    
}
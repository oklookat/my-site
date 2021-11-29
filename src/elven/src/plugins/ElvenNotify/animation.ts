/** contains cool animations / transitions(?) */
export default class Animation {

    /** opacity: 0 to 1. Resolves the promise when animation ended */
    public static fadeIn(div: HTMLDivElement): Promise<void> {
        let opacity = 0
        return new Promise(resolve => {
            const anim = () => {
                opacity = opacity + 0.1
                div.style.opacity = opacity.toString(10)
                if (opacity > 1) {
                    clearInterval(interval)
                    resolve()
                }
            }
            const interval = setInterval(anim, 20)
        })
    }

    /** opacity: 1 to 0. Resolves the promise when animation ended */
    public static fadeOut(div: HTMLDivElement): Promise<void> {
        let opacity = 1
        return new Promise(resolve => {
            const anim = () => {
                opacity = opacity - 0.1
                div.style.opacity = opacity.toString(10)
                if (opacity < 0) {
                    clearInterval(interval)
                    resolve()
                }
            }
            const interval = setInterval(anim, 20)
        })
    }

}
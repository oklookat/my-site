/** prints information */
export default class Logger {

    private static readonly prefix = `[elvenPlayer]`

    public static info(message: string) {
        console.log(`${this.prefix} ${message}`)
    }

    public static warn(message: string) {
        console.warn(`${this.prefix} ${message}`)
    }

    public static error(message: string) {
        console.error(`${this.prefix} ${message}`)
    }

}
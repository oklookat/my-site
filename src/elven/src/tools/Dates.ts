export default class Dates {

    private static readonly hoursTitles = ['hour', 'hours']
    private static readonly minTitles = ['minute', 'minutes']
    private static readonly secTitles = ['second', 'seconds']
    private static readonly months = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec']

    /** convert date to string like: '12 minutes ago' or '12 jan 1970 at 12:22' */
    public static convert(date: string | number | Date): string {
        let d = new Date(date)
        const currentDate = new Date()
        const isCurrentYear = currentDate.getFullYear() === d.getFullYear()
        const isCurrentMonth = currentDate.getMonth() === d.getMonth()
        const today = isCurrentYear && isCurrentMonth && currentDate.getDate() === d.getDate()
        if (today) {
            const secondsAgo = Math.round((currentDate.getTime() - d.getTime()) / 1000)
            // n seconds ago
            if (secondsAgo < 60) {
                const word = this.declensionOfNumbers(secondsAgo, this.secTitles)
                return `${secondsAgo} ${word}`
            }
            const minutesAgo = Math.round(secondsAgo / 60)
            // n minutes ago
            if (minutesAgo < 60) {
                const word = this.declensionOfNumbers(minutesAgo, this.minTitles)
                return `${minutesAgo} ${word}`
            }
            const hoursAgo = Math.round(minutesAgo / 60)
            // hour ago
            if (hoursAgo < 2) {
                return `hour`
            }
            // n hours ago
            if (hoursAgo < 9) {
                const word = this.declensionOfNumbers(hoursAgo, this.hoursTitles)
                return `${hoursAgo} ${word}`
            }
            // today
            const hours = this.numberWithZero(d.getHours())
            const minutes = this.numberWithZero(d.getMinutes())
            return `${hours}:${minutes}`
        }
        const yesterday = isCurrentYear && isCurrentMonth && currentDate.getDate() - d.getDate() === 1
        if (yesterday) {
            const hours = this.numberWithZero(d.getHours())
            const minutes = this.numberWithZero(d.getMinutes())
            return `yesterday at ${hours}:${minutes}`
        }
        // current year?
        const day = this.numberWithZero(d.getDate())
        const month = d.getMonth()
        // in JS Date() January starts with 0, not 1, so we can use it in months array without subtraction
        const monthName = this.months[month]
        const hours = this.numberWithZero(d.getHours())
        const minutes = this.numberWithZero(d.getMinutes())
        if (isCurrentYear) {
            return `${day} ${monthName} at ${hours}:${minutes}`
        } else {
            return `${day} ${monthName} ${d.getFullYear()} at ${hours}:${minutes}`
        }
    }

    /** second(s)? minute(s)? hour(s)? */
    private static declensionOfNumbers(number: number, titles: string[]): string {
        number = Math.abs(number)
        if (number === 1) {
            return titles[0]
        }
        return titles[1]
    }

    /** example: 9 = 09 */
    private static numberWithZero(value: number | string) {
        return ("0" + value).slice(-2)
    }

}
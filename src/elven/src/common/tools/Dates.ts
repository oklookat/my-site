const hoursTitles = ['hour', 'hours']
const minTitles = ['minute', 'minutes']
const secTitles = ['second', 'seconds']
const months = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec']

class Dates {
    public static convert(date) {
        date = new Date(date)
        const currentDate = new Date()
        const isYesterday = currentDate.getFullYear() === date.getFullYear() && currentDate.getMonth() === date.getMonth() && currentDate.getDate() > date.getDate() && currentDate.getDate() - date.getDate() === 1
        const isToday = currentDate.getFullYear() === date.getFullYear() && currentDate.getMonth() === date.getMonth() && currentDate.getDate() === date.getDate()
        if (isToday) {
            const secondsAgo = Math.round((currentDate.getTime() - date.getTime()) / 1000)
            const minutesAgo = Math.round(secondsAgo / 60)
            const hoursAgo = Math.round(minutesAgo / 60)
            if (secondsAgo < 60) {
                // 1min-
                date = `${secondsAgo} ${this.declensionOfNumbers(secondsAgo, secTitles)} ago`
            } else if (minutesAgo < 60) {
                // 1h-
                date = `${minutesAgo} ${this.declensionOfNumbers(minutesAgo, minTitles)} ago`
            } else if (hoursAgo < 2) {
                // 2h-
                date = `hour ago`
            } else if (hoursAgo < 8) {
                // 8h-
                date = `${hoursAgo} ${this.declensionOfNumbers(hoursAgo, hoursTitles)} ago`
            } else if (hoursAgo < 24) {
                // 24h-
                const hours = this.numberWithZero(date.getHours())
                const minutes = this.numberWithZero(date.getMinutes())
                date = `today at ${hours}:${minutes}`
            }
        } else if (isYesterday) {
            date = `yesterday at ${this.numberWithZero(date.getHours())}:${this.numberWithZero(date.getMinutes())}`
        } else {
            const day = this.numberWithZero(date.getDate())
            const month = date.getMonth() // in JS Date() January starts with 0, not 1
            const hours = this.numberWithZero(date.getHours())
            const minutes = this.numberWithZero(date.getMinutes())
            if (date.getFullYear() === currentDate.getFullYear()) {
                date = `${day} ${months[month]} at ${hours}:${minutes}`
            } else {
                date = `${day} ${months[month]} ${date.getFullYear()} at ${hours}:${minutes}`
            }
        }
        return date
    }

    public static declensionOfNumbers(number, titles): string {
        number = Math.abs(number)
        if (number === 1) {
            return titles[0]
        }
        return titles[1]
    }

    public static numberWithZero(date) {
        return ("0" + date).slice(-2)
    }
}

export default Dates
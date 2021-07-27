const hoursTitles = ['час', 'часа', 'часов']
const minTitles = ['минуту', 'минуты', 'минут']
const secTitles = ['секунду', 'секунды', 'секунд']
const months = ['янв', 'фев', 'мар', 'апр', 'мая', 'июн', 'июл', 'авг', 'сен', 'ноя', 'дек']

class ElvenDates {
    public static convert(date) {
        date = new Date(date)
        const currentDate = new Date()
        const isYesterday = currentDate.getFullYear() === date.getFullYear() && currentDate.getMonth() === date.getMonth() && currentDate.getDay() > date.getDay() && currentDate.getDay() - date.getDay() === 1
        const isToday = currentDate.getFullYear() === date.getFullYear() && currentDate.getMonth() === date.getMonth() && currentDate.getDay() === date.getDay()
        if (isToday) {
            const secondsAgo = Math.round((currentDate.getTime() - date.getTime()) / 1000)
            const minutesAgo = Math.round(secondsAgo / 60)
            const hoursAgo = Math.round(minutesAgo / 60)
            if (secondsAgo < 60) {
                // 1min-
                date = `${secondsAgo} ${this.declensionOfNumbers(secondsAgo, secTitles)} назад`
            } else if (minutesAgo < 60) {
                // 1h-
                date = `${minutesAgo} ${this.declensionOfNumbers(minutesAgo, minTitles)} назад`
            } else if (hoursAgo < 2) {
                // 2h-
                date = `час назад`
            } else if (hoursAgo < 8) {
                // 8h-
                date = `${hoursAgo} ${this.declensionOfNumbers(hoursAgo, hoursTitles)} назад`
            } else if (hoursAgo < 24) {
                // 24h-
                const hours = this.numberWithZero(date.getHours())
                const minutes = this.numberWithZero(date.getMinutes())
                date = `сегодня в ${hours}:${minutes}`
            }
        } else if (isYesterday) {
            date = `вчера в ${this.numberWithZero(date.getHours())}:${this.numberWithZero(date.getMinutes())}`
        } else {
            const day = this.numberWithZero(date.getDate())
            const month = date.getMonth() // in JS Date() January starts with 0, not 1
            const hours = this.numberWithZero(date.getHours())
            const minutes = this.numberWithZero(date.getMinutes())
            if (date.getFullYear() === currentDate.getFullYear()) {
                date = `${day} ${months[month]} в ${hours}:${minutes}`
            } else {
                date = `${day} ${months[month]} ${date.getFullYear()} в ${hours}:${minutes}`
            }
        }
        return date
    }

    public static declensionOfNumbers(number, titles) {
        // https://gist.github.com/realmyst/1262561#gistcomment-3443551
        number = Math.abs(number)
        if (Number.isInteger(number)) {
            const cases = [2, 0, 1, 1, 1, 2]
            return titles[(number % 100 > 4 && number % 100 < 20) ? 2 : cases[(number % 10 < 5) ? number % 10 : 5]]
        }
        return titles[1]
    }

    public static numberWithZero(date) {
        return ("0" + date).slice(-2)
    }
}

export default ElvenDates
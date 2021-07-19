const hoursTitles = ['час', 'часа', 'часов']
const minTitles = ['минуту', 'минуты', 'минут']
const secTitles = ['секунду', 'секунды', 'секунд']

class ElvenTools{
    public static convertDate(date){
        date = new Date(date)
        const currentDate = new Date()
        const isYesterday = currentDate.getFullYear() === date.getFullYear() && currentDate.getMonth() === date.getMonth() && currentDate.getDay() > date.getDay() && currentDate.getDay() - date.getDay() === 1
        const isToday = currentDate.getFullYear() === date.getFullYear() && currentDate.getMonth() === date.getMonth() && currentDate.getDay() === date.getDay()
        if(isYesterday){
            date = `вчера в ${date.getHours()}:${date.getMinutes()}`
        } else if(isToday){
            const secondsAgo = Math.round((currentDate.getTime() - date.getTime()) / 1000)
            const minutesAgo = Math.round(secondsAgo / 60)
            const hoursAgo = Math.round(minutesAgo / 60)
            if(secondsAgo < 60){
                // 1min-
                date = `${secondsAgo} ${this.declensionOfNumbers(secondsAgo, secTitles)} назад`
            } else if(minutesAgo < 60){
                // 1h-
                date = `${minutesAgo} ${this.declensionOfNumbers(minutesAgo, minTitles)} назад`
            } else if(hoursAgo === 1){
                // 2h-
                date = `час назад`
            } else if(hoursAgo < 8){
                // 8h-
                date = `${hoursAgo} ${this.declensionOfNumbers(hoursAgo, hoursTitles)} назад`
            } else if(hoursAgo < 24){
                // 24h-
                date = `в ${date.getHours()}:${("0" + date.getMinutes()).slice(-2)}`
            }
        } else{
            console.log('хухухух')
            date = `${("0" + date.getDay()).slice(-2)}.${("0" + date.getMonth()).slice(-2)}.${date.getFullYear()} в ${("0" + date.getHours()).slice(-2)}:${("0" + date.getMinutes()).slice(-2)}`
        }
        return date
    }

    public static declensionOfNumbers(number, titles){
        // https://gist.github.com/realmyst/1262561#gistcomment-3443551
        number = Math.abs(number);
        if (Number.isInteger(number)) {
            const cases = [2, 0, 1, 1, 1, 2];
            return titles[ (number%100>4 && number%100<20)? 2 : cases[(number%10<5)?number%10:5] ];
        }
        return titles[1]
    }

    private static numberWithZero(number){
        return ("0" + number).slice(-2)
    }
}

export default ElvenTools
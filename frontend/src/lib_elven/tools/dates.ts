const hoursTitles = ['hour', 'hours']
const minTitles = ['minute', 'minutes']
const secTitles = ['second', 'seconds']
const months = ['jan', 'feb', 'mar', 'apr', 'may', 'jun', 'jul', 'aug', 'sep', 'oct', 'nov', 'dec']

/** convert date to string like: '12 minutes ago' or '12 jan 1970 at 12:22' */
export function dateToReadable(date: string | number | Date): string {
	const argDate = new Date(date)
	const currentDate = new Date()

	const isCurrentYear = currentDate.getFullYear() === argDate.getFullYear()
	const isCurrentMonth = currentDate.getMonth() === argDate.getMonth()

	const isToday = isCurrentYear && isCurrentMonth && currentDate.getDate() === argDate.getDate()
	if (isToday) {
		// n seconds ago
		const secondsAgo = Math.round((currentDate.getTime() - argDate.getTime()) / 1000)
		if (secondsAgo < 60) {
			const word = declensionOfNumbers(secondsAgo, secTitles)
			return `${secondsAgo} ${word}`
		}

		// n minutes ago
		const minutesAgo = Math.round(secondsAgo / 60)
		if (minutesAgo < 60) {
			const word = declensionOfNumbers(minutesAgo, minTitles)
			return `${minutesAgo} ${word}`
		}

		// hour ago
		const hoursAgo = Math.round(minutesAgo / 60)
		if (hoursAgo < 2) {
			return `hour`
		}

		// n hours ago
		if (hoursAgo < 9) {
			const word = declensionOfNumbers(hoursAgo, hoursTitles)
			return `${hoursAgo} ${word}`
		}

		// today
		const hours = numberWithZero(argDate.getHours())
		const minutes = numberWithZero(argDate.getMinutes())
		return `${hours}:${minutes}`
	}

	const isYesterday = isCurrentYear && isCurrentMonth && currentDate.getDate() - argDate.getDate() === 1
	if (isYesterday) {
		const hours = numberWithZero(argDate.getHours())
		const minutes = numberWithZero(argDate.getMinutes())
		return `yesterday at ${hours}:${minutes}`
	}

	// current year?
	const day = numberWithZero(argDate.getDate())
	const month = argDate.getMonth()
	const monthName = months[month]
	const hours = numberWithZero(argDate.getHours())
	const minutes = numberWithZero(argDate.getMinutes())
	if (isCurrentYear) {
		return `${day} ${monthName} at ${hours}:${minutes}`
	}
	return `${day} ${monthName} ${argDate.getFullYear()} at ${hours}:${minutes}`
}

/** second(s)? minute(s)? hour(s)? */
function declensionOfNumbers(number: number, titles: string[]): string {
	number = Math.abs(number)
	if (number === 1) {
		return titles[0]
	}
	return titles[1]
}

/** example: 9 = 09 */
function numberWithZero(value: number | string) {
	return ('0' + value).slice(-2);
}

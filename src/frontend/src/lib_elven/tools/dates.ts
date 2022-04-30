import { stringToNumber } from "$lib_elven/tools";
import { getDateFormatter, getMessageFormatter, getTimeFormatter, _, json } from "svelte-i18n";


let i18n_hourAgo = 'hour ago'
let i18n_hoursAgo = 'hours ago'
let i18n_minutesAgo = 'minutes ago'
let i18n_secondsAgo = 'seconds ago'
let i18n_yesterday = 'yesterday'
// let i18n_months = [
// 	'jan',
// 	'feb',
// 	'mar',
// 	'apr',
// 	'may',
// 	'jun',
// 	'jul',
// 	'aug',
// 	'sep',
// 	'oct',
// 	'nov',
// 	'dec'
// ];

// i18n
json.subscribe(v => {
	i18n_hourAgo = v('general.hourAgo') as any
	i18n_hoursAgo = v('general.hoursAgo') as any
	i18n_minutesAgo = v('general.minutesAgo') as any
	i18n_secondsAgo = v('general.secondsAgo') as any
	i18n_yesterday = v('general.yesterday') as any
	//i18n_months = v('general.months') as any
})()

/** convert date to string like: '12 minutes ago' or '12 jan 1970 at 12:22' */
export function dateToReadable(date: string | number | Date): string {
	let d = new Date(date);
	const currentDate = new Date();

	const dateFormatter = getDateFormatter({
		'month': 'short',
		'day': '2-digit',
		'year': 'numeric',
		'hour': '2-digit',
	})

	let monthName = ''
	let year = ''
	let dayPeriod = ''
	for(const form of dateFormatter.formatToParts(d)) {
		switch(form.type) {
			case 'month': 
				monthName = form.value
				continue
			case 'year':
				year = form.value
				continue
			case 'dayPeriod':
				dayPeriod = ' '+ form.value
				continue
		}
	}

	const isCurrentYear = currentDate.getFullYear() === d.getFullYear();
	const isCurrentMonth = currentDate.getMonth() === d.getMonth();

	const isToday = isCurrentYear && isCurrentMonth && currentDate.getDate() === d.getDate();

	const hours = numberWithZero(d.getHours());
	const minutes = numberWithZero(d.getMinutes());
	const hoursAndMinutes = `${hours}:${minutes}${dayPeriod}`

	if (isToday) {
		const secondsAgo = Math.round((currentDate.getTime() - d.getTime()) / 1000);

		// n seconds ago
		if (secondsAgo < 60) {
			return `${secondsAgo} ${i18n_secondsAgo}`;
		}
		const minutesAgo = Math.round(secondsAgo / 60);

		// n minutes ago
		if (minutesAgo < 60) {
			return `${minutesAgo} ${i18n_minutesAgo}`;
		}
		const hoursAgo = Math.round(minutesAgo / 60);

		// hour ago
		if (hoursAgo < 2) {
			return `${i18n_hourAgo}`;
		}

		// n hours ago
		if (hoursAgo < 9) {
			return `${hoursAgo} ${i18n_hoursAgo}`;
		}

		// today
		return `${hoursAndMinutes}`;
	}

	const isYesterday = isCurrentYear && isCurrentMonth && currentDate.getDate() - d.getDate() === 1;
	if (isYesterday) {
		return `${i18n_yesterday}, ${hoursAndMinutes}`;
	}

	// current year?
	const day = numberWithZero(d.getDate());

	if (isCurrentYear) {
		return `${day} ${monthName} ${hoursAndMinutes}`;
	}

	return `${day} ${monthName} ${year}, ${hoursAndMinutes}`;
}


/** example: 9 = 09 */
function numberWithZero(value: number | string) {
	return ('0' + value).slice(-2);
}
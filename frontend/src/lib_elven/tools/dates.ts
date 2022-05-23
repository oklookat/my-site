import { browser } from '$app/env';

let formatter: Intl.DateTimeFormat;
let locale = 'en-US';
if (browser) {
	locale = navigator.language;
	formatter = new Intl.DateTimeFormat(locale, {
		month: 'short',
		day: '2-digit',
		year: 'numeric',
		hour: '2-digit',
		second: '2-digit'
	});
}

/** convert date to string like: '12 minutes ago' or '12 jan 1970 at 12:22' */
export function dateToReadable(date: string | number | Date): string {
	if (!browser) {
		return '';
	}
	let d = new Date(date);
	const currentDate = new Date();

	let monthName = '';
	let year = '';
	let dayPeriod = '';
	let second = '';
	for (const form of formatter.formatToParts(d)) {
		switch (form.type) {
			case 'month':
				monthName = form.value;
				continue;
			case 'year':
				year = form.value;
				continue;
			case 'dayPeriod':
				dayPeriod = ' ' + form.value;
				continue;
			case 'second':
				second = form.value;
				continue;
		}
	}
	console.log(second);

	const isCurrentYear = currentDate.getFullYear() === d.getFullYear();
	const isCurrentMonth = currentDate.getMonth() === d.getMonth();

	const isToday = isCurrentYear && isCurrentMonth && currentDate.getDate() === d.getDate();

	const hours = numberWithZero(d.getHours());
	const minutes = numberWithZero(d.getMinutes());
	const hoursAndMinutes = `${hours}:${minutes}${dayPeriod}`;

	if (isToday) {
		const secondsAgo = Math.round((currentDate.getTime() - d.getTime()) / 1000);

		// n seconds ago
		if (secondsAgo < 60) {
			return `${secondsAgo} ${'seconds ago'}`;
		}
		const minutesAgo = Math.round(secondsAgo / 60);

		// n minutes ago
		if (minutesAgo < 60) {
			return `${minutesAgo} ${'minutes ago'}`;
		}
		const hoursAgo = Math.round(minutesAgo / 60);

		// hour ago
		if (hoursAgo < 2) {
			return `${'hour ago'}`;
		}

		// n hours ago
		if (hoursAgo < 9) {
			return `${hoursAgo} ${'hours ago'}`;
		}

		// today
		return `${hoursAndMinutes}`;
	}

	const isYesterday = isCurrentYear && isCurrentMonth && currentDate.getDate() - d.getDate() === 1;
	if (isYesterday) {
		return `${'yesterday'}, ${hoursAndMinutes}`;
	}

	// current year?
	const day = numberWithZero(d.getDate());

	if (isCurrentYear) {
		return `${day} ${monthName}, ${hoursAndMinutes}`;
	}

	return `${day} ${monthName} ${year}, ${hoursAndMinutes}`;
}

/** example: 9 = 09 */
function numberWithZero(value: number | string) {
	return ('0' + value).slice(-2);
}

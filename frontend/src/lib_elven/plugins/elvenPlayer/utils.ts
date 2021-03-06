type ConvertSecondsMode = 'auto' | 'hours' | 'minutes';

/** get percents of value */
export function getPercents(value: number, total: number): number {
	let percents = (value / total) * 100;
	percents = Math.round(percents);
	if (percents > 100) {
		percents = 100;
	} else if (percents < 0) {
		percents = 0;
	}
	return percents;
}

/** convert seconds to string like '01:23' */
export function getPretty(seconds: number, mode: ConvertSecondsMode = 'auto'): string {
	// https://stackoverflow.com/a/1322771/16762009
	seconds = Math.round(seconds);
	const modes = {
		hours: [11, 19],
		minutes: [14, 19]
	};
	let sub = [0, 0];
	switch (mode) {
		case 'auto':
			sub = seconds < 3600 ? modes.minutes : modes.hours;
			break;
		case 'hours':
			sub = modes.hours;
			break;
		case 'minutes':
			sub = modes.minutes;
			break;
	}
	let pretty = '00:00';
	try {
		pretty = new Date(seconds * 1000).toISOString().substring(sub[0], sub[1]);
	} catch (err) {
		// Logger.warn('getPretty(): ' + err)
	}
	return pretty;
}

/** get buffered percents */
export function getBufferedPercents(
	currentTime: number,
	duration: number,
	buffered: TimeRanges
): number {
	if (duration < 0) {
		return 0;
	}
	for (let i = 0; i < buffered.length; i++) {
		const len = buffered.length - 1 - i;
		if (buffered.start(len) > currentTime) {
			continue;
		}
		const perc = getPercents(buffered.end(len), duration);
		return Math.round(perc);
	}
	return 0;
}

/** get current time like '01:23' or '01:23:23' by seconds and duration seconds */
export function convertCurrentTimePretty(currentTime: number, duration: number): string {
	const mode: ConvertSecondsMode = duration < 3600 ? 'minutes' : 'hours';
	return getPretty(currentTime, mode);
}

/** convert current time percents to seconds position by duration */
export function convertPercentsToCurrentTime(percents: number, duration: number): number {
	return (duration / 100) * percents;
}

export function log(message: string) {
	console.log(`[elvenPlayer] ${message}`)
}

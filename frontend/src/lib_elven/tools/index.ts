import type { LoadEvent } from '@sveltejs/kit';
import * as cookie from 'cookie';

// @ts-ignore
const uploadsURL = import.meta.env.VITE_UPLOADS_URL as string;
// @ts-ignore
const apiURL = import.meta.env.VITE_API_URL as string;

/** convert bytes count to string like '10.4 MB' */
export function bytesToReadable(bytes: number): string {
	let converted = 0;
	try {
		const conv = stringToNormal(bytes);
		if (typeof conv === 'number') {
			converted = conv;
		} else {
			return `unknown size`;
		}
	} catch (err) {
		return `unknown size`;
	}

	const exponent = 3;
	const system = 1000;

	const kb = saveTwoAfterDot(converted / Math.pow(10, exponent));
	if (kb < system) {
		return `${kb} KB`;
	}

	const mb = saveTwoAfterDot(converted / Math.pow(10, exponent * 2));
	if (mb < system) {
		return `${mb} MB`;
	}

	const gb = saveTwoAfterDot(converted / Math.pow(10, exponent * 3));
	if (gb < system) {
		return `${gb} GB`;
	}

	const tb = saveTwoAfterDot(converted / Math.pow(10, exponent * 4));
	if (tb < system) {
		return `${tb} TB`;
	}

	const pb = saveTwoAfterDot(converted / Math.pow(10, exponent * 5));
	return `${pb} PB`;
}

/** if string > 31 length = cut it and add '...' */
export function cutString(str: string, maxLength = 31): string {
	if (str.length < maxLength) {
		return str;
	}
	return str.substring(0, maxLength) + '...';
}

/** get Record<> keys length */
export function getRecordLength<T>(record: Record<any, T>): number {
	const isObject = typeof record === 'object' && record !== null;
	if (!isObject) {
		throw Error('not a object');
	}
	return Object.keys(record).length;
}

/** execute function only after ms. All other calls before will be ignored */
export function debounce(f: Function, ms: number) {
	let isCooldown = false;
	return function () {
		if (isCooldown) {
			return;
		}

		// @ts-ignore
		f.apply(this, arguments);

		isCooldown = true;

		setTimeout(() => {
			isCooldown = false;
		}, ms);
	};
}

/** get token from request headers (cookie) */
export function getTokenFromRequestHeaders(headers: Headers): string | null {
	const isValid = !!headers && headers instanceof Headers && headers.has('cookie');
	if (!isValid) {
		return null;
	}
	const cookiesStr = headers.get('cookie');
	if (!cookiesStr) {
		return null;
	}

	let token = '';
	try {
		const parsed = cookie.parse(cookiesStr);
		if (!parsed || !parsed.token) {
			return null;
		}
		token = parsed.token;
	} catch (err) {
		return null;
	}
	return token;
}

export function addTokenToHeaders(token: string, headers: Headers) {
	if (!token || !(headers instanceof Headers)) {
		return;
	}
	headers.append('Authorization', `Elven ${token}`);
}

export function getTokenFromSession(e: LoadEvent): string {
	return e.session.user.token || '';
}

/** returns normalized XY coords depend on click position and element
 *
 * for what? For example: move element on click position
 */
export function correctElementOverflow(el: HTMLElement, evt: MouseEvent): { x: number; y: number } {
	let x = evt.clientX;
	let y = evt.clientY;
	const moveOffset = 10;

	// left-right (X)
	const popupWidth = el.offsetWidth;
	const overflowDifferenceX = x + popupWidth - document.body.clientWidth;
	if (overflowDifferenceX > 0) {
		x = x - overflowDifferenceX - moveOffset;
	}

	// top-bottom (Y)
	const popupHeight = el.offsetHeight;
	const overflowDifferenceY = y + popupHeight - document.body.clientHeight;
	if (overflowDifferenceY > 0) {
		y = y - overflowDifferenceY - moveOffset;
	}

	return { x, y };
}

/** convert string to value depend on type */
export function stringToNormal(value: any): boolean | number | string {
	try {
		const converted = stringToUndefinedOrNull(value);
		return converted;
	} catch (err) { }

	// try to bool
	try {
		const converted = stringToBool(value);
		return converted;
	} catch (err) { }

	// try to number
	try {
		const converted = stringToNumber(value);
		return converted;
	} catch (err) { }

	return `${value}`;
}

/** convert string to boolean. Throws error if convert failed */
export function stringToBool(value: any): boolean {
	if (typeof value === 'boolean') {
		return value;
	}

	if (typeof value !== 'string') {
		throw Error('value is not a string');
	}

	const valued = value.toUpperCase();

	const trueVals = ['T', 'TRUE', 'YES', 'Y', 'ON'];
	if (trueVals.includes(valued)) {
		return true;
	}

	const falseVals = ['F', 'FALSE', 'N', 'NO', 'OFF'];
	if (falseVals.includes(valued)) {
		return false;
	}

	throw Error('value is not convertible');
}

/** convert string to number. Throws error if convert failed */
export function stringToNumber(value: any): number {
	if (typeof value === 'number') {
		return value;
	}

	if (typeof value !== 'string') {
		throw Error('value is not a string');
	}

	const converted = Number(value);
	if (isNaN(converted) || value === '') {
		throw Error('value is not convertible');
	}

	return converted;
}

export function stringToUndefinedOrNull(value: any): string {
	if (typeof value === 'undefined' || value === null) {
		return '';
	}

	if (typeof value !== 'string') {
		throw Error('value is not a string');
	}

	const valued = value.toUpperCase();

	if (valued === 'UNDEFINED' || valued === 'NULL') {
		return '';
	}

	throw Error('value is not convertible');
}

/** get uploads URL (for uploading files etc) */
export function getUploadsURL(): URL {
	return new URL(uploadsURL);
}

/**  get API URL */
export function getApiURL(): URL {
	return new URL(apiURL);
}

/** api uploads + path. Like: 'https://uploads.example.com/yourpath' */
export function getUploadsWith(path: string): URL {
	let final = path;
	const uploads = getUploadsURL().toString();
	if (final.endsWith('/')) {
		final = `${uploads}${final}`;
	} else {
		final = `${uploads}/${final}`;
	}
	return new URL(final);
}

/** get path like: '/elven/yourpath' */
export function getPathWithElvenPrefix(path: string): string {
	let pathd = path;
	if (!pathd.startsWith('/')) {
		pathd = '/' + pathd;
	}
	return `/elven${pathd}`;
}

/** is device with touchscreen & its default input? */
export function isTouchDevice(): boolean {
	return matchMedia('(hover: none)').matches;
}

/** is admin panel page? */
export function isAdminPanelPage(url: URL): boolean {
	if (!(url instanceof URL)) {
		return false;
	}
	return url.pathname.startsWith('/elven');
}

/** is admin panel login page? */
export function isAdminPanelLoginPage(url: URL): boolean {
	if (!(url instanceof URL)) {
		return false;
	}
	return url.pathname.startsWith('/elven/login');
}

/** is admin panel logout page? */
export function isAdminPanelLogoutPage(url: URL): boolean {
	if (!(url instanceof URL)) {
		return false;
	}
	return url.pathname.startsWith('/elven/logout');
}

/** set searchparam to URL */
export function setSearchParam(params: URLSearchParams, name: string, value: any) {
	if (!(params instanceof URLSearchParams)) {
		throw Error('not a searchparams');
	}

	const isString = typeof value === 'string';

	let finalValue = value;

	const isInvalid =
		finalValue === undefined || finalValue === null || (isString && finalValue.length < 1);
	if (isInvalid) {
		params.delete(name);
		return;
	}

	if (!isString) {
		if (typeof finalValue === 'object' && 'toString' in finalValue) {
			finalValue = finalValue.toString();
		} else {
			finalValue = `${finalValue}`;
		}
	}

	params.set(name, value);
}

/** convert URLSearchParams to object */
export function searchParamsToObject(params: URLSearchParams): Object | undefined {
	if (!(params instanceof URLSearchParams)) {
		return;
	}

	const result: Record<any, any> = {};

	params.forEach((value, key) => {
		// convert if needed.
		try {
			const valNormal = stringToNormal(value);
			result[key] = valNormal;
			return;
		} catch (err) {
			return;
		}
	});

	return result;
}

/** add params to URLSearchParams by object */
export function searchParamsByObject(data: Record<string | number, any>): URLSearchParams {
	const params = new URLSearchParams();

	for (const key in data) {
		const value = data[key];
		if (value === undefined || value === null) {
			continue;
		}
		params.append(key, value);
	}

	return params;
}

/** first call = store default body 'no-scroll'
 *
 * second call = set no-scroll if not exists
 * 
 * third call = remove no-scroll if body not have it on first call
 * 
 * @returns body scroll toggle function
 */
export function createBodyScrollToggler(): () => void {
	let active = false
	return () => {
		if (typeof document === 'undefined') {
			console.warn("[createBodyScrollToggler] toggle on client-side, not server")
			return
		}

		const classNames = document.body.className.split(" ")

		if (!active) {
			classNames.push("no-scroll")
			document.body.className = classNames.join(" ").trim()
			active = true
			return
		}

		for (let i = 0; i < classNames.length; i++) {
			const className = classNames[i]
			if (!className) {
				delete classNames[i]
				continue
			}
			if (className === "no-scroll") {
				delete classNames[i]
				break
			}
		}
		document.body.className = classNames.join(" ").trim()
		active = false
	};
}

/**
 * example:
 *
 * val: 528.3333333333334
 *
 * returns: 528.33
 */
export function saveTwoAfterDot(val: number) {
	// val = 528.3333333333334
	// [528, 3333333333334]
	const dotsSplit = val.toString(10).split('.');
	if (dotsSplit.length > 1) {
		// 3333333333334
		let second = dotsSplit[1];
		if (second.length > 2) {
			second = second.slice(0, 2);
		}
		dotsSplit[1] = second;
	}
	const resultString = dotsSplit.join('.');
	return Number(resultString);
}

/**
 * generate random background and contrast text
 * 
 * @see https://stackoverflow.com/questions/11867545/change-text-color-based-on-brightness-of-the-covered-background-area
 * @returns color for text and background (already in rgb format)
 */
export function randomColors(): { text: string; background: string } {
	const rgb = [255, 0, 0];

	// Randomly update colours
	rgb[0] = Math.round(Math.random() * 255);
	rgb[1] = Math.round(Math.random() * 255);
	rgb[2] = Math.round(Math.random() * 255);

	// http://www.w3.org/TR/AERT#color-contrast
	const brightness = Math.round((rgb[0] * 299 + rgb[1] * 587 + rgb[2] * 114) / 1000);
	const textColour = brightness > 125 ? 'black' : 'white';
	const backgroundColour = 'rgb(' + rgb[0] + ',' + rgb[1] + ',' + rgb[2] + ')';
	return {
		text: textColour,
		background: backgroundColour
	};
}

/** opacity: 0 to 1. Resolves promise when animation ended
 *
 * optional: set opacity on your div to 0 before call this animation
 */
 export function fadeIn(div: HTMLDivElement, speed: number = 20): Promise<void> {
	let opacity = 0;
	const setOpacity = (val: number) => {
		div.style.opacity = val.toString(10);
	};
	return new Promise((resolve) => {
		const anim = () => {
			opacity = opacity + 0.1;
			setOpacity(opacity);
			if (opacity < 1) {
				return;
			}
			setOpacity(1);
			clearInterval(interval);
			resolve();
		};
		const interval = setInterval(anim, speed);
	});
}

/** opacity: 1 to 0. Resolves promise when animation ended */
export function fadeOut(div: HTMLDivElement, speed: number = 20): Promise<void> {
	let opacity = 1;
	const setOpacity = (val: number) => {
		div.style.opacity = val.toString(10);
	};
	return new Promise((resolve) => {
		const anim = () => {
			opacity = opacity - 0.1;
			div.style.opacity = opacity.toString(10);
			if (opacity > 0) {
				return;
			}
			setOpacity(0);
			clearInterval(interval);
			resolve();
		};
		const interval = setInterval(anim, speed);
	});
}
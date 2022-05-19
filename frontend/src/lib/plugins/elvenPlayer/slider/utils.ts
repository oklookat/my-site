/** get pageX (horizontal mouse position) by mouse or touch event */
export function getPageX(e: MouseEvent | TouchEvent): number {
	const movedByTouch = window.TouchEvent && e instanceof TouchEvent && e.touches.length > 0;
	if (movedByTouch) {
		// get first touch
		return e.touches[0].pageX;
	}
	const movedByMouse = e instanceof MouseEvent;
	if (movedByMouse) {
		return e.pageX;
	}
	return 0;
}

/**
 * get click value (width)
 * @param pageX target pageX
 * @param container container where target placed
 * @returns width value.
 * value < 0 - end of element (left); value > 1 end of element (right). Values between 0 and 1 means you inside target.
 * Multiply by 100 gives you percents.
 */
export function getClickPercentsWidth(pageX: number, container: HTMLElement): number {
	let pos = (pageX - container.offsetLeft) / container.offsetWidth;
	if (pos < 0) {
		pos = 0;
	} else if (pos > 1) {
		pos = 1;
	}
	pos = pos * 100;
	return pos;
}

/** get percents of current param by setting the total param */
export function computePercents(current: number, total: number): number {
	let percents = (current / total) * 100;
	if (percents > 100) {
		percents = 100;
	} else if (percents < 0) {
		percents = 0;
	}
	return percents;
}

/** contains cool animations / transitions(?) */
export default class Animation {
	/** opacity: 0 to 1. Resolves promise when animation ended
	 *
	 * optional: set opacity on your div to 0 before call this animation
	 */
	public static fadeIn(div: HTMLDivElement, speed: number = 20): Promise<void> {
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
	public static fadeOut(div: HTMLDivElement, speed: number = 20): Promise<void> {
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
}

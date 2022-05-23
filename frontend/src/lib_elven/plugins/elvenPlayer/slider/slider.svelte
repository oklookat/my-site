<script lang="ts">
	import { computePercents, getClickPercentsWidth, getPageX } from './utils';
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';

	let container: HTMLDivElement;

	const dispatch = createEventDispatcher<{
		/** on slider percents changed */
		slide: number;

		mouseDown: void;
		mouseUp: void;
	}>();

	/** set percents */
	export let percents: number = 100;

	/** is mouse down on slider? */
	let isMouseDown = false;
	$: watchMouseDown(isMouseDown);
	function watchMouseDown(val: boolean) {
		if (val) {
			dispatch('mouseDown');
			return;
		}
		dispatch('mouseUp');
	}

	let finalPercents = percents;

	/** watch when percents prop changed */
	$: watchPercents(percents);
	function watchPercents(perc: number) {
		if (isMouseDown) {
			return;
		}
		if (!perc) {
			perc = 0;
		} else if (perc > 100) {
			perc = 100;
		}
		finalPercents = perc;
	}

	onMount(() => {
		if (!browser) {
			return;
		}
		container.onpointerdown = beginSliding;
		container.onpointerup = stopSliding;
		container.onpointercancel = stopSliding;
		container.oncontextmenu = (e) => {
			e.preventDefault();
		};
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		container.onpointerdown = null;
		container.onpointerup = null;
		container.onpointermove = null;
		container.onpointercancel = null;
		isMouseDown = false;
	});

	// click
	function beginSliding(e: PointerEvent) {
		// disallow dragging slider with any mouse buttons except LMB
		if (e.pointerType === 'mouse' && e.button !== 0) {
			return;
		}
		isMouseDown = true;
		container.onpointermove = slide;
		container.setPointerCapture(e.pointerId);

		// start slide because user already clicked
		slide(e);
	}

	// unclick
	function stopSliding(e: PointerEvent) {
		container.onpointermove = null;
		container.releasePointerCapture(e.pointerId);
		isMouseDown = false;
	}

	// sliding now
	function slide(e: PointerEvent) {
		const containerWidth = container.clientWidth;
		const rect = container.getBoundingClientRect();
		const position = e.clientX - rect.left;
		let perc = computePercents(position, containerWidth);
		if (perc > 100) {
			perc = 100;
		} else if (perc < 0) {
			perc = 0;
		}
		finalPercents = perc;
		dispatch('slide', finalPercents);
	}
</script>

<div class="slider" bind:this={container}>
	<div class="itself" style="width: {finalPercents}%" />
	<div class="line" style="left: calc({finalPercents}% - 1%)" />
</div>

<style lang="scss">
	.slider {
		// important for PointerEvent
		touch-action: pan-y;
		user-select: none;

		cursor: pointer;
		box-sizing: border-box;
		position: relative;

		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;

		.itself,
		.line {
			background-color: #918ce6;
		}
		.itself {
			border-top-left-radius: 2px;
			border-bottom-left-radius: 2px;
			top: auto;
			bottom: 0;
			width: 0;
			height: 100%;
		}
		.line {
			position: absolute;
			width: 6px;
			height: 130%;
		}
	}
</style>

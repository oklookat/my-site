<script lang="ts">
	import { correctElementOverflow } from '$lib/tools';
	import { onMount } from 'svelte';

	export let mouseEvent: MouseEvent;

	export let onDisabled: () => void;

	let popupEL: HTMLDivElement;
	let initialClick = true;

	onMount(() => {
		enable(mouseEvent);
	});

	function enable(evt: MouseEvent) {
		if (!evt) {
			return;
		}
		const { x, y } = correctElementOverflow(popupEL, evt);
		// set styles
		popupEL.style.left = `${x}px`;
		popupEL.style.top = `${y}px`;
	}

	function watchClick(e: MouseEvent) {
		if (initialClick) {
			initialClick = false;
			return;
		}
		if (!onDisabled) {
			return;
		}
		e.preventDefault();
		const target = e.target;
		if (target instanceof HTMLElement) {
			const isChild = isChildOfPopup(target);
			if (isChild) {
				return;
			}
		}
		onDisabled();
	}

	function isChildOfPopup(el: HTMLElement): boolean {
		if (el === popupEL) {
			return true;
		}
		while (el.parentNode && el.parentNode.nodeName.toLowerCase() !== 'body') {
			el = el.parentNode as any;
			if (el === popupEL) {
				return true;
			}
		}
		return false;
	}
</script>

<svelte:window on:click={watchClick} />
<div class="popup with-border" bind:this={popupEL}>
	<slot />
</div>

<style lang="scss">
	.popup {
		font-size: 0.9rem;
		z-index: 9999;
		display: block;
		position: absolute;
		background-color: var(--color-level-1);
		border-radius: var(--border-radius);
		min-width: 104px;
		height: max-content;
	}
</style>

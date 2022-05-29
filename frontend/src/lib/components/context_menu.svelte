<script lang="ts">
	import { correctElementOverflow } from '$elven/tools';
	import { onMount } from 'svelte';
	import Portal from 'svelte-portal/src/Portal.svelte';

	export let mouseEvent: MouseEvent;

	export let onDisabled: () => void;

	let contextEL: HTMLDivElement;
	let initialClick = true;

	onMount(() => {
		if (!mouseEvent) {
			return;
		}
		const { x, y } = correctElementOverflow(contextEL, mouseEvent);
		contextEL.style.left = `${x}px`;
		contextEL.style.top = `${y}px`;
	});

	function onDocumentClick(e: MouseEvent) {
		e.preventDefault();
		if (initialClick) {
			initialClick = false;
			return;
		}

		const target = e.target;
		if (target === contextEL) {
			return;
		}

		// avoid close context menu when clicked on context menu content
		if (target instanceof HTMLElement && isChildOfSlot(target)) {
			return;
		}

		onDisabled();
	}

	function isChildOfSlot(el: HTMLElement): boolean {
		if (contextEL === el) {
			return true;
		}
		while (el.parentNode && el.parentNode.nodeName.toLowerCase() !== 'body') {
			el = el.parentNode as any;
			if (contextEL === el) {
				return true;
			}
		}
		return false;
	}
</script>

<svelte:window on:click={onDocumentClick} />
<Portal target="body">
	<div class="context" bind:this={contextEL}>
		<slot />
	</div>
</Portal>

<style lang="scss">
	.context {
		background-color: var(--color-level-1);
		border-radius: 4px;
		border: 1px solid var(--color-border);
		font-size: 0.9rem;
		z-index: 9999;
		min-width: 104px;
		min-height: 32px;
		position: absolute;
		display: block;
		height: max-content;
	}
</style>

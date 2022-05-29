<script lang="ts">
	import { createBodyScrollToggler } from '$elven/tools';
	import { fade } from 'svelte/transition';
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import Portal from 'svelte-portal/src/Portal.svelte';

	export let onClose: (e: MouseEvent) => void;

	let toggleScroll: () => void;
	onMount(() => {
		if (!browser) {
			return;
		}
		toggleScroll = createBodyScrollToggler();
		toggleScroll();
	});

	onDestroy(() => {
		if (!browser) {
			return;
		}
		toggleScroll();
	});
</script>

<Portal target="body">
	<div class="overlay" transition:fade={{ duration: 120 }} on:click|stopPropagation|self={onClose}>
		<slot />
	</div>
</Portal>

<style lang="scss">
	.overlay {
		background-color: rgba(0, 0, 0, 0.4);
		z-index: 9998;
		max-width: 100vw;
		width: 100%;
		height: 100%;
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: grid;
	}
</style>

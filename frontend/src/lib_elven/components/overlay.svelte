<script lang="ts">
	import { toggleBodyScroll } from '$elven/tools';

	import { onDestroy, onMount } from 'svelte';

	export let onClose: (e: MouseEvent) => void;

	let setDefScroll: () => void;
	onMount(() => {
		setDefScroll = toggleBodyScroll();
	});

	onDestroy(() => {
		setDefScroll();
	});
</script>

<div class="overlay" on:click|stopPropagation|self={onClose}>
	<slot />
</div>

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

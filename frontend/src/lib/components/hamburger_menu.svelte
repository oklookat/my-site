<script lang="ts">
	import { browser } from '$app/env';

	import { createBodyScrollToggler } from '$elven/tools';

	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import { fly, fade } from 'svelte/transition';

	const dispatch = createEventDispatcher<{ closed: void }>();
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

<div
	class="hamburger"
	on:click|self|stopPropagation={() => dispatch('closed')}
	transition:fade={{ duration: 324 }}
>
	<div class="menu" transition:fly={{ x: -100, duration: 324 }}>
		<slot />
	</div>
</div>

<style lang="scss">
	.hamburger {
		overflow: hidden;
		background-color: rgba(0, 0, 0, 0.4);
		z-index: 9999;
		max-width: 100vw;
		width: 100%;
		height: 100%;
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		.menu {
			border-right: 1px solid var(--color-border);
			height: 100%;
			max-width: max-content;
		}
	}
</style>

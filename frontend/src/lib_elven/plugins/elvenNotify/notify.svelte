<script lang="ts">
	import { ElvenNotify } from '.';
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';

	let container: HTMLDivElement;
	let plugin: ElvenNotify | undefined;

	onMount(() => {
		if (!browser) {
			return;
		}
		plugin = new ElvenNotify(container, 2000);
	});

	onDestroy(() => {
		if (!browser || !plugin) {
			return;
		}
		plugin.destroy();
		plugin = undefined
	});
</script>

<div class="notify">
	<div class="notifications" bind:this={container} />
</div>

<style lang="scss">
	@mixin _desktop {
		@media (min-width: 765px) {
			@content;
		}
	}

	.notify {
		width: 100%;
		z-index: 9999;
		bottom: 0;
		margin-bottom: 8px;
		position: fixed;
		overflow: hidden;

		.notifications {
			width: 100%;
			gap: 8px;
			height: max-content;
			display: flex;
			flex-direction: column-reverse;
			box-sizing: border-box;
			position: relative;
		}
		
		@include _desktop() {
			margin-right: 12px;
			height: min-content;
			width: 224px;
			right: 0;
			bottom: 0;
			
			.notifications {
				height: max-content;
				width: max-content;
				flex-direction: column;
			}
		}
	}
</style>

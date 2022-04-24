<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import Notify from '$lib_elven/plugins/elvenNotify/notify.svelte';
	import Player from '$lib_elven/plugins/elvenPlayer/components/index.svelte';

	export let isAdmin = false;

	let isPlayerReady = false;
	let playerCore;

	onMount(async () => {
		const ElvenPlayer = await import('$lib_elven/plugins/elvenPlayer');
		const { ElvenChoose } = await import('$lib_elven/plugins/elvenChoose');
		new ElvenChoose();

		playerCore = new ElvenPlayer.default();
		isPlayerReady = true;
	});

	function destroyPlugins() {
		if (!playerCore) {
			return;
		}
		playerCore.destroy();
		playerCore = null;
	}

	onDestroy(() => {
		destroyPlugins();
	});
</script>

<div class="service">
	<Notify />

	{#if isPlayerReady && isAdmin}
		<Player core={playerCore} />
	{/if}
</div>

<style lang="scss">
	.service {
		display: flex;
		height: fit-content;
		background-color: red;
		bottom: 0;
		width: 100%;
		z-index: 7777;
		position: sticky;
		bottom: 0;
	}
</style>

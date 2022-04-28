<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import Notify from '$lib_elven/plugins/elvenNotify/notify.svelte';
	import Player from '$lib_elven/plugins/elvenPlayer/components/index.svelte';
	import Confirm from '$lib_elven/plugins/elvenChoose/confirm.svelte';

	export let isAdmin = false;

	let isPlayerReady = false;
	let elvenPlayer: any;

	onMount(async () => {
		const { default: ElvenPlayerModule } = await import('$lib_elven/plugins/elvenPlayer');
		elvenPlayer = new ElvenPlayerModule();
		isPlayerReady = true
	});

	function destroyPlugins() {
		if (!elvenPlayer) {
			return;
		}
		elvenPlayer.destroy();
		elvenPlayer = undefined;
		isPlayerReady = false;
	}

	onDestroy(() => {
		destroyPlugins();
	});
</script>

<div class="service">
	<Notify />
	<Confirm />

	{#if isPlayerReady && isAdmin}
		<Player core={elvenPlayer} />
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

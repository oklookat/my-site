<script lang="ts">
	import Notify from '$lib_elven/plugins/elvenNotify/notify.svelte';
	import Player from '$lib_elven/plugins/elvenPlayer/player.svelte';
	import Confirm from '$lib_elven/plugins/elvenChoose/confirm.svelte';

	export let isAdmin = false;
</script>

<div class="service">
	<Notify />
	<Confirm />

	{#if isAdmin}
		<Player />
	{/if}
</div>

<style lang="scss">
	.service {
		display: flex;
		height: fit-content;
		bottom: 0;
		width: 100%;
		z-index: 7777;
		position: sticky;
		bottom: 0;
	}
</style>

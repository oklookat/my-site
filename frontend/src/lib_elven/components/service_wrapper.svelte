<script lang="ts">
	import Notify from '$elven/plugins/elvenNotify/notify.svelte';
	import Player from '$elven/plugins/elvenPlayer/player.svelte';
	import Confirm from '$elven/plugins/elvenChoose/confirm.svelte';

	export let isAdmin = false;
</script>

<div class="service">
	<Notify />
	{#if isAdmin}
		<Confirm />
		<Player />
	{/if}
</div>

<style lang="scss">
	.service {
		display: flex;
		height: fit-content;
		width: 100%;
		z-index: 7777;
		position: sticky;
		bottom: 0;
	}
</style>

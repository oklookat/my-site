<script lang="ts">
	import ElvenLink from '$lib_elven/components/elven_link.svelte';
	import AccountSettings from '$lib_elven/components/settings_account.svelte';
	import ToolbarBig from '$lib_elven/components/toolbar_big.svelte';
	import Utils from '$lib_elven/tools';
</script>

<svelte:head>
	<title>{Utils.setTitleElven('settings')}</title>
</svelte:head>

<div class="settings">
	<ToolbarBig>
		<ElvenLink path="/logout">logout</ElvenLink>
	</ToolbarBig>

	<AccountSettings />
</div>

<style lang="scss">
	.settings {
		width: 95%;
		height: 100%;
		min-height: 64px;
		max-width: 512px;
		margin: auto;
		display: flex;
		flex-direction: column;
		gap: 14px;
	}
</style>

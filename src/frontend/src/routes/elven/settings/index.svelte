<script lang="ts">
	import ElvenLink from '$lib_elven/components/link.svelte';
	import AccountSettings from '$lib_elven/components/settings_account.svelte';
	import ToolbarBig from '$lib/components/toolbar_big.svelte';
	import { setTitleElven } from '$lib/tools';
	import { _ } from 'svelte-i18n'
</script>

<svelte:head>
	<title>{setTitleElven($_('elven.routes.settings.title'))}</title>
</svelte:head>

<div class="settings">
	<ToolbarBig>
		<ElvenLink path="/logout">{$_('elven.routes.settings.logout')}</ElvenLink>
	</ToolbarBig>

	<AccountSettings />
</div>

<style lang="scss">
	.settings {
		width: var(--width-95);
		height: 100%;
		min-height: 64px;
		max-width: 512px;
		margin: auto;
		display: flex;
		flex-direction: column;
		gap: 14px;
	}
</style>

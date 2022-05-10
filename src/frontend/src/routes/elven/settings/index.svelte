<script lang="ts">
	import AccountSettings from '$lib_elven/components/settings_account.svelte';
	import ToolbarBig from '$lib/components/toolbar_big.svelte';
	import { t } from '$lib/locale';
</script>

<svelte:head>
	<title>{$t('elven.settings.title')}</title>
</svelte:head>

<div class="settings">
	<ToolbarBig>
		<a href="/elven/logout">{$t('elven.settings.logout')}</a>
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

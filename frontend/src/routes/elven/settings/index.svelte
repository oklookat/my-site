<script context="module" lang="ts">
	export const load: Load = async (e) => {
		const stuff = e.stuff;
		const output = {
			stuff: stuff
		};
		stuff.title = t.get('elven.settings.title');
		return output;
	};
</script>

<script lang="ts">
	import AccountSettings from '$elven/components/settings_account.svelte';
	import ToolbarBig from '$elven/components/toolbar_big.svelte';
	import { t } from '$lib/locale';
	import type { Load } from '@sveltejs/kit';
</script>

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

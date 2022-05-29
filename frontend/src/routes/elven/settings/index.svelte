<script context="module" lang="ts">
	export const load: Load = async (e) => {
		const stuff = e.stuff;
		const output = {
			stuff: stuff
		};
		stuff.title = 'settings';
		return output;
	};
</script>

<script lang="ts">
	import AccountSettings from '$elven/components/settings_account.svelte';
	import Toolbar from '$lib/components/toolbar.svelte';

	import type { Load } from '@sveltejs/kit';
</script>

<div class="settings">
	<Toolbar>
		<a href="/elven/logout">log out</a>
	</Toolbar>

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

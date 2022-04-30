<script lang="ts">
	import { session } from '$app/stores';
	import OkLink from '$lib_oklookat/components/link.svelte';
	import ElvenLink from '$lib_elven/components/link.svelte';
	import { setTitleElven } from '$lib_elven/tools';
	import { _ } from 'svelte-i18n'
</script>

<svelte:head>
	<title>{setTitleElven($_('elven.routes.index.title'))}</title>
</svelte:head>

<div class="hello base__container">
	<div class="hi">{$_('elven.routes.index.hello')} {$session.user.username}</div>
	<div class="things">
		<OkLink path="/">{$_('elven.routes.index.mainSite')}</OkLink>
		<ElvenLink path="/settings">{$_('elven.routes.index.settings')}</ElvenLink>
	</div>
</div>

<style lang="scss">
	.hello {
		display: flex;
		flex-direction: column;
		flex-wrap: wrap;
		gap: 24px;
		height: 100%;
		border-radius: var(--border-radius);
		.hi {
			font-size: 2rem;
		}
		.things {
			display: flex;
			flex-direction: row;
			flex-wrap: wrap;
			gap: 14px;
			:global(*) {
				padding: 8px;
				background-color: var(--color-level-1);
				border-radius: var(--border-radius);
			}
		}
	}
</style>

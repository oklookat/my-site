<script lang="ts">
	import { session } from '$app/stores';
	import ElvenLink from '$lib_elven/components/link.svelte';
	import { setTitleElven } from '$lib/tools';
	import { _ } from 'svelte-i18n'
</script>

<svelte:head>
	<title>{setTitleElven($_('elven.routes.index.title'))}</title>
</svelte:head>

<div class="hello base__container">
	<div class="hi">{$_('elven.routes.index.hello')} {$session.user.username}</div>
	<div class="things">
		<a href="/">{$_('elven.routes.index.mainSite')}</a>
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
		}
	}
</style>

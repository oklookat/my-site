<script lang="ts">
	import { session } from '$app/stores';
	import { t } from '$lib/locale';
</script>

<svelte:head>
	<title>{$t('elven.index.title')}</title>
</svelte:head>

<div class="hello base__container">
	<div class="hi">{$t('elven.index.hello')} {$session.user.username}</div>
	<div class="things">
		<a class="card" href="/">{$t('elven.index.mainSite')}</a>
		<a class="card" href="/elven/settings">{$t('elven.index.settings')}</a>
	</div>
</div>

<style lang="scss">
	.card {
		height: 48px;
		width: 100%;
		background-color: var(--color-level-1);
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--border-radius);
	}
	.hello {
		display: flex;
		flex-direction: column;
		flex-wrap: wrap;
		gap: 24px;
		height: 100%;
		border-radius: var(--border-radius);
		.hi {
			font-size: 2rem;
			align-self: center;
		}
		.things {
			display: flex;
			flex-direction: row;
			flex-wrap: wrap;
			gap: 14px;
		}
	}
</style>

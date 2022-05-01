<script lang="ts">
	import { page } from '$app/stores';
	import ElvenLink from '$lib_elven/components/link.svelte';
	import { onDestroy } from 'svelte';
	import { _ } from 'svelte-i18n'
	// {$_('elven.header.elven')}

	let isUnknown = false;
	let isArticles = false;
	let isFiles = false;

	const unsub = page.subscribe((v) => {
		if (!v || !v.url || !v.url.pathname) {
			return;
		}
		onPathChanged(v.url.pathname);
	});

	onDestroy(() => {
		unsub();
	});

	function onPathChanged(path: string) {
		if (!path) {
			isUnknown = true;
			return;
		}
		isFiles = path.includes('/elven/files');
		isArticles = path.includes('/elven/articles');
		isUnknown = !isFiles && !isArticles;
	}
</script>

<nav class="header base__container">
	<div class="header__items">
		<ElvenLink path="">
			<div class={isUnknown ? 'route-active' : ''}>{$_('elven.components.header.elven')}</div>
		</ElvenLink>

		<ElvenLink path="/articles">
			<div class={isArticles ? 'route-active' : ''}>{$_('elven.components.header.articles')}</div>
		</ElvenLink>

		<ElvenLink path="/files">
			<div class={isFiles ? 'route-active' : ''}>{$_('elven.components.header.files')}</div>
		</ElvenLink>
	</div>
</nav>

<style lang="scss">
	:global(.route-active) {
		text-decoration: underline 1px;
	}

	.header {
		font-weight: bold;
		height: max-content;
		display: flex;
		align-items: center;
		&__items {
			padding: 14px;
			width: 100%;
			display: flex;
			align-items: center;
			flex-wrap: wrap;
			justify-content: center;
			gap: 42px;
		}
	}
</style>

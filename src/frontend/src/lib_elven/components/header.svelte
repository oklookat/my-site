<script lang="ts">
	import { page } from '$app/stores';
	//
	import ElvenLink from '$lib_elven/components/elven_link.svelte';

	let isUnknown = false;
	let isArticles = false;
	let isFiles = false;

	$: onPathChanged($page.url.pathname);
	function onPathChanged(path: string) {
		isFiles = path.includes('/elven/files');
		isArticles = path.includes('/elven/articles');
		isUnknown = !isFiles && !isArticles;
	}
</script>

<nav class="header">
	<div class="header__items">
		<ElvenLink path="">
			<div class={isUnknown ? 'route-active' : ''}>elven</div>
		</ElvenLink>

		<ElvenLink path="/articles">
			<div class={isArticles ? 'route-active' : ''}>articles</div>
		</ElvenLink>

		<ElvenLink path="/files">
			<div class={isFiles ? 'route-active' : ''}>files</div>
		</ElvenLink>
	</div>
</nav>

<style lang="scss">
	:global(.route-active) {
		text-decoration: underline 1px;
	}

	.header {
		font-weight: bold;
		color: var(--color-header-text);
		border-bottom: var(--color-border) 1px solid;
		height: max-content;
		width: 100%;
		display: flex;
		align-items: center;
		&__items {
			padding-top: 16px;
			padding-bottom: 16px;
			width: 100%;
			display: flex;
			align-items: center;
			flex-wrap: wrap;
			gap: 42px;
			margin-left: 2vw;
			margin-right: 2vw;
			@media screen and (max-width: 1023px) {
				justify-content: center;
				margin-left: 0;
				margin-right: 0;
			}
		}
	}
</style>

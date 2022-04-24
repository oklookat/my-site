<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	// ui
	import Pagination from '$lib_elven/components/pagination.svelte';
	// article
	import type { Article, Params } from '$lib_elven/types/articles';
	import type { Items } from '$lib_elven/types';
	import ArticlesToolbars from '$lib_elven/components/articles_toolbars.svelte';
	import ArticlesList from '$lib_elven/components/articles_list.svelte';
	import NetworkArticle from '$lib_elven/network/network_article';
	import { setTitleElven } from '$lib_elven/tools';
	import { HandleRouteParam, Refresh } from '$lib_elven/tools/routes';

	const networkArticle = new NetworkArticle('');

	/** articles data */
	export let items: Items<Article>;

	/** request params */
	export let params: Params;

	/** url searchparams */
	let searchparams = $page.url.searchParams;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: 'page', val: page });
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string | boolean | number }) {
		const data = await HandleRouteParam<Article>(networkArticle, event, {
			items,
			params,
			searchparams
		});
		items = data.items;
		params = data.params;
		searchparams = data.searchparams;
	}

	async function refresh() {
		const data = await Refresh<Article>(networkArticle, {
			items,
			params,
			searchparams
		});
		items = data.items;
		params = data.params;
		searchparams = data.searchparams;
	}
</script>

<svelte:head>
	<title>{setTitleElven('articles')}</title>
</svelte:head>

<div class="articles base__container">
	<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />

	<ArticlesList {items} onDeleted={async () => await refresh()} />

	<div class="pages">
		{#if items.meta}
			<Pagination
				bind:total={items.meta.total_pages}
				bind:current={items.meta.current_page}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>

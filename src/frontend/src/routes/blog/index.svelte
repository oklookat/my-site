<script context="module" lang="ts">
	import NetworkArticle from '$lib_elven/network/network_article';

	import { HandleRouteParam, Params, type RPH_Data, type RPH_Event } from '$lib_elven/tools/params';
	import type { Items } from '$lib_elven/types';
	import type { Article } from '$lib_elven/types/articles';

	import type { Load } from '@sveltejs/kit';

	export const load: Load = async (event) => {
		let requestParams = new Params<Article>('article', event.url.searchParams);
		requestParams.setParam('drafts', false);
		requestParams.setParam('by', 'published');

		let response: Response;

		let items: Items<Article>;

		const networkArticle = new NetworkArticle('');

		const fetchData = async () => {
			response = await networkArticle.getAll(requestParams.toObject(), event.fetch);
			if (response.ok) {
				items = (await response.json()) as Items<Article>;
				return;
			}
			throw Error(response.statusText);
		};

		try {
			await fetchData();
			const pageParam = requestParams.getParam('page');
			if (pageParam > items.meta.total_pages) {
				// @ts-ignore
				requestParams.setParam('page', items.meta.total_pages);
				await fetchData();
			}
		} catch (err) {}

		return {
			status: 200,
			props: {
				items: items,
				params: requestParams
			}
		};
	};
</script>

<script lang="ts">
	import { _ } from 'svelte-i18n';
	import Pagination from '$lib_elven/components/pagination.svelte';
	import ArticlesToolbars from '$lib_oklookat/components/articles_toolbars.svelte';
	import ArticlesList from '$lib_oklookat/components/articles_list.svelte';

	/** articles data */
	export let items: Items<Article>;

	/** request params */
	export let params: Params<Article>;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: 'page', val: page });
	}

	/** on request param changed */
	async function onParamChanged(event: RPH_Event<Article>) {
		const data: RPH_Data<Article> = {
			params: params,
			items: items
		};
		await HandleRouteParam<Article>(event, data);
	}
</script>

<svelte:head>
	<title>блог - oklookat</title>
</svelte:head>

<div class="articles base__container">
	<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />

	<ArticlesList bind:items />

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

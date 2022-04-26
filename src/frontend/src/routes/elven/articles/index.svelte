<script context="module" lang="ts">
	export const load: Load = async (event) => {
		let requestParams = new Params<Article>('article', event.url.searchParams);

		const setParam = (name: string, val: any) => {
			// @ts-ignore
			requestParams.setParam(name, val);
		};

		let response: Response;
		let items: Items<Article>;
		const networkArticle = new NetworkArticle(event.session.user.token);

		const fetchData = async () => {
			response = await networkArticle.getAll(requestParams.toObject(), event.fetch);
			if (response.ok) {
				items = (await response.json()) as Items<Article>;
			} else {
				throw Error(response.statusText);
			}
		};

		try {
			await fetchData();
			if (requestParams.getParam('page') > items.meta.total_pages) {
				setParam('page', items.meta.total_pages);
				await fetchData();
			}
		} catch (err) {}

		return {
			status: response.status,
			props: {
				items: items,
				params: requestParams
			}
		};
	};
</script>

<script lang="ts">
	// ui
	import Pagination from '$lib_elven/components/pagination.svelte';
	// article
	import type { Article } from '$lib_elven/types/articles';
	import type { Items } from '$lib_elven/types';
	import ArticlesToolbars from '$lib_elven/components/articles_toolbars.svelte';
	import ArticlesList from '$lib_elven/components/articles_list.svelte';
	import NetworkArticle from '$lib_elven/network/network_article';
	import { setTitleElven } from '$lib_elven/tools';
	import { HandleRouteParam, Refresh } from '$lib_elven/tools/routes';
	import { Params } from '$lib_elven/tools/params';
	import type { Load } from '@sveltejs/kit';

	const networkArticle = new NetworkArticle('');

	/** articles data */
	export let items: Items<Article>;

	/** request params */
	export let params: Params<Article>;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: 'page', val: page });
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string | boolean | number }) {
		const data = await HandleRouteParam<Article>(networkArticle, event, {
			items,
			params
		});
		items = data.items;
		params = data.params;
	}

	async function refresh() {
		const data = await Refresh<Article>(networkArticle, {
			items,
			params
		});
		items = data.items;
		params = data.params;
	}
</script>

<svelte:head>
	<title>{setTitleElven('articles')}</title>
</svelte:head>

<div class="articles base__container">
	<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />

	<ArticlesList bind:items on:deleted={async () => await refresh()} />

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

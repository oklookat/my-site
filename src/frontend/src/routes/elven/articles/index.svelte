<script context="module" lang="ts">
	export const load: Load = async (e) => {
		let requestParams = new Params<Article>('article', e.url.searchParams);
		let items: Items<Article>;

		let resp: Response | null = null;

		const stuff = e.stuff;
		stuff.title = t.get('elven.articles.title');

		const networkArticle = new NetworkArticle(getTokenFromSession(e), e.fetch);
		try {
			resp = await networkArticle.getAll(requestParams.toObject());
			if (resp.ok) {
				items = (await resp.json()) as Items<Article>;
			}
		} catch (err) {
			throw Error(resp?.statusText);
		}

		return {
			status: resp?.status || 200,
			stuff: stuff,
			props: {
				items: items,
				params: requestParams
			}
		};
	};
</script>

<script lang="ts">
	// ui
	import Pagination from '$lib/components/pagination.svelte';
	// article
	import type { Article } from '$lib/types/articles';
	import type { Items } from '$lib/types';
	import ArticlesToolbars from '$lib/components/elven/articles_toolbars.svelte';
	import ArticlesList from '$lib/components/elven/articles_list.svelte';
	import NetworkArticle from '$lib/network/article';
	import {
		HandleRouteParam,
		Params,
		Refresh,
		type RPH_Data,
		type RPH_Event
	} from '$lib/tools/params';
	import type { Load } from '@sveltejs/kit';
	import { t } from '$lib/locale';
	import { getTokenFromSession } from '$lib/tools';

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

	async function refresh() {
		const getData = async () => {
			return Promise.resolve({
				items: items,
				params: params
			});
		};
		await Refresh<Article>(getData);
	}
</script>

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

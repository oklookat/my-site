<script context="module" lang="ts">
	export const load: Load = async (e) => {
		let requestParams = new Params<RAW>('article', e.url.searchParams);
		let items: Items<RAW>;

		let resp: Response | null = null;

		const stuff = e.stuff;
		stuff.title = 'articles';

		const networkArticle = new NetworkArticle(getTokenFromSession(e), e.fetch);
		try {
			resp = await networkArticle.getAll(requestParams.toObject());
			if (resp.ok) {
				items = (await resp.json()) as Items<RAW>;
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
	import Pagination from '$lib/components/pagination.svelte';
	import type { RAW } from '$elven/types/article';
	import type { Items } from '$elven/types';
	import ArticlesToolbars from '$elven/components/articles_toolbars.svelte';
	import ArticlesList from '$elven/components/articles_list.svelte';
	import NetworkArticle from '$elven/network/article';
	import {
		HandleRouteParam,
		Params,
		Refresh,
		type RPH_Data,
		type RPH_Event
	} from '$elven/tools/params';
	import type { Load } from '@sveltejs/kit';
	import { getTokenFromSession } from '$elven/tools';
	import ItemsContainer from '$elven/components/items_container.svelte';

	/** articles data */
	export let items: Items<RAW>;

	/** request params */
	export let params: Params<RAW>;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: 'page', val: page });
	}

	/** on request param changed */
	async function onParamChanged(event: RPH_Event<RAW>) {
		const data: RPH_Data<RAW> = {
			params: params,
			items: items
		};
		await HandleRouteParam<RAW>(event, data);
	}

	async function refresh() {
		const getData = async () => {
			return Promise.resolve({
				items: items,
				params: params
			});
		};
		await Refresh<RAW>(getData);
	}
</script>

<ItemsContainer>
	<div slot="up">
		<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />
	</div>
	
	<div slot="list">
		<ArticlesList bind:items on:deleted={async () => await refresh()} />
	</div>

	<div slot="pages">
		{#if items.meta}
			<Pagination
				bind:total={items.meta.total_pages}
				bind:current={items.meta.current_page}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</ItemsContainer>

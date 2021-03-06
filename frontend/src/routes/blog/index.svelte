<script context="module" lang="ts">
	import NetworkArticle from '$elven/network/article';
	import { HandleRouteParam, Params, type RPH_Data, type RPH_Event } from '$elven/tools/params';
	import type { Items } from '$elven/types';
	import type { RAW } from '$elven/types/article';

	import type { Load, LoadOutput } from '@sveltejs/kit';

	export const load: Load = async (e) => {
		const requestParams = new Params<RAW>('article', e.url.searchParams);
		requestParams.setParam('drafts', false);
		requestParams.setParam('by', 'published');

		let resp: Response;
		const networkArticle = new NetworkArticle('', e.fetch);

		const stuff = e.stuff;
		const output: LoadOutput = {
			status: 200,
			stuff: stuff,
			props: {
				items: null,
				params: requestParams
			}
		};

		stuff.title = 'блог';

		try {
			resp = await networkArticle.getAll(requestParams.toObject());
			output.status = resp.status;
			if (resp.ok && output.props) {
				output.props.items = (await resp.json()) as Items<RAW>;
			}
		} catch (err) {}

		return output;
	};
</script>

<script lang="ts">
	import Pagination from '$lib/components/pagination.svelte';
	import ArticlesToolbars from '$oklookat/components/articles_toolbars.svelte';
	import ArticlesList from '$oklookat/components/articles_list.svelte';
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
</script>

<ItemsContainer>
	<div slot="up">
		<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />
	</div>

	<div slot="list">
		<ArticlesList bind:items />
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

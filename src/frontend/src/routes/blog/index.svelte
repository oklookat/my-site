<script context="module" lang="ts">
	import NetworkArticle from '$lib_elven/network/network_article';
	import { HandleRouteParam, Params, type RPH_Data, type RPH_Event } from '$lib_elven/tools/params';
	import type { Items } from '$lib_elven/types';
	import type { Article } from '$lib_elven/types/articles';

	import type { Load, LoadOutput } from '@sveltejs/kit';

	export const load: Load = async (e) => {
		const requestParams = new Params<Article>('article', e.url.searchParams);
		requestParams.setParam('drafts', false);
		requestParams.setParam('by', 'published');

		let resp: Response;
		const networkArticle = new NetworkArticle('', e.fetch);

		const stuff = e.stuff
		const output: LoadOutput = {
			status: 200,
			stuff: stuff,
			props: {
				items: null,
				params: requestParams
			}
		}

		stuff.title = "блог"

		try {
			resp = await networkArticle.getAll(requestParams.toObject());
			output.status = resp.status
			if (resp.ok && output.props) {
				output.props.items = (await resp.json()) as Items<Article>;
			}
		} catch (err) {}

		return output
	};
</script>

<script lang="ts">
	import { t } from '$lib/locale';
	import Pagination from '$lib/components/pagination.svelte';
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

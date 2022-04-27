<script context="module" lang="ts">
	export const load: Load = async (event) => {
		let requestParams = new Params<Article>(
			"article",
			event.url.searchParams
		);

		const setParam = (name: string, val: any) => {
			// @ts-ignore
			requestParams.setParam(name, val);
		};

		let response: Response;

		let items: Items<Article>;

		const networkArticle = new NetworkArticle(event.session.user.token);

		const fetchData = async () => {
			response = await networkArticle.getAll(
				requestParams.toObject(),
				event.fetch
			);
			if (response.ok) {
				items = (await response.json()) as Items<Article>;
				return;
			}
			throw Error(response.statusText);
		};

		try {
			await fetchData();
			const pageParam = requestParams.getParam("page");
			if (pageParam > items.meta.total_pages) {
				setParam("page", items.meta.total_pages);
				await fetchData();
			}
		} catch (err) {}

		return {
			status: response.status,
			props: {
				items: items,
				params: requestParams,
			},
		};
	};
</script>

<script lang="ts">
	// ui
	import Pagination from "$lib_elven/components/pagination.svelte";
	// article
	import type { Article } from "$lib_elven/types/articles";
	import type { Items } from "$lib_elven/types";
	import ArticlesToolbars from "$lib_elven/components/articles_toolbars.svelte";
	import ArticlesList from "$lib_elven/components/articles_list.svelte";
	import NetworkArticle from "$lib_elven/network/network_article";
	import { setTitleElven } from "$lib_elven/tools";
	import {
		HandleRouteParam,
		Params,
		Refresh,
		type RPH_Data,
		type RPH_Event,
	} from "$lib_elven/tools/params";
	import type { Load } from "@sveltejs/kit";

	/** articles data */
	export let items: Items<Article>;

	/** request params */
	export let params: Params<Article>;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: "page", val: page });
	}

	/** on request param changed */
	async function onParamChanged(event: RPH_Event<Article>) {
		const data: RPH_Data<Article> = {
			params: params,
			items: items,
		};
		await HandleRouteParam<Article>(event, data);
	}

	async function refresh() {
		const getData = async () => {
			return Promise.resolve({
				items: items,
				params: params
			})
		}
		await Refresh<Article>(getData);
	}
</script>

<svelte:head>
	<title>{setTitleElven("articles")}</title>
</svelte:head>

<div class="articles base__container">
	<ArticlesToolbars
		bind:params
		on:paramChanged={async (e) => await onParamChanged(e.detail)}
	/>

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

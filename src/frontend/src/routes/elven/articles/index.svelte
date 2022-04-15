<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	// ui
	import Pagination from '$lib_elven/components/pagination.svelte';
	// utils
	import Utils from '$lib_elven/tools';
	// article
	import type { Article, Params } from '$lib_elven/types/articles';
	import type { Data } from '$lib_elven/types';
	import ArticlesToolbars from '$lib_elven/components/articles_toolbars.svelte';
	import ArticlesList from '$lib_elven/components/articles_list.svelte';
	import NetworkArticle from '$lib_elven/network/network_article';

	const networkArticle = new NetworkArticle('')

	/** articles data */
	export let items: Data<Article>;

	/** request params */
	export let params: Params;

	/** url searchparams */
	const urlParams = $page.url.searchParams;

	onMount(async () => {
		// @ts-ignore
		await goto(`?${urlParams.toString()}`, { replaceState: true });
	});

	/** when page changed */
	async function onPageChanged(page: number) {
		// refresh if page not changed
		if (page === params.page) {
			const resp = await networkArticle.getAll(params);
			if (resp.status === 200) {
				items = await resp.json();
			}
			return;
		}
		params.page = page;
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	async function refresh() {
		const getPage = async () => {
			return params.page;
		};
		const setPage = async (newPage: number) => {
			params.page = newPage;
		};
		const fetchItems = async (initial: boolean) => {
			if(initial) {
				return items
			}
			await onPageChanged(await getPage());
			return items;
		};
		await Utils.refresh(getPage, setPage, fetchItems);
	}
</script>

<svelte:head>
	<title>elven: articles</title>
</svelte:head>

<div class="articles base__container">
	<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />

	<ArticlesList {items} onDeleted={async () => await refresh()} />

	<div class="pages">
		{#if items.meta}
			<Pagination
				total={items.meta.total_pages}
				current={items.meta.current_page}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>

<style lang="scss">
</style>

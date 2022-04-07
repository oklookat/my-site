<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	// ui
	import Pagination from '$lib/components/pagination.svelte';
	// article
	import { getDefaultParams, type Article, type Params } from '$lib/types/articles';
	import type { Data, Meta } from '$lib/types';
	import CArticle from '$lib/components/article.svelte';
	import ArticleActions from '$lib/components/article_actions.svelte';
	import ArticlesToolbars from '$lib/components/articles_toolbars.svelte';
import Utils from '$lib/tools';

	/** is article selected? */
	let isSelected = false;

	/** selected article */
	let selected: {
		counter: number | null;
		article: Article | null;
		mouseEvent: MouseEvent;
	} = { counter: null, article: null, mouseEvent: null };

	export let items: Data<Article>;
	export let params: Params;

	const urlParams = $page.url.searchParams;
	onMount(async () => {
    // @ts-ignore
    Utils.searchParamsByObject($page.url.searchParams, getDefaultParams())
		await goto(`?${urlParams.toString()}`, { replaceState: true });
	});

	function onDeleted() {
		isSelected = false;
		deleteFromArray(selected.counter);
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`);
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);
		await goto(`?${urlParams.toString()}`);
	}

	/** select article */
	function select(article: Article, mouseEvent: MouseEvent, counter: number) {
		selected.counter = counter;
		selected.mouseEvent = mouseEvent;
		selected.article = items.data[counter];
		isSelected = true;
	}

	/** delete article from articles array and refresh articles */
	async function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
		await refresh();
	}

	/** refresh articles */
	async function refresh() {
		// const getData = async () => {
		//   await getAll()
		//   return articles
		// }
		// const setPage = (val: number) => (requestParams.page = val);
		// await Utils.refresh(requestParams.page, setPage, getData);
	}
</script>

<svelte:head>
	<title>elven: articles</title>
</svelte:head>

{#if isSelected}
	<ArticleActions
		article={selected.article}
		mouseEvent={selected.mouseEvent}
		onDisabled={() => (isSelected = false)}
		onDeleted={() => onDeleted()}
	/>
{/if}

<div class="articles base__container">
	<ArticlesToolbars bind:params on:paramChanged={async (e) => await onParamChanged(e.detail)} />

	<div class="list">
		{#if items.data}
			{#each Object.entries(items.data) as [counter, article]}
				<CArticle
					{article}
					onSelected={(article, event) => select(article, event, parseInt(counter, 10))}
				/>
			{/each}
		{/if}
	</div>

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

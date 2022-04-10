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
	import CArticle from '$lib_elven/components/article.svelte';
	import ArticleActions from '$lib_elven/components/article_actions.svelte';
	import ArticlesToolbars from '$lib_elven/components/articles_toolbars.svelte';

	/** is article selected? */
	let isSelected = false;

	/** selected article */
	let selected: {
		counter: number | null;
		article: Article | null;
		mouseEvent: MouseEvent;
	} = { counter: null, article: null, mouseEvent: null };

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

	function onDeleted() {
		isSelected = false;
		deleteFromArray(selected.counter);
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	/** select article */
	function select(mouseEvent: MouseEvent, counter: number) {
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
		const getData = async () => {
			await onPageChanged(params.page);
			return items.data;
		};
		const setPage = (val: number) => (params.page = val);
		await Utils.refresh(params.page, setPage, getData);
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
				<CArticle {article} onSelected={(article, event) => select(event, parseInt(counter, 10))} />
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

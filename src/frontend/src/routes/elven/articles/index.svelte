<script lang="ts">
	import { onMount } from 'svelte';
	import { goto, invalidate } from '$app/navigation';
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
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}


	/** refresh files */
	const refresh = getRefresher();
	function getRefresher() {
		let force = false;
		let prevPage = params.page;
		let isFirstCall = true;
		const getData = async () => {
			if (isFirstCall && !force) {
				isFirstCall = false;
				return items.data;
			}
			if (prevPage < 2 && params.page < 2) {
				params.page = 1;
				await invalidate('');
				return items.data;
			}
			await onPageChanged(params.page);
			return items.data;
		};
		const setPage = (val: number) => {
			prevPage = params.page;
			params.page = val;
		};
		return async (isForce = false) => {
			force = isForce;
			await Utils.refresh(params.page, setPage, getData);
		};
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

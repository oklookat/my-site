<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	// tools
	import type { Data } from '$lib/types';
	import Utils from '$lib/tools';
	// ui
	import Pagination from '$lib/components/pagination.svelte';
	// file
	import { getDefaultParams, type File, type Params } from '$lib/types/files';
	import CFile from '$lib/components/file.svelte';
	import FileActions from '$lib/components/file_actions.svelte';
	import FilesToolbars from '$lib/components/files_toolbars.svelte';
	import NetworkFile from '$lib/network/network_file';

	const networkFile = new NetworkFile('');

	export let isSelectMode = false;

	const dispatch = createEventDispatcher<{
		/** on 'select' option clicked on file */
		selected: File;
	}>();

	/** selected file */
	let selected: {
		counter: number | null;
		file: File | null;
		mouseEvent: MouseEvent;
	} = { counter: null, file: null, mouseEvent: null };

	/** files data */
	export let items: Data<File> | undefined = undefined;

	/** request params */
	export let params: Params;

	/** url searchparams */
	let urlParams = $page.url.searchParams;

	onMount(async () => {
		// @ts-ignore
		Utils.searchParamsByObject(urlParams, getDefaultParams());
		if (!items) {
			await getAll();
			return;
		}
		await goto(`?${urlParams.toString()}`, { replaceState: true });
	});

	/** get files without SSR */
	async function getAll(requestParams = params) {
		items = await networkFile.getAll(requestParams);
	}

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		if (isSelectMode) {
			params.page = 1;
			params[event.name] = event.val;
			await getAll();
			return;
		}
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);
		await goto(`?${urlParams.toString()}`);
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		if (isSelectMode) {
			params.page = page;
			await getAll();
			return;
		}
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`);
	}

	/** is file selected? */
	let isSelected = false;

	/** on selected file deleted */
	function onDeleted() {
		isSelected = false;
		deleteFromArray(selected.counter);
	}

	/** select file */
	function select(mouseEvent: MouseEvent, counter: number) {
		selected.counter = counter;
		selected.mouseEvent = mouseEvent;
		selected.file = items.data[counter];
		isSelected = true;
	}

	/** refresh files */
	async function refresh() {
		const getData = async () => {
			await onPageChanged(params.page);
			return items.data;
		};
		const setPage = (val: number) => (params.page = val);
		await Utils.refresh(params.page, setPage, getData);
	}

	/** delete file from files array */
	async function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
		await refresh();
	}

	/** on 'select' button clicked on selected file */
	function onSelectClicked() {
		if (!isSelectMode) {
			return;
		}
		dispatch('selected', selected.file);
	}
</script>

<svelte:head>
	<title>elven: files</title>
</svelte:head>

{#if isSelected}
	<FileActions
		file={selected.file}
		mouseEvent={selected.mouseEvent}
		onDisabled={() => (isSelected = false)}
		onDeleted={() => onDeleted()}
		{onSelectClicked}
	/>
{/if}

<div class="files base__container">
	<FilesToolbars bind:params on:paramChanged={async (e) => onParamChanged(e.detail)} />

	<div class="list">
		{#if items && items.data}
			{#each Object.entries(items.data) as [counter, file]}
				<CFile {file} onSelected={(e) => select(e, parseInt(counter, 10))} />
			{/each}
		{/if}
	</div>

	<div class="pages">
		{#if items && items.meta}
			<Pagination
				total={items.meta.total_pages}
				current={items.meta.current_page}
				on:changed={(e) => onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>

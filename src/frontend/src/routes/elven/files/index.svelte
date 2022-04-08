<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	// tools
	import type { Data, Meta } from '$lib/types';
	import Utils from '$lib/tools';
	// ui
	import Pagination from '$lib/components/pagination.svelte';
	// file
	import { getDefaultParams, type File, type Params } from '$lib/types/files';
	import CFile from '$lib/components/file.svelte';
	import FileActions from '$lib/components/file_actions.svelte';
	import FilesToolbars from '$lib/components/files_toolbars.svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	/** add "select" option to selection overlay and dispatch event if this button clicked */
	export let withSelect: boolean = false;

	const dispatch = createEventDispatcher<{
		/** on 'select' option clicked on file */
		selected: File;
	}>();

	/** files loaded? */
	let loaded = false;

	/** response information */
	let meta: Meta;

	/** is file selected? */
	let isSelected = false;

	/** selected file */
	let selected: {
		counter: number | null;
		file: File | null;
		mouseEvent: MouseEvent;
	} = { counter: null, file: null, mouseEvent: null };

	/** files data */
	export let items: Data<File>;

	/** request params */
	export let params: Params;

	/** url searchparams */
	const urlParams = $page.url.searchParams;

	onMount(async () => {
		// @ts-ignore
		Utils.searchParamsByObject($page.url.searchParams, getDefaultParams());
		await goto(`?${urlParams.toString()}`, { replaceState: true });
	});

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);
		await goto(`?${urlParams.toString()}`);
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`);
	}

	/** on selected file deleted */
	function onDeleted() {
		isSelected = false;
		deleteFromArray(selected.counter);
	}

	/** select file */
	function select(file: File, mouseEvent: MouseEvent, counter: number) {
		selected.counter = counter;
		selected.mouseEvent = mouseEvent;
		selected.file = items.data[counter];
		isSelected = true;
	}

	/** refresh files */
	async function refresh() {
		// const getData = async () => {
		//   await getAll();
		//   return files;
		// };
		// const setPage = (val: number) => (requestParams.page = val);
		// await Utils.refresh(requestParams.page, setPage, getData);
	}

	/** delete file from files array */
	async function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
		await refresh();
	}

	/** on 'select' button clicked on selected file */
	function onSelectClicked() {
		if (!withSelect) {
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
		{withSelect}
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
		{#if loaded && Utils.getRecordLength(items.data) > 0}
			{#each Object.entries(items.data) as [counter, file]}
				<CFile {file} onSelected={(e) => select(file, e, parseInt(counter, 10))} />
			{/each}
		{/if}
	</div>

	<div class="pages">
		{#if loaded && meta && meta.total_pages && meta.current_page}
			<Pagination
				total={meta.total_pages}
				current={meta.current_page}
				on:changed={(e) => onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>

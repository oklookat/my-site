<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto, invalidate } from '$app/navigation';
	// tools
	import type { Data } from '$lib_elven/types';
	import Utils from '$lib_elven/tools';
	// ui
	import Pagination from '$lib_elven/components/pagination.svelte';
	// file
	import type { File, Params } from '$lib_elven/types/files';
	import FileActions from '$lib_elven/components/file_actions.svelte';
	import FilesToolbars from '$lib_elven/components/files_toolbars.svelte';
	import FilesList from '$lib_elven/components/files_list.svelte';

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
		await goto(`?${urlParams.toString()}`, { replaceState: true });
	});

	/** on request param changed */
	async function onParamChanged(event: { name: string; val: string }) {
		urlParams.set('page', '1');
		urlParams.set(event.name, event.val);

		// remove filename param if empty
		if (urlParams.has('filename') && !urlParams.get('filename')) {
			urlParams.delete('filename');
		}

		// keepfocus for search/page input
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	/** when page changed */
	async function onPageChanged(page: number) {
		urlParams.set('page', page.toString());
		await goto(`?${urlParams.toString()}`, { keepfocus: true });
	}

	/** is file selected? */
	let isSelected = false;

	/** on selected file deleted */
	async function onDeleted() {
		isSelected = false;
		deleteFromArray(selected.counter);
		await refresh();
	}

	/** delete file from files array */
	function deleteFromArray(counter: number) {
		delete items.data[counter];
		items = items;
	}

	/** select file */
	function select(mouseEvent: MouseEvent, counter: number) {
		selected.counter = counter;
		selected.mouseEvent = mouseEvent;
		selected.file = items.data[counter];
		isSelected = true;
	}

	/** refresh files */
	const refresh = getRefresher()
	function getRefresher() {
		let force = false
		let prevPage = params.page
		let isFirstCall = true
		const getData = async () => {
			if(isFirstCall && !force) {
				isFirstCall = false
				return items.data
			}
			if(prevPage < 2 && params.page < 2) {
				params.page = 1
				await invalidate("");
				return items.data
			}
			await onPageChanged(params.page);
			return items.data;
		};
		const setPage = (val: number) => {
			prevPage = params.page
			params.page = val
		};
		return async (isForce = false) => {
			force = isForce
			await Utils.refresh(params.page, setPage, getData);
		}
	}

	/** on file uploaded */
	async function onUploaded() {
		params.page = 1;
		await refresh(true);
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
		onDeleted={async () => await onDeleted()}
	/>
{/if}

<div class="files base__container">
	<FilesToolbars
		bind:params
		on:uploaded={async () => await onUploaded()}
		on:paramChanged={async (e) => onParamChanged(e.detail)}
	/>

	<FilesList {items} onSelected={(file, counter, e) => select(e, counter)} />

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

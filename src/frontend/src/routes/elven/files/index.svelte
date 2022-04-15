<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	// tools
	import type { Data } from '$lib_elven/types';
	// ui
	import Pagination from '$lib_elven/components/pagination.svelte';
	// file
	import type { File, Params } from '$lib_elven/types/files';
	import FilesToolbars from '$lib_elven/components/files_toolbars.svelte';
	import FilesList from '$lib_elven/components/files_list.svelte';
	import Utils from '$lib_elven/tools';
	import NetworkFile from '$lib_elven/network/network_file';

	const networkFile = new NetworkFile('');

	/** files data */
	export let items: Data<File> | undefined = undefined;

	/** request params */
	export let params: Params;

	/** url searchparams */
	let urlParams = $page.url.searchParams;

	onMount(async () => {
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
		// refresh if page not changed
		if (page === params.page) {
			const resp = await networkFile.getAll(params);
			if (resp.status === 200) {
				items = await resp.json();
			}
			return;
		}
		params.page = page;
		urlParams.set('page', page.toString());
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
			if (initial) {
				return items;
			}
			await onPageChanged(await getPage());
			return items;
		};
		await Utils.refresh(getPage, setPage, fetchItems);
	}
</script>

<svelte:head>
	<title>elven: files</title>
</svelte:head>

<div class="files base__container">
	<FilesToolbars bind:params {refresh} on:paramChanged={async (e) => onParamChanged(e.detail)} />

	<FilesList {items} onDeleted={async () => await refresh()} />

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

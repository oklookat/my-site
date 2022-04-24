<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import type { Items } from '$lib_elven/types';
	import Pagination from '$lib_elven/components/pagination.svelte';
	import type { File, Params } from '$lib_elven/types/files';
	import FilesToolbars from '$lib_elven/components/files_toolbars.svelte';
	import FilesList from '$lib_elven/components/files_list.svelte';
	import NetworkFile from '$lib_elven/network/network_file';
	import { setTitleElven } from '$lib_elven/tools';
	import { HandleRouteParam, Refresh, type RPH_Event } from '$lib_elven/tools/routes';

	const networkFile = new NetworkFile('');

	/** files data */
	export let items: Items<File> | undefined = undefined;

	/** request params */
	export let params: Params;

	/** url searchparams */
	let searchparams = $page.url.searchParams;

	async function onPageChanged(page: number) {
		await onParamChanged({ name: 'page', val: page });
	}

	async function onUploaded() {
		await onParamChanged({ name: 'page', val: 1 });
	}

	/** on request param changed */
	async function onParamChanged(event: RPH_Event) {
		const data = await HandleRouteParam<File>(networkFile, event, {
			items,
			params,
			searchparams
		});
		items = data.items;
		params = data.params;
		searchparams = data.searchparams;
	}

	async function refresh() {
		const data = await Refresh(networkFile, {
			items,
			params,
			searchparams
		});
		items = data.items;
		params = data.params;
		searchparams = data.searchparams;
	}
</script>

<svelte:head>
	<title>{setTitleElven('files')}</title>
</svelte:head>

<div class="files base__container">
	<FilesToolbars
		bind:params
		on:uploaded={async () => await onUploaded()}
		on:paramChanged={async (e) => await onParamChanged(e.detail)}
	/>

	<FilesList bind:items onDeleted={async () => await refresh()} />

	<div class="pages">
		{#if items && items.meta}
			<Pagination
				total={items.meta.total_pages}
				current={params.page}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>

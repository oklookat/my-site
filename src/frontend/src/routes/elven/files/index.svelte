<script context="module" lang="ts">
	export const load: Load = async (event) => {
		let requestParams = new Params<File>('file', event.url.searchParams);

		const setParam = (name: string, val: any) => {
			// @ts-ignore
			requestParams.setParam(name, val);
		};

		let response: Response;
		let items: Items<File>;
		const networkFile = new NetworkFile(event.session.user.token);

		const fetchData = async () => {
			response = await networkFile.getAll(requestParams.toObject(), event.fetch);
			if (response.ok) {
				items = (await response.json()) as Items<File>;
			} else {
				throw Error(response.statusText);
			}
		};

		try {
			await fetchData();
			if (requestParams.getParam('page') > items.meta.total_pages) {
				setParam('page', items.meta.total_pages);
				await fetchData();
			}
		} catch (err) {}

		return {
			status: response.status,
			props: {
				items: items,
				params: requestParams
			}
		};
	};
</script>

<script lang="ts">
	import type { Items } from '$lib_elven/types';
	import Pagination from '$lib_elven/components/pagination.svelte';
	import type { File } from '$lib_elven/types/files';
	import FilesToolbars from '$lib_elven/components/files_toolbars.svelte';
	import FilesList from '$lib_elven/components/files_list.svelte';
	import NetworkFile from '$lib_elven/network/network_file';
	import { setTitleElven } from '$lib_elven/tools';
	import { HandleRouteParam, Refresh, type RPH_Event } from '$lib_elven/tools/routes';
	import { Params } from '$lib_elven/tools/params';
	import type { Load } from '@sveltejs/kit';

	const networkFile = new NetworkFile('');

	/** files data */
	export let items: Items<File>;

	/** request params */
	export let params: Params<File>;

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
			params
		});
		items = data.items;
		params = data.params;
	}

	async function refresh() {
		const data = await Refresh(networkFile, {
			items,
			params
		});
		items = data.items;
		params = data.params;
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

	<FilesList {items} on:deleted={async () => await refresh()} />

	<div class="pages">
		{#if items && items.meta}
			<Pagination
				total={items.meta.total_pages}
				current={params.getParam('page')}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</div>

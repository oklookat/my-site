<script context="module" lang="ts">
	export const load: Load = async (e) => {
		let requestParams = new Params<File>('file', e.url.searchParams);

		let resp: Response | null = null;
		let items: Items<File>;

		const stuff = e.stuff;
		stuff.title = 'files';

		const networkFile = new NetworkFile(getTokenFromSession(e));

		try {
			resp = await networkFile.getAll(requestParams.toObject(), e.fetch);
			if (resp.ok) {
				items = (await resp.json()) as Items<File>;
			}
		} catch (err) {
			throw Error(resp?.statusText);
		}

		return {
			status: resp?.status || 200,
			stuff: stuff,
			props: {
				items: items,
				params: requestParams
			}
		};
	};
</script>

<script lang="ts">
	import type { Items } from '$elven/types';
	import Pagination from '$lib/components/pagination.svelte';
	import type { File } from '$elven/types/file';
	import FilesToolbars from '$elven/components/files_toolbars.svelte';
	import FilesList from '$elven/components/files_list.svelte';
	import NetworkFile from '$elven/network/file';
	import {
		HandleRouteParam,
		Params,
		Refresh,
		type RPH_Data,
		type RPH_Event
	} from '$elven/tools/params';
	import type { Load } from '@sveltejs/kit';

	import { getTokenFromSession } from '$elven/tools';
	import ItemsContainer from '$elven/components/items_container.svelte';

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
	async function onParamChanged(event: RPH_Event<File>) {
		const data: RPH_Data<File> = {
			params: params,
			items: items
		};
		await HandleRouteParam<File>(event, data);
	}

	async function refresh() {
		const getData = async () => {
			return Promise.resolve({
				items: items,
				params: params
			});
		};
		await Refresh<File>(getData);
	}
</script>

<ItemsContainer>
	<div slot="up">
		<FilesToolbars
			bind:params
			on:uploaded={async () => await onUploaded()}
			on:paramChanged={async (e) => await onParamChanged(e.detail)}
		/>
	</div>

	<div slot="list">
		<FilesList {items} on:deleted={async () => await refresh()} />
	</div>
	
	<div slot="pages">
		{#if items && items.meta}
			<Pagination
				total={items.meta.total_pages}
				current={params.getParam('page')}
				on:changed={async (e) => await onPageChanged(e.detail)}
			/>
		{/if}
	</div>
</ItemsContainer>
